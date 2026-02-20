package chess

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAbs(t *testing.T) {
	tests := []struct {
		name     string
		n        int
		expected int
	}{
		{
			name:     "abs of positive number",
			n:        5,
			expected: 5,
		},
		{
			name:     "abs of negative number",
			n:        -5,
			expected: 5,
		},
		{
			name:     "abs of zero",
			n:        0,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := abs(tt.n)
			assert.Equal(t, tt.expected, result)
		})
	}
}

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

func TestPopLSB(t *testing.T) {
	tests := []struct {
		name          string
		board         Bitboard
		expectedIndex int
		expectedBoard Bitboard
	}{
		{
			name:          "pop LSB from single bit",
			board:         Bitboard(0x1),
			expectedIndex: 0,
			expectedBoard: Bitboard(0x0),
		},
		{
			name:          "pop LSB from multiple bits",
			board:         Bitboard(0xF),
			expectedIndex: 0,
			expectedBoard: Bitboard(0xE),
		},
		{
			name:          "pop LSB from higher bit position",
			board:         Bitboard(0x80),
			expectedIndex: 7,
			expectedBoard: Bitboard(0x0),
		},
		{
			name:          "pop LSB with multiple set bits",
			board:         Bitboard(0x104),
			expectedIndex: 2,
			expectedBoard: Bitboard(0x100),
		},
		{
			name:          "pop LSB from zero bitboard",
			board:         Bitboard(0x0),
			expectedIndex: 64, // bits.TrailingZeros64 returns 64 for x=0
			expectedBoard: Bitboard(0x0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			board := tt.board
			index := popLSB(&board)
			assert.Equal(t, tt.expectedIndex, index)
			assert.Equal(t, tt.expectedBoard, board)
		})
	}
}
