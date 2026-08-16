package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"go-chess/internal/cache"
	"go-chess/internal/chess"
	"go-chess/internal/db/sqlc"
	"go-chess/internal/server"

	"github.com/jackc/pgx/v5"
)

const (
	matchCacheCleanupInterval = 1 * time.Hour
	matchCacheTTL             = 24 * time.Hour
)

type MatchHandler struct {
	queries   *sqlc.Queries
	generator *chess.Generator
	cache     *cache.Cache[string, []byte]
}

func NewMatchHandler(ctx context.Context, queries *sqlc.Queries) (*MatchHandler, error) {
	cache, err := cache.New[string](cache.Options[[]byte]{
		Cleanup: &cache.CleanupConfig[[]byte]{
			Interval: matchCacheCleanupInterval,
			ShouldEvict: func(_ []byte, lastUsed time.Time) bool {
				return time.Since(lastUsed) > matchCacheTTL
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create match cache: %w", err)
	}
	cache.StartCleanup(ctx)

	return &MatchHandler{
		queries:   queries,
		generator: chess.NewGenerator(),
		cache:     cache,
	}, nil
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
		http.Error(w, "Missing match ID", http.StatusBadRequest)
		return
	}
	if !server.IsValidPublicMatchID(publicID) {
		http.Error(w, "Invalid match ID", http.StatusBadRequest)
		return
	}

	if cached, ok := h.cache.Get(publicID); ok {
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write(cached); err != nil {
			log.Printf("ERROR [MatchHandler]: failed to write cached response: %v", err)
		}
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

	respBytes, err := json.Marshal(resp)
	if err != nil {
		log.Printf("ERROR [MatchHandler]: failed to encode response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	h.cache.Set(publicID, respBytes)

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(respBytes); err != nil {
		log.Printf("ERROR [MatchHandler]: failed to write response: %v", err)
	}
}
