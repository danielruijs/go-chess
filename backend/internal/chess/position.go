package chess

import "fmt"

func (p *Position) ApplyMove(move Move, color Color) error {
	err := p.ValidateMove(move, color)
	if err != nil {
		return fmt.Errorf("invalid move: %w", err)
	}
	fmt.Println("apply move", move.From, move.To)
	// TODO: apply move to board
	return nil
}
