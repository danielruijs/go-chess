package chess

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFenToPositionStart(t *testing.T) {
	pos, err := StartingPositionFEN.ToPosition()

	assert.Nil(t, err)

	// Rooks
	assert.NotZero(t, pos.WhiteRooks&coordMask(0, 0))
	assert.NotZero(t, pos.WhiteRooks&coordMask(7, 0))
	assert.NotZero(t, pos.BlackRooks&coordMask(0, 7))
	assert.NotZero(t, pos.BlackRooks&coordMask(7, 7))

	// Knights
	assert.NotZero(t, pos.WhiteKnights&coordMask(1, 0))
	assert.NotZero(t, pos.WhiteKnights&coordMask(6, 0))
	assert.NotZero(t, pos.BlackKnights&coordMask(1, 7))
	assert.NotZero(t, pos.BlackKnights&coordMask(6, 7))

	// Bishops
	assert.NotZero(t, pos.WhiteBishops&coordMask(2, 0))
	assert.NotZero(t, pos.WhiteBishops&coordMask(5, 0))
	assert.NotZero(t, pos.BlackBishops&coordMask(2, 7))
	assert.NotZero(t, pos.BlackBishops&coordMask(5, 7))

	// Queens
	assert.NotZero(t, pos.WhiteQueens&coordMask(3, 0))
	assert.NotZero(t, pos.BlackQueens&coordMask(3, 7))

	// Kings
	assert.NotZero(t, pos.WhiteKing&coordMask(4, 0))
	assert.NotZero(t, pos.BlackKing&coordMask(4, 7))

	// Pawns
	for file := 0; file < 8; file++ {
		assert.NotZero(t, pos.WhitePawns&coordMask(file, 1))
		assert.NotZero(t, pos.BlackPawns&coordMask(file, 6))
	}

	assert.Equal(t, White, pos.ActiveColor)
	assert.Equal(t, true, pos.CastlingRights.WhiteOO)
	assert.Equal(t, true, pos.CastlingRights.WhiteOOO)
	assert.Equal(t, true, pos.CastlingRights.BlackOO)
	assert.Equal(t, true, pos.CastlingRights.BlackOOO)
	assert.Nil(t, pos.EnPassant)
	assert.Equal(t, uint(0), pos.Halfmove)
	assert.Equal(t, uint(1), pos.Fullmove)
}

func TestFenToPositionEndgame(t *testing.T) {
	fen := Fen("3Q4/5pk1/pp4pp/2pB4/P1Pn4/2N3PP/5KP1/8 b - - 0 37")
	pos, err := fen.ToPosition()

	assert.Nil(t, err)

	assert.NotZero(t, pos.WhiteQueens&coordMask(3, 7))
	assert.NotZero(t, pos.BlackKing&coordMask(6, 6))
	assert.NotZero(t, pos.BlackPawns&coordMask(1, 5))
	assert.NotZero(t, pos.BlackKnights&coordMask(3, 3))
	assert.NotZero(t, pos.WhiteKnights&coordMask(2, 2))
	assert.NotZero(t, pos.WhiteKing&coordMask(5, 1))

	assert.Equal(t, Black, pos.ActiveColor)
	assert.Equal(t, false, pos.CastlingRights.WhiteOO)
	assert.Equal(t, false, pos.CastlingRights.WhiteOOO)
	assert.Equal(t, false, pos.CastlingRights.BlackOO)
	assert.Equal(t, false, pos.CastlingRights.BlackOOO)
	assert.Nil(t, pos.EnPassant)
	assert.Equal(t, uint(0), pos.Halfmove)
	assert.Equal(t, uint(37), pos.Fullmove)
}

func TestFenToPositionEnPassantWhite(t *testing.T) {
	fen := Fen("rnbqkbnr/ppp2ppp/3p4/3Pp3/8/8/PPP1PPPP/RNBQKBNR w KQkq e6 0 3")

	pos, err := fen.ToPosition()

	assert.Nil(t, err)

	assert.NotZero(t, pos.WhitePawns&coordMask(3, 4))
	assert.NotZero(t, pos.BlackPawns&coordMask(4, 4))

	assert.Equal(t, White, pos.ActiveColor)
	assert.Equal(t, true, pos.CastlingRights.WhiteOO)
	assert.Equal(t, true, pos.CastlingRights.WhiteOOO)
	assert.Equal(t, true, pos.CastlingRights.BlackOO)
	assert.Equal(t, true, pos.CastlingRights.BlackOOO)
	assert.Equal(t, Square(44), *pos.EnPassant)
	assert.Equal(t, uint(0), pos.Halfmove)
	assert.Equal(t, uint(3), pos.Fullmove)
}

func TestFenToPositionEnPassantBlack(t *testing.T) {
	fen := Fen("rnbqkbnr/ppp1pppp/8/8/P2pP3/8/1PPP1PPP/RNBQKBNR b KQkq e3 0 3")

	pos, err := fen.ToPosition()

	assert.Nil(t, err)

	assert.NotZero(t, pos.WhitePawns&coordMask(4, 3))
	assert.NotZero(t, pos.BlackPawns&coordMask(3, 3))

	assert.Equal(t, Black, pos.ActiveColor)
	assert.Equal(t, true, pos.CastlingRights.WhiteOO)
	assert.Equal(t, true, pos.CastlingRights.WhiteOOO)
	assert.Equal(t, true, pos.CastlingRights.BlackOO)
	assert.Equal(t, true, pos.CastlingRights.BlackOOO)
	assert.Equal(t, Square(20), *pos.EnPassant)
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
