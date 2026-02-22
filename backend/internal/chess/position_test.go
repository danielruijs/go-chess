package chess

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPiece(t *testing.T) {
	position := NewInitialPosition()

	assert.Equal(t, Piece{Type: Rook, Color: White}, position.GetPiece(0))
	assert.Equal(t, Piece{Type: Knight, Color: White}, position.GetPiece(1))
	assert.Equal(t, Piece{Type: Bishop, Color: White}, position.GetPiece(2))
	assert.Equal(t, Piece{Type: Queen, Color: White}, position.GetPiece(3))
	assert.Equal(t, Piece{Type: King, Color: White}, position.GetPiece(4))

	assert.Equal(t, Piece{Type: Pawn, Color: White}, position.GetPiece(8))
	assert.Equal(t, Piece{Type: Pawn, Color: White}, position.GetPiece(15))

	assert.Equal(t, Piece{}, position.GetPiece(16))
	assert.Equal(t, Piece{}, position.GetPiece(27))
	assert.Equal(t, Piece{}, position.GetPiece(47))

	assert.Equal(t, Piece{Type: Pawn, Color: Black}, position.GetPiece(48))
	assert.Equal(t, Piece{Type: Pawn, Color: Black}, position.GetPiece(55))

	assert.Equal(t, Piece{Type: Queen, Color: Black}, position.GetPiece(59))
	assert.Equal(t, Piece{Type: King, Color: Black}, position.GetPiece(60))
	assert.Equal(t, Piece{Type: Bishop, Color: Black}, position.GetPiece(61))
	assert.Equal(t, Piece{Type: Knight, Color: Black}, position.GetPiece(62))
	assert.Equal(t, Piece{Type: Rook, Color: Black}, position.GetPiece(63))
}

func TestGetBoard(t *testing.T) {
	position := NewInitialPosition()

	board := position.GetBoard()

	assert.Equal(t, 8, len(board))
	for _, rank := range board {
		assert.Equal(t, 8, len(rank))
	}

	assert.Equal(t, Rook, board[0][0].Type)
	assert.Equal(t, White, board[0][0].Color)
	assert.Equal(t, Queen, board[3][0].Type)
	assert.Equal(t, White, board[3][0].Color)
	assert.Equal(t, King, board[4][7].Type)
	assert.Equal(t, Black, board[4][7].Color)
	assert.Equal(t, Pawn, board[2][6].Type)
	assert.Equal(t, Black, board[2][6].Color)
}

func TestGetCopy(t *testing.T) {
	position := NewInitialPosition()
	position.Halfmove = 5
	position.Fullmove = 10
	enPassantSquare := Square(16)
	position.EnPassant = &enPassantSquare

	copy := position.GetCopy()

	assert.Equal(t, position.WhitePawns, copy.WhitePawns)
	assert.Equal(t, position.BlackPawns, copy.BlackPawns)

	assert.Equal(t, position.ActiveColor, copy.ActiveColor)
	assert.Equal(t, position.CastlingRights, copy.CastlingRights)
	assert.Equal(t, position.Halfmove, copy.Halfmove)
	assert.Equal(t, position.Fullmove, copy.Fullmove)

	assert.NotNil(t, copy.EnPassant)
	assert.Equal(t, *position.EnPassant, *copy.EnPassant)
	assert.NotSame(t, position.EnPassant, copy.EnPassant)

	// Changing en passant should not change original
	*copy.EnPassant = 20
	assert.Equal(t, Square(16), *position.EnPassant)

	// Changing a piece should not change original
	copy.WhitePawns |= squareMask(24)
	assert.Equal(t, Bitboard(1<<24), copy.WhitePawns&squareMask(24))
	assert.Equal(t, Bitboard(0), position.WhitePawns&squareMask(24))
}

func TestGetCopyWithNilEnPassant(t *testing.T) {
	position := NewInitialPosition()

	copy := position.GetCopy()

	assert.Nil(t, copy.EnPassant)
}
