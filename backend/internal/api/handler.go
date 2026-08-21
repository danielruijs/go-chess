package api

import (
	"context"
	"fmt"
	"net/http"

	"go-chess/internal/db/sqlc"
)

type Handler struct {
	matchHandler *matchHandler
	userHandler  *userHandler
}

func NewHandler(ctx context.Context, queries *sqlc.Queries) (*Handler, error) {
	matchHandler, err := newMatchHandler(ctx, queries)
	if err != nil {
		return nil, fmt.Errorf("failed to create match handler: %w", err)
	}

	return &Handler{
		matchHandler: matchHandler,
		userHandler:  newUserHandler(queries),
	}, nil
}

func (h *Handler) GetMatch(w http.ResponseWriter, r *http.Request) {
	h.matchHandler.getMatch(w, r)
}

func (h *Handler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	h.userHandler.getUserProfile(w, r)
}
