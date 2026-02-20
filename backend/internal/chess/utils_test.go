package chess

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShift(t *testing.T) {
	tests := []struct {
		name     string
		board    Bitboard
		n        int
		expected Bitboard
	}{
		{
			name:     "shift left",
			board:    Bitboard(0x1),
			n:        3,
			expected: Bitboard(0x8),
		},
		{
			name:     "shift right",
			board:    Bitboard(0x8),
			n:        -3,
			expected: Bitboard(0x1),
		},
		{
			name:     "shift by 0",
			board:    Bitboard(0x4),
			n:        0,
			expected: Bitboard(0x4),
		},
		{
			name:     "shift zero bitboard",
			board:    Bitboard(0x0),
			n:        5,
			expected: Bitboard(0x0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := shift(tt.board, tt.n)
			assert.Equal(t, tt.expected, result)
		})
	}
}
