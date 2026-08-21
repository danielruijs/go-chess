package api

import (
	"context"
	"fmt"
	"net/http"

	"go-chess/internal/db/sqlc"
)

type Handler struct {
	matchHandler *matchHandler
}

func NewHandler(ctx context.Context, queries *sqlc.Queries) (*Handler, error) {
	matchHandler, err := newMatchHandler(ctx, queries)
	if err != nil {
		return nil, fmt.Errorf("failed to create match handler: %w", err)
	}

	return &Handler{
		matchHandler: matchHandler,
	}, nil
}

func (h *Handler) GetMatch(w http.ResponseWriter, r *http.Request) {
	h.matchHandler.getMatch(w, r)
}
