package chess

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFenToPositionStart(t *testing.T) {
	pos, err := StartingPositionFEN.ToPosition()

	board := pos.Board
	assert.Equal(t, nil, err)
	assert.Equal(t, BoardSize, len(board))
	for _, row := range board {
		assert.Equal(t, BoardSize, len(row))
	}
	assert.Equal(t, Rook, board[0][0].Type)
	assert.Equal(t, Black, board[0][0].Color)
	assert.Equal(t, Queen, board[0][3].Type)
	assert.Equal(t, Black, board[0][3].Color)
	assert.Equal(t, King, board[7][4].Type)
	assert.Equal(t, White, board[7][4].Color)
	assert.Equal(t, Pawn, board[6][7].Type)
	assert.Equal(t, White, board[6][7].Color)

	assert.Equal(t, White, pos.ActiveColor)
	assert.Equal(t, true, pos.CastlingRights.WhiteOO)
	assert.Equal(t, true, pos.CastlingRights.WhiteOOO)
	assert.Equal(t, true, pos.CastlingRights.BlackOO)
	assert.Equal(t, true, pos.CastlingRights.BlackOOO)
	assert.Equal(t, Square(""), pos.EnPassant)
	assert.Equal(t, uint(0), pos.Halfmove)
	assert.Equal(t, uint(1), pos.Fullmove)
}

func TestFenToPositionEndgame(t *testing.T) {
	fen := Fen("3Q4/5pk1/pp4pp/2pB4/P1Pn4/2N3PP/5KP1/8 b - - 0 37")
	pos, err := fen.ToPosition()

	board := pos.Board
	assert.Equal(t, BoardSize, len(board))
	for _, row := range board {
		assert.Equal(t, BoardSize, len(row))
	}

	assert.Equal(t, Queen, board[0][3].Type)
	assert.Equal(t, White, board[0][3].Color)
	assert.Equal(t, King, board[1][6].Type)
	assert.Equal(t, Black, board[1][6].Color)
	assert.Equal(t, Pawn, board[2][0].Type)
	assert.Equal(t, Black, board[2][0].Color)
	assert.Equal(t, Knight, board[4][3].Type)
	assert.Equal(t, Black, board[4][3].Color)
	assert.Equal(t, King, board[6][5].Type)
	assert.Equal(t, White, board[6][5].Color)

	assert.Equal(t, Black, pos.ActiveColor)
	assert.Equal(t, false, pos.CastlingRights.WhiteOO)
	assert.Equal(t, false, pos.CastlingRights.WhiteOOO)
	assert.Equal(t, false, pos.CastlingRights.BlackOO)
	assert.Equal(t, false, pos.CastlingRights.BlackOOO)
	assert.Equal(t, Square(""), pos.EnPassant)
	assert.Equal(t, uint(0), pos.Halfmove)
	assert.Equal(t, uint(37), pos.Fullmove)

	assert.Equal(t, nil, err)
}

func TestFenToPositionInvalid(t *testing.T) {
	invalidFens := []Fen{
		"8/8/8/8/ w KQkq - 0 1",                  // Not enough rows
		"8/8/8/8/8/8/8/8/8 w KQkq - 0 1",         // Too many rows
		"8/8/8/8/8/8/8/9 w KQkq - 0 1",           // Too many columns
		"8/8/8/8/ppppppppppp/8/8/8 w KQkq - 0 1", // Too many columns
		"8/8/8/8/8/8/8/7X w KQkq - 0 1",          // Invalid character
		"8/8/8/8/8/8/8/8 x KQkq - 0 1",           // Invalid active color
		"8/8/8/8/8/8/8/8 w KQkq q1 0 1",          // Invalid en passant square
		"8/8/8/8/8/8/8/8 w KQkq e3 -1 1",         // Invalid halfmove clock
		"8/8/8/8/8/8/8/8 w KQkq e3 0 0",          // Invalid fullmove number
		"invalid fen string",                     // Completely invalid format
	}
	for _, fen := range invalidFens {
		_, err := fen.ToPosition()
		assert.NotEqual(t, nil, err, "Expected error for invalid FEN: %s", fen)
	}
}
