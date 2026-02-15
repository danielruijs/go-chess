package chess

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFenToPositionStart(t *testing.T) {
	pos, err := StartingPositionFEN.ToPosition()

	board := pos.Board
	assert.Nil(t, err)
	assert.Equal(t, BoardSize, len(board))
	for _, rank := range board {
		assert.Equal(t, BoardSize, len(rank))
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
	assert.Nil(t, err)
	assert.Equal(t, BoardSize, len(board))
	for _, rank := range board {
		assert.Equal(t, BoardSize, len(rank))
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
}

func TestFenToPositionEnPassantWhite(t *testing.T) {
	fen := Fen("rnbqkbnr/ppp2ppp/3p4/3Pp3/8/8/PPP1PPPP/RNBQKBNR w KQkq e6 0 3")

	pos, err := fen.ToPosition()

	board := pos.Board
	assert.Nil(t, err)
	assert.Equal(t, BoardSize, len(board))
	for _, rank := range board {
		assert.Equal(t, BoardSize, len(rank))
	}

	assert.Equal(t, Pawn, board[3][3].Type)
	assert.Equal(t, White, board[3][3].Color)
	assert.Equal(t, Pawn, board[3][4].Type)
	assert.Equal(t, Black, board[3][4].Color)

	assert.Equal(t, White, pos.ActiveColor)
	assert.Equal(t, true, pos.CastlingRights.WhiteOO)
	assert.Equal(t, true, pos.CastlingRights.WhiteOOO)
	assert.Equal(t, true, pos.CastlingRights.BlackOO)
	assert.Equal(t, true, pos.CastlingRights.BlackOOO)
	assert.Equal(t, Square("e6"), pos.EnPassant)
	assert.Equal(t, uint(0), pos.Halfmove)
	assert.Equal(t, uint(3), pos.Fullmove)
}

func TestFenToPositionEnPassantBlack(t *testing.T) {
	fen := Fen("rnbqkbnr/ppp1pppp/8/8/P2pP3/8/1PPP1PPP/RNBQKBNR b KQkq e3 0 3")

	pos, err := fen.ToPosition()

	board := pos.Board
	assert.Nil(t, err)
	assert.Equal(t, BoardSize, len(board))
	for _, rank := range board {
		assert.Equal(t, BoardSize, len(rank))
	}

	assert.Equal(t, Pawn, board[4][4].Type)
	assert.Equal(t, White, board[4][4].Color)
	assert.Equal(t, Pawn, board[4][3].Type)
	assert.Equal(t, Black, board[4][3].Color)

	assert.Equal(t, Black, pos.ActiveColor)
	assert.Equal(t, true, pos.CastlingRights.WhiteOO)
	assert.Equal(t, true, pos.CastlingRights.WhiteOOO)
	assert.Equal(t, true, pos.CastlingRights.BlackOO)
	assert.Equal(t, true, pos.CastlingRights.BlackOOO)
	assert.Equal(t, Square("e3"), pos.EnPassant)
	assert.Equal(t, uint(0), pos.Halfmove)
	assert.Equal(t, uint(3), pos.Fullmove)
}

func TestFenToPositionInvalid(t *testing.T) {
	invalidFens := []Fen{
		"8/8/8/8/ w KQkq - 0 1",                  // Not enough ranks
		"8/8/8/8/8/8/8/8/8 w KQkq - 0 1",         // Too many ranks
		"8/8/8/8/8/8/8/9 w KQkq - 0 1",           // Too many files
		"8/8/8/8/ppppppppppp/8/8/8 w KQkq - 0 1", // Too many files
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
