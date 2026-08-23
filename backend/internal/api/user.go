package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"go-chess/internal/chess"
	"go-chess/internal/db/sqlc"
	"go-chess/internal/server"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type userHandler struct {
	queries *sqlc.Queries
}

func newUserHandler(queries *sqlc.Queries) *userHandler {
	return &userHandler{
		queries: queries,
	}
}

type UserInfo struct {
	Username    string    `json:"username"`
	DisplayName string    `json:"displayName"`
	CreatedAt   time.Time `json:"createdAt"`
}

type GameRecord struct {
	Wins   int32 `json:"wins"`
	Losses int32 `json:"losses"`
	Draws  int32 `json:"draws"`
}

type UserStats struct {
	White GameRecord `json:"white"`
	Black GameRecord `json:"black"`
}

type UserMatchItem struct {
	PublicID            string              `json:"publicId"`
	PlayedColor         chess.Color         `json:"playedColor"`
	OpponentDisplayName string              `json:"opponentDisplayName"`
	OpponentUsername    string              `json:"opponentUsername,omitempty"`
	Result              chess.Result        `json:"result"`
	TimeFormat          server.TimeFormatMs `json:"timeFormat"`
	MoveCount           int64               `json:"moveCount"`
	CreatedAt           time.Time           `json:"createdAt"`
}

type UserProfileResponse struct {
	User    UserInfo        `json:"user"`
	Stats   UserStats       `json:"stats"`
	Matches []UserMatchItem `json:"matches"`
}

func (h *userHandler) getUserProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username := r.PathValue("username")
	if username == "" {
		http.Error(w, "Missing username", http.StatusBadRequest)
		return
	}

	user, err := h.queries.GetUserByUsername(r.Context(), username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		log.Printf("ERROR [Handler]: failed to get user %q: %v", username, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	userID := pgtype.Int8{Int64: user.ID, Valid: true}
	matchRows, err := h.queries.GetUserEndedMatches(r.Context(), userID)
	if err != nil {
		log.Printf("ERROR [Handler]: failed to get matches for user %d: %v", user.ID, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	matches, stats := summarizeMatches(user.ID, matchRows)

	resp := UserProfileResponse{
		User: UserInfo{
			Username:    user.Username,
			DisplayName: user.DisplayName,
			CreatedAt:   user.CreatedAt.Time,
		},
		Stats:   stats,
		Matches: matches,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("ERROR [Handler]: failed to encode profile response: %v", err)
	}
}

func summarizeMatches(userID int64, matchRows []sqlc.GetUserEndedMatchesRow) ([]UserMatchItem, UserStats) {
	var stats UserStats
	matches := make([]UserMatchItem, 0, len(matchRows))

	for _, matchRow := range matchRows {
		match, err := buildUserMatchItem(userID, matchRow)
		if err != nil {
			log.Printf("ERROR [Handler]: failed to build match item %q for user %d: %v", matchRow.PublicID, userID, err)
			continue
		}
		matches = append(matches, match)

		switch match.PlayedColor {
		case chess.White:
			switch match.Result.Outcome {
			case chess.WhiteWin:
				stats.White.Wins++
			case chess.BlackWin:
				stats.White.Losses++
			case chess.Draw:
				stats.White.Draws++
			}
		case chess.Black:
			switch match.Result.Outcome {
			case chess.BlackWin:
				stats.Black.Wins++
			case chess.WhiteWin:
				stats.Black.Losses++
			case chess.Draw:
				stats.Black.Draws++
			}
		}
	}

	return matches, stats
}

func buildUserMatchItem(userID int64, matchRow sqlc.GetUserEndedMatchesRow) (UserMatchItem, error) {
	gameEnded, err := server.ParseGameEndedPayload(matchRow.EndedPayload)
	if err != nil {
		return UserMatchItem{}, fmt.Errorf("failed to parse game ended payload: %w", err)
	}

	userIsWhite := matchRow.WhiteUserID.Valid && matchRow.WhiteUserID.Int64 == userID
	userIsBlack := matchRow.BlackUserID.Valid && matchRow.BlackUserID.Int64 == userID

	if !userIsWhite && !userIsBlack {
		return UserMatchItem{}, fmt.Errorf("user %d is not a player in match %s", userID, matchRow.PublicID)
	}

	playedColor := chess.White
	opponentDisplayName := matchRow.BlackDisplayName
	opponentUsername := matchRow.BlackUsername.String
	if !userIsWhite {
		playedColor = chess.Black
		opponentDisplayName = matchRow.WhiteDisplayName
		opponentUsername = matchRow.WhiteUsername.String
	}

	return UserMatchItem{
		PublicID:            matchRow.PublicID,
		PlayedColor:         playedColor,
		OpponentDisplayName: opponentDisplayName,
		OpponentUsername:    opponentUsername,
		Result:              gameEnded.Result,
		TimeFormat: server.TimeFormatMs{
			InitialMs:   matchRow.InitialTimeMs,
			IncrementMs: matchRow.IncrementMs,
		},
		MoveCount: matchRow.MoveCount,
		CreatedAt: matchRow.CreatedAt.Time,
	}, nil
}
