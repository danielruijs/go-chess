package chess

import (
	"fmt"
)

type Engine struct {
	position  *Position
	generator *Generator
	moves     []Move
	positions map[PositionKey]int
}

func NewEngine() *Engine {
	e := &Engine{
		position:  NewInitialPosition(),
		generator: NewGenerator(),
		moves:     []Move{},
		positions: make(map[PositionKey]int),
	}
	e.positions[e.position.Key()]++
	return e
}

func (e *Engine) GetBoard() Board {
	return e.position.GetBoard()
}

func (e *Engine) GetLegalMoves(color Color) []Move {
	if color != e.position.ActiveColor {
		return []Move{}
	}
	return e.generator.GenerateMoves(e.position, color)
}

func (e *Engine) ApplyMove(move Move, color Color) (*Result, error) {
	err := e.validateMove(move, color)
	if err != nil {
		return nil, fmt.Errorf("illegal move: %w", err)
	}

	e.position.MakeMove(move)

	e.positions[e.position.Key()]++
	e.moves = append(e.moves, move)

	result := e.GetResult()
	return result, nil
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
	if !ContainsMove(moves, move) {
		return fmt.Errorf("move is not legal")
	}
	return nil
}
