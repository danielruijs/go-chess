package server

import (
	"fmt"
	"go-chess/internal/chess"
)

func moveDataToMove(data MoveData) (chess.Move, error) {
	from, err := chess.StrToSquare(data.From)
	if err != nil {
		return chess.Move{}, fmt.Errorf("invalid from square: %v", err)
	}
	to, err := chess.StrToSquare(data.To)
	if err != nil {
		return chess.Move{}, fmt.Errorf("invalid to square: %v", err)
	}

	return chess.Move{
		From:      from,
		To:        to,
		Promotion: data.Promotion,
	}, nil
}

func moveListToLegalMoves(moves []chess.Move) map[string][]LegalMove {
	legalMoves := make(map[string][]LegalMove)
	for _, move := range moves {
		fromStr := chess.SquareToStr(move.From)
		legalMove := LegalMove{
			To:        chess.SquareToStr(move.To),
			Promotion: move.Promotion,
		}
		legalMoves[fromStr] = append(legalMoves[fromStr], legalMove)
	}
	return legalMoves
}
