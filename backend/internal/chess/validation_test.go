package chess

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateMoveNoPiece(t *testing.T) {
	position := NewInitialPosition()
	move := Move{From: "e3", To: "e4"}
	err := position.ValidateMove(move, White)
	assert.NotNil(t, err)
}

func TestValidateMoveWrongActiveColor(t *testing.T) {
	position := NewInitialPosition()
	move := Move{From: "e7", To: "e6"}
	err := position.ValidateMove(move, Black)
	assert.NotNil(t, err)
}

func TestValidateMoveNotPlayersPiece(t *testing.T) {
	position := NewInitialPosition()
	move := Move{From: "e7", To: "e6"}
	err := position.ValidateMove(move, White)
	assert.NotNil(t, err)
}

func TestValidatePawnMove(t *testing.T) {
	tests := []struct {
		name  string
		fen   Fen
		move  Move
		color Color
		valid bool
	}{
		{
			name:  "one square forward - white",
			fen:   StartingPositionFEN,
			move:  Move{From: "e2", To: "e3"},
			color: White,
			valid: true,
		},
		{
			name:  "one square forward - black",
			fen:   StartingPositionFEN,
			move:  Move{From: "e7", To: "e6"},
			color: Black,
			valid: true,
		},
		{
			name:  "two squares forward - white",
			fen:   StartingPositionFEN,
			move:  Move{From: "e2", To: "e4"},
			color: White,
			valid: true,
		},
		{
			name:  "two squares forward - black",
			fen:   StartingPositionFEN,
			move:  Move{From: "e7", To: "e5"},
			color: Black,
			valid: true,
		},
		{
			name:  "too far - white",
			fen:   StartingPositionFEN,
			move:  Move{From: "e2", To: "e5"},
			color: White,
			valid: false,
		},
		{
			name:  "too far - black",
			fen:   StartingPositionFEN,
			move:  Move{From: "e7", To: "e4"},
			color: Black,
			valid: false,
		},
		{
			name:  "nothing to capture - white",
			fen:   StartingPositionFEN,
			move:  Move{From: "e2", To: "d3"},
			color: White,
			valid: false,
		},
		{
			name:  "nothing to capture - black",
			fen:   StartingPositionFEN,
			move:  Move{From: "e7", To: "f6"},
			color: Black,
			valid: false,
		},
		{
			name:  "cant move backwards - white",
			fen:   "8/7P/1k6/4pp2/4PP2/1K6/3p4/8 w - - 0 1",
			move:  Move{From: "e4", To: "e3"},
			color: White,
			valid: false,
		},
		{
			name:  "cant move backwards - black",
			fen:   "8/7P/1k6/4pp2/4PP2/1K6/3p4/8 w - - 0 1",
			move:  Move{From: "e5", To: "e6"},
			color: Black,
			valid: false,
		},
		{
			name:  "cant move sideways - white",
			fen:   "8/7P/1k6/4pp2/4PP2/1K6/3p4/8 w - - 0 1",
			move:  Move{From: "e4", To: "d4"},
			color: White,
			valid: false,
		},
		{
			name:  "cant move sideways - black",
			fen:   "8/7P/1k6/4pp2/4PP2/1K6/3p4/8 w - - 0 1",
			move:  Move{From: "e5", To: "d5"},
			color: Black,
			valid: false,
		},
		{
			name:  "valid capture - white",
			fen:   "8/7P/1k6/4pp2/4PP2/1K6/3p4/8 w - - 0 1",
			move:  Move{From: "e4", To: "f5"},
			color: White,
			valid: true,
		},
		{
			name:  "valid capture - black",
			fen:   "8/7P/1k6/4pp2/4PP2/1K6/3p4/8 w - - 0 1",
			move:  Move{From: "e5", To: "f4"},
			color: Black,
			valid: true,
		},
		{
			name:  "valid promotion - white",
			fen:   "8/7P/1k6/4pp2/4PP2/1K6/3p4/8 w - - 0 1",
			move:  Move{From: "h7", To: "h8"},
			color: White,
			valid: true,
		},
		{
			name:  "valid promotion - black",
			fen:   "8/7P/1k6/4pp2/4PP2/1K6/3p4/8 w - - 0 1",
			move:  Move{From: "d2", To: "d1"},
			color: Black,
			valid: true,
		},
		{
			name:  "cant move forward if blocked - white",
			fen:   "8/7P/1k6/4pp2/4PP2/1K6/3p4/8 w - - 0 1",
			move:  Move{From: "f4", To: "f5"},
			color: White,
			valid: false,
		},
		{
			name:  "cant move forward if blocked - black",
			fen:   "8/7P/1k6/4pp2/4PP2/1K6/3p4/8 w - - 0 1",
			move:  Move{From: "f5", To: "f4"},
			color: Black,
			valid: false,
		},
		{
			name:  "valid en passant - white",
			fen:   "rnbqkbnr/ppp2ppp/3p4/3Pp3/8/8/PPP1PPPP/RNBQKBNR w KQkq e6 0 3",
			move:  Move{From: "d5", To: "e6"},
			color: White,
			valid: true,
		},
		{
			name:  "invalid en passant - white",
			fen:   "rnbqkbnr/pppp1ppp/8/3Pp3/8/8/PPP1PPPP/RNBQKBNR w KQkq - 0 3",
			move:  Move{From: "d5", To: "e6"},
			color: White,
			valid: false,
		},
		{
			name:  "valid en passant - black",
			fen:   "rnbqkbnr/ppp1pppp/8/8/P2pP3/8/1PPP1PPP/RNBQKBNR b KQkq e3 0 3",
			move:  Move{From: "d4", To: "e3"},
			color: Black,
			valid: true,
		},
		{
			name:  "invalid en passant - black",
			fen:   "rnbqkbnr/ppp1pppp/8/8/3pP3/P7/1PPP1PPP/RNBQKBNR b KQkq - 0 3",
			move:  Move{From: "d4", To: "e3"},
			color: Black,
			valid: false,
		},
	}

	for _, tt := range tests {
		position, err := tt.fen.ToPosition()
		if err != nil {
			t.Fatalf("%s: failed to parse FEN: %v", tt.name, err)
		}
		position.ActiveColor = tt.color // ignore active color
		err = position.ValidateMove(tt.move, tt.color)
		if tt.valid && err != nil {
			t.Errorf("%s: should be valid, got error: %v", tt.name, err)
		} else if !tt.valid && err == nil {
			t.Errorf("%s: should not be valid, got no error", tt.name)
		}
	}
}
