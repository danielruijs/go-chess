package chess

import (
	"fmt"
)

type Engine struct {
	position  *Position
	generator *Generator
	pgn       PGN
	positions map[PositionKey]int
}

func NewEngine() *Engine {
	e := &Engine{
		position:  NewInitialPosition(),
		generator: NewGenerator(),
		pgn:       "",
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

func (e *Engine) GetPGN() PGN {
	return e.pgn
}

func (e *Engine) GetActiveColor() Color {
	return e.position.ActiveColor
}

func (e *Engine) ApplyMove(move Move, color Color) (*Result, error) {
	err := e.validateMove(move, color)
	if err != nil {
		return nil, fmt.Errorf("illegal move: %w", err)
	}

	e.pgn = e.pgn.AppendMove(move, e.position, e.generator)
	e.position.MakeMove(move)
	e.positions[e.position.Key()]++

	result := getResult(e.generator, e.position, e.positions)
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

func (e *Engine) ApplyResult(result *Result) {
	if result == nil {
		return
	}
	e.pgn = e.pgn.AppendResult(result)
}
