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
		color  Color
	}{
		{0, "a1", 0, 0, 1 << 0, Black},
		{1, "b1", 1, 0, 1 << 1, White},
		{14, "g2", 6, 1, 1 << 14, White},
		{27, "d4", 3, 3, 1 << 27, Black},
		{36, "e5", 4, 4, 1 << 36, Black},
		{49, "b7", 1, 6, 1 << 49, White},
		{63, "h8", 7, 7, 1 << 63, Black},
	}

	for _, tt := range squareTests {
		// StrToSquare
		sq, err := StrToSquare(tt.str)
		assert.Nil(t, err)
		if sq != tt.square {
			t.Errorf("StrToSquare(%s) = %d; want %d", tt.str, sq, tt.square)
		}

		// SquareToStr
		str := SquareToStr(tt.square)
		if str != tt.str {
			t.Errorf("SquareToStr(%d) = %s; want %s", tt.square, str, tt.str)
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

		// Color
		if c := tt.square.Color(); c != tt.color {
			t.Errorf("Square(%d).Color() = %v; want %v", tt.square, c, tt.color)
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
