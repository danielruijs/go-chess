package chess

import (
	"fmt"
	"log"
)

type CastlingRights struct {
	WhiteOO  bool `json:"white_oo"`
	WhiteOOO bool `json:"white_ooo"`
	BlackOO  bool `json:"black_oo"`
	BlackOOO bool `json:"black_ooo"`
}

type Position struct {
	Board          Board          `json:"board"`           // TODO: remove, replace with bitboards
	ActiveColor    Color          `json:"active_color"`    // Color to move
	CastlingRights CastlingRights `json:"castling_rights"` // Castling rights
	EnPassant      Square         `json:"en_passant"`      // En passant target square, square over which pawn just moved when moving two squares
	Halfmove       uint           `json:"halfmove"`        // Halfmove clock, number of halfmoves since last capture or pawn move, for fifty-move rule
	Fullmove       uint           `json:"fullmove"`        // Fullmove number
}

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
	return e.position.Board
}

func (e *Engine) ApplyMove(move Move, color Color) error {
	e.moves = append(e.moves, move)
	fmt.Printf("Apply move %s -> %s\n", move.From, move.To)
	return nil
}
