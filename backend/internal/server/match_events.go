package server

import (
	"encoding/json"
	"fmt"

	"go-chess/internal/chess"
)

type GameStartedPayload struct{}

type MovePayload struct {
	From        string           `json:"from"`
	To          string           `json:"to"`
	Promotion   *chess.PieceType `json:"promotion,omitempty"`
	WhiteTimeMs int64            `json:"whiteTimeMs"`
	BlackTimeMs int64            `json:"blackTimeMs"`
}

func ParseMovePayload(payload []byte) (MovePayload, error) {
	var mp MovePayload
	if err := json.Unmarshal(payload, &mp); err != nil {
		return MovePayload{}, fmt.Errorf("failed to unmarshal move payload: %w", err)
	}
	return mp, nil
}

func (mp MovePayload) ToMove() (chess.Move, error) {
	from, err := chess.StrToSquare(mp.From)
	if err != nil {
		return chess.Move{}, fmt.Errorf("invalid from square %q: %w", mp.From, err)
	}

	to, err := chess.StrToSquare(mp.To)
	if err != nil {
		return chess.Move{}, fmt.Errorf("invalid to square %q: %w", mp.To, err)
	}

	return chess.Move{
		From:      from,
		To:        to,
		Promotion: mp.Promotion,
	}, nil
}

type GameEndedPayload struct {
	Result chess.Result `json:"result"`
}

func ParseGameEndedPayload(payload []byte) (GameEndedPayload, error) {
	var gep GameEndedPayload
	if err := json.Unmarshal(payload, &gep); err != nil {
		return GameEndedPayload{}, fmt.Errorf("failed to unmarshal game ended payload: %w", err)
	}
	return gep, nil
}
