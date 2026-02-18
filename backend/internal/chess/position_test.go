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

func TestValidateMoveNoPiece(t *testing.T) {
	position := NewInitialPosition()
	move := Move{From: 20, To: 28} // e3 to e4
	err := position.ValidateMove(move, White)
	assert.NotNil(t, err)
}

func TestValidateMoveWrongActiveColor(t *testing.T) {
	position := NewInitialPosition()
	move := Move{From: 52, To: 44} // e7 to e6
	err := position.ValidateMove(move, Black)
	assert.NotNil(t, err)
}

func TestValidateMoveNotPlayersPiece(t *testing.T) {
	position := NewInitialPosition()
	move := Move{From: 52, To: 44} // e7 to e6
	err := position.ValidateMove(move, White)
	assert.NotNil(t, err)
}
