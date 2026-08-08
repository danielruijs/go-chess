package server

import "go-chess/internal/chess"

type GameStartedPayload struct{}

type MovePayload struct {
	From        string           `json:"from"`
	To          string           `json:"to"`
	Promotion   *chess.PieceType `json:"promotion,omitempty"`
	WhiteTimeMs int64            `json:"whiteTimeMs"`
	BlackTimeMs int64            `json:"blackTimeMs"`
}

type GameEndedPayload struct {
	Outcome chess.Outcome `json:"outcome"`
	Reason  chess.Reason  `json:"reason"`
}
