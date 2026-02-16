package chess

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSquareConversions(t *testing.T) {
	squareTests := []struct {
		square Square
		str    string
		file   int
		rank   int
		mask   Bitboard
	}{
		{0, "a1", 0, 0, 1 << 0},
		{1, "b1", 1, 0, 1 << 1},
		{14, "g2", 6, 1, 1 << 14},
		{27, "d4", 3, 3, 1 << 27},
		{36, "e5", 4, 4, 1 << 36},
		{49, "b7", 1, 6, 1 << 49},
		{63, "h8", 7, 7, 1 << 63},
	}

	for _, tt := range squareTests {
		// StrToSquare
		sq, err := StrToSquare(tt.str)
		assert.Nil(t, err)
		if sq != tt.square {
			t.Errorf("StrToSquare(%s) = %d; want %d", tt.str, sq, tt.square)
		}

		// coordsToSquare
		if sq := coordsToSquare(tt.file, tt.rank); sq != tt.square {
			t.Errorf("coordsToSquare(%d,%d) = %d; want %d", tt.file, tt.rank, sq, tt.square)
		}

		// squareMask
		if m := squareMask(tt.square); m != tt.mask {
			t.Errorf("squareMask(%d) = %d; want %d", tt.square, m, tt.mask)
		}

		// coordMask
		if m := coordMask(tt.file, tt.rank); m != tt.mask {
			t.Errorf("coordMask(%d,%d) = %d; want %d", tt.file, tt.rank, m, tt.mask)
		}
	}
}

func TestIsValidStr(t *testing.T) {
	validSquares := []string{"a1", "h8", "d4", "e5", "b7", "g2"}
	invalidSquares := []string{"a0", "i5", "d9", "z3", "aa", "1b", "", "e"}

	for _, s := range validSquares {
		valid := isValidStr(s)
		assert.True(t, valid)
	}

	for _, s := range invalidSquares {
		valid := isValidStr(s)
		assert.False(t, valid)
	}
}
