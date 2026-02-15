package chess

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPiece(t *testing.T) {
	board := NewInitialPosition().Board

	tests := []struct {
		name     string
		square   Square
		expected Piece
	}{
		{
			name:   "a1 should be white rook",
			square: "a1",
			expected: Piece{
				Type:  Rook,
				Color: White,
			},
		},
		{
			name:   "d1 should be white queen",
			square: "d1",
			expected: Piece{
				Type:  Queen,
				Color: White,
			},
		},
		{
			name:   "e1 should be white king",
			square: "e1",
			expected: Piece{
				Type:  King,
				Color: White,
			},
		},
		{
			name:   "g1 should be white knight",
			square: "g1",
			expected: Piece{
				Type:  Knight,
				Color: White,
			},
		},
		{
			name:   "b2 should be white pawn",
			square: "b2",
			expected: Piece{
				Type:  Pawn,
				Color: White,
			},
		},
		{
			name:   "a8 should be black rook",
			square: "a8",
			expected: Piece{
				Type:  Rook,
				Color: Black,
			},
		},
		{
			name:   "d8 should be black queen",
			square: "d8",
			expected: Piece{
				Type:  Queen,
				Color: Black,
			},
		},
		{
			name:   "e8 should be black king",
			square: "e8",
			expected: Piece{
				Type:  King,
				Color: Black,
			},
		},
		{
			name:   "g8 should be black knight",
			square: "g8",
			expected: Piece{
				Type:  Knight,
				Color: Black,
			},
		},
		{
			name:   "b7 should be black pawn",
			square: "b7",
			expected: Piece{
				Type:  Pawn,
				Color: Black,
			},
		},
		{
			name:     "d4 should be empty",
			square:   "d4",
			expected: Piece{},
		},
	}

	for _, tt := range tests {
		piece, err := board.GetPiece(tt.square)
		assert.Nil(t, err)
		if piece != tt.expected {
			t.Errorf("%s: expected piece: %v, got: %v", tt.name, tt.expected, piece)
		}
	}

	// invalid square
	_, err := board.GetPiece("i5")
	assert.NotNil(t, err)
}
