package chess

import (
	"fmt"
	"log"
)

type Engine struct {
	position *Position
	moves    []Move
}

func NewInitialPosition() *Position {
	pos, err := StartingPositionFEN.ToPosition()
	if err != nil {
		log.Fatalf("failed to create initial position: %v", err)
	}
	return &pos
}

func NewEngine() Engine {
	return Engine{
		position: NewInitialPosition(),
		moves:    []Move{},
	}
}

func (e *Engine) GetBoard() Board {
	return e.position.GetBoard()
}

func (e *Engine) ApplyMove(move Move, color Color) error {
	err := e.position.ValidateMove(move, color)
	if err != nil {
		return fmt.Errorf("invalid move: %w", err)
	}
	// TODO: apply move to board
	e.moves = append(e.moves, move)
	return nil
}
