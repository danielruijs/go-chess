package chess

import "fmt"

type Validator interface {
	ValidateMove(move Move, position *Position, color Color) error
}

var pieceToValidator = map[PieceType]Validator{
	Pawn: PawnMoveValidator{},
}

func (p *Position) ValidateMove(move Move, color Color) error {
	pieceToMove, err := p.Board.GetPiece(move.From)
	if err != nil {
		return fmt.Errorf("failed to get piece at source square: %w", err)
	}
	if pieceToMove == (Piece{}) {
		return fmt.Errorf("no piece at source square")
	}
	if color != p.ActiveColor {
		return fmt.Errorf("not %s's turn to move", color)
	}
	if pieceToMove.Color != color {
		return fmt.Errorf("piece at source square does not belong to player")
	}

	validator := pieceToValidator[pieceToMove.Type]
	return validator.ValidateMove(move, p, color)
}

type PawnMoveValidator struct{}

func (v PawnMoveValidator) ValidateMove(move Move, position *Position, color Color) error {
	// move forward by 1 if not blocked
	// move forward by 2 if on starting rank and not blocked
	// capture diagonally by 1 if occupied by opponent piece
	// en passant
	// promotion
	return nil
}
