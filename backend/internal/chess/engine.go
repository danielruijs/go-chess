package chess

import (
	"fmt"
)

type Engine struct {
	position  *Position
	generator *Generator
	moves     []Move
}

func NewEngine() *Engine {
	return &Engine{
		position:  NewInitialPosition(),
		generator: NewGenerator(),
		moves:     []Move{},
	}
}

func (e *Engine) GetBoard() Board {
	return e.position.GetBoard()
}

func (e *Engine) ApplyMove(move Move, color Color) error {
	if !e.isMoveLegal(move, color) {
		return fmt.Errorf("illegal move")
	}
	// TODO: apply move to board
	e.moves = append(e.moves, move)
	return nil
}

func (e *Engine) isMoveLegal(move Move, color Color) bool {
	pieceToMove := e.position.GetPiece(move.From)
	if pieceToMove == (Piece{}) {
		return false
	}
	if color != e.position.ActiveColor {
		return false
	}
	if pieceToMove.Color != color {
		return false
	}

	moves := e.generator.GeneratePieceMoves(e.position, color, pieceToMove.Type)
	for _, m := range moves {
		if m == move {
			return true
		}
	}
	return false
}
