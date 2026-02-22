package chess

import (
	"fmt"
	"slices"
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

func (e *Engine) GetLegalMoves(color Color) []Move {
	return e.generator.GenerateMoves(e.position, color)
}

func (e *Engine) ApplyMove(move Move, color Color) error {
	err := e.validateMove(move, color)
	if err != nil {
		return fmt.Errorf("illegal move: %w", err)
	}
	e.position.MakeMove(move)
	e.moves = append(e.moves, move)
	return nil
}

func (e *Engine) validateMove(move Move, color Color) error {
	pieceToMove := e.position.GetPiece(move.From)
	if pieceToMove == nil {
		return fmt.Errorf("no piece at from square")
	}
	if color != e.position.ActiveColor {
		return fmt.Errorf("wrong active color")
	}
	if color != pieceToMove.Color {
		return fmt.Errorf("piece does not belong to player")
	}

	moves := e.generator.GeneratePieceMoves(e.position, color, pieceToMove.Type)
	if !slices.Contains(moves, move) {
		return fmt.Errorf("move is not legal")
	}
	return nil
}
