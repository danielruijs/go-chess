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

// func TestValidatePawnMove(t *testing.T) {
// 	tests := []struct {
// 		name  string
// 		fen   Fen
// 		move  Move
// 		color Color
// 		valid bool
// 	}{
// 		{
// 			name:  "one square forward - white",
// 			fen:   StartingPositionFEN,
// 			move:  Move{From: "e2", To: "e3"},
// 			color: White,
// 			valid: true,
// 		},
// 		{
// 			name:  "one square forward - black",
// 			fen:   StartingPositionFEN,
// 			move:  Move{From: "e7", To: "e6"},
// 			color: Black,
// 			valid: true,
// 		},
// 		{
// 			name:  "two squares forward - white",
// 			fen:   StartingPositionFEN,
// 			move:  Move{From: "e2", To: "e4"},
// 			color: White,
// 			valid: true,
// 		},
// 		{
// 			name:  "two squares forward - black",
// 			fen:   StartingPositionFEN,
// 			move:  Move{From: "e7", To: "e5"},
// 			color: Black,
// 			valid: true,
// 		},
// 		{
// 			name:  "too far - white",
// 			fen:   StartingPositionFEN,
// 			move:  Move{From: "e2", To: "e5"},
// 			color: White,
// 			valid: false,
// 		},
// 		{
// 			name:  "too far - black",
// 			fen:   StartingPositionFEN,
// 			move:  Move{From: "e7", To: "e4"},
// 			color: Black,
// 			valid: false,
// 		},
// 		{
// 			name:  "nothing to capture - white",
// 			fen:   StartingPositionFEN,
// 			move:  Move{From: "e2", To: "d3"},
// 			color: White,
// 			valid: false,
// 		},
// 		{
// 			name:  "nothing to capture - black",
// 			fen:   StartingPositionFEN,
// 			move:  Move{From: "e7", To: "f6"},
// 			color: Black,
// 			valid: false,
// 		},
// 		{
// 			name:  "cant move backwards - white",
// 			fen:   "8/7P/1k6/4pp2/4PP2/1K6/3p4/8 w - - 0 1",
// 			move:  Move{From: "e4", To: "e3"},
// 			color: White,
// 			valid: false,
// 		},
// 		{
// 			name:  "cant move backwards - black",
// 			fen:   "8/7P/1k6/4pp2/4PP2/1K6/3p4/8 w - - 0 1",
// 			move:  Move{From: "e5", To: "e6"},
// 			color: Black,
// 			valid: false,
// 		},
// 		{
// 			name:  "cant move sideways - white",
// 			fen:   "8/7P/1k6/4pp2/4PP2/1K6/3p4/8 w - - 0 1",
// 			move:  Move{From: "e4", To: "d4"},
// 			color: White,
// 			valid: false,
// 		},
// 		{
// 			name:  "cant move sideways - black",
// 			fen:   "8/7P/1k6/4pp2/4PP2/1K6/3p4/8 w - - 0 1",
// 			move:  Move{From: "e5", To: "d5"},
// 			color: Black,
// 			valid: false,
// 		},
// 		{
// 			name:  "valid capture - white",
// 			fen:   "8/7P/1k6/4pp2/4PP2/1K6/3p4/8 w - - 0 1",
// 			move:  Move{From: "e4", To: "f5"},
// 			color: White,
// 			valid: true,
// 		},
// 		{
// 			name:  "valid capture - black",
// 			fen:   "8/7P/1k6/4pp2/4PP2/1K6/3p4/8 w - - 0 1",
// 			move:  Move{From: "e5", To: "f4"},
// 			color: Black,
// 			valid: true,
// 		},
// 		{
// 			name:  "valid promotion - white",
// 			fen:   "8/7P/1k6/4pp2/4PP2/1K6/3p4/8 w - - 0 1",
// 			move:  Move{From: "h7", To: "h8"},
// 			color: White,
// 			valid: true,
// 		},
// 		{
// 			name:  "valid promotion - black",
// 			fen:   "8/7P/1k6/4pp2/4PP2/1K6/3p4/8 w - - 0 1",
// 			move:  Move{From: "d2", To: "d1"},
// 			color: Black,
// 			valid: true,
// 		},
// 		{
// 			name:  "cant move forward if blocked - white",
// 			fen:   "8/7P/1k6/4pp2/4PP2/1K6/3p4/8 w - - 0 1",
// 			move:  Move{From: "f4", To: "f5"},
// 			color: White,
// 			valid: false,
// 		},
// 		{
// 			name:  "cant move forward if blocked - black",
// 			fen:   "8/7P/1k6/4pp2/4PP2/1K6/3p4/8 w - - 0 1",
// 			move:  Move{From: "f5", To: "f4"},
// 			color: Black,
// 			valid: false,
// 		},
// 		{
// 			name:  "valid en passant - white",
// 			fen:   "rnbqkbnr/ppp2ppp/3p4/3Pp3/8/8/PPP1PPPP/RNBQKBNR w KQkq e6 0 3",
// 			move:  Move{From: "d5", To: "e6"},
// 			color: White,
// 			valid: true,
// 		},
// 		{
// 			name:  "invalid en passant - white",
// 			fen:   "rnbqkbnr/pppp1ppp/8/3Pp3/8/8/PPP1PPPP/RNBQKBNR w KQkq - 0 3",
// 			move:  Move{From: "d5", To: "e6"},
// 			color: White,
// 			valid: false,
// 		},
// 		{
// 			name:  "valid en passant - black",
// 			fen:   "rnbqkbnr/ppp1pppp/8/8/P2pP3/8/1PPP1PPP/RNBQKBNR b KQkq e3 0 3",
// 			move:  Move{From: "d4", To: "e3"},
// 			color: Black,
// 			valid: true,
// 		},
// 		{
// 			name:  "invalid en passant - black",
// 			fen:   "rnbqkbnr/ppp1pppp/8/8/3pP3/P7/1PPP1PPP/RNBQKBNR b KQkq - 0 3",
// 			move:  Move{From: "d4", To: "e3"},
// 			color: Black,
// 			valid: false,
// 		},
// 	}

// 	for _, tt := range tests {
// 		position, err := tt.fen.ToPosition()
// 		if err != nil {
// 			t.Fatalf("%s: failed to parse FEN: %v", tt.name, err)
// 		}
// 		position.ActiveColor = tt.color // ignore active color
// 		err = position.ValidateMove(tt.move, tt.color)
// 		if tt.valid && err != nil {
// 			t.Errorf("%s: should be valid, got error: %v", tt.name, err)
// 		} else if !tt.valid && err == nil {
// 			t.Errorf("%s: should not be valid, got no error", tt.name)
// 		}
// 	}
// }
