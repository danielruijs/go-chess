package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"go-chess/internal/chess"
	"go-chess/internal/db/sqlc"
	"go-chess/internal/server"

	"github.com/jackc/pgx/v5"
)

type MatchHandler struct {
	queries   *sqlc.Queries
	generator *chess.Generator
}

func NewMatchHandler(queries *sqlc.Queries) *MatchHandler {
	return &MatchHandler{
		queries:   queries,
		generator: chess.NewGenerator(),
	}
}

type Position struct {
	Index       int         `json:"index"`
	Board       chess.Board `json:"board"`
	SAN         string      `json:"san,omitempty"` // omitted for starting position (index 0)
	WhiteTimeMs int64       `json:"whiteTimeMs"`
	BlackTimeMs int64       `json:"blackTimeMs"`
}

type Match struct {
	WhitePlayerName string       `json:"whitePlayerName"`
	BlackPlayerName string       `json:"blackPlayerName"`
	Result          chess.Result `json:"result"`
	Positions       []Position   `json:"positions"`
}

func (h *MatchHandler) GetMatch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	publicID := r.PathValue("publicId")
	if publicID == "" {
		publicID = strings.TrimPrefix(r.URL.Path, "/api/match/")
	}
	if publicID == "" {
		http.Error(w, "Missing match ID", http.StatusBadRequest)
		return
	}

	match, err := h.queries.GetMatchByPublicID(r.Context(), publicID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "Match not found", http.StatusNotFound)
			return
		}
		log.Printf("ERROR [MatchHandler]: failed to get match: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	events, err := h.queries.GetMatchEventsByMatchID(r.Context(), match.ID)
	if err != nil {
		log.Printf("ERROR [MatchHandler]: failed to get match events: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if len(events) == 0 {
		http.Error(w, "Match is in progress", http.StatusConflict)
		return
	}

	lastEvent := events[len(events)-1]
	if lastEvent.EventType != sqlc.MatchEventTypeGameEnded {
		http.Error(w, "Match is in progress", http.StatusConflict)
		return
	}

	gameEndedPayload, err := server.ParseGameEndedPayload(lastEvent.Payload)
	if err != nil {
		log.Printf("ERROR [MatchHandler]: failed to parse game ended payload: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	pos := chess.NewInitialPosition()

	positions := make([]Position, 0, len(events)+1)
	positions = append(positions, Position{
		Index:       0,
		Board:       pos.GetBoard(),
		WhiteTimeMs: match.InitialTimeMs,
		BlackTimeMs: match.InitialTimeMs,
	})

	for _, event := range events {
		if event.EventType != sqlc.MatchEventTypeMove {
			continue
		}

		movePayload, err := server.ParseMovePayload(event.Payload)
		if err != nil {
			log.Printf("ERROR [MatchHandler]: failed to parse move payload: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		move, err := movePayload.ToMove()
		if err != nil {
			log.Printf("ERROR [MatchHandler]: failed to parse move: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		san := move.ToSAN(pos, h.generator)
		pos.MakeMove(move)

		positions = append(positions, Position{
			Index:       len(positions),
			Board:       pos.GetBoard(),
			SAN:         san,
			WhiteTimeMs: movePayload.WhiteTimeMs,
			BlackTimeMs: movePayload.BlackTimeMs,
		})
	}

	resp := Match{
		WhitePlayerName: match.WhiteDisplayName,
		BlackPlayerName: match.BlackDisplayName,
		Result:          gameEndedPayload.Result,
		Positions:       positions,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Printf("ERROR [MatchHandler]: failed to encode response: %v", err)
	}
}
