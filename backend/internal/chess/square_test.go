package chess

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSquareToStr(t *testing.T) {
	tests := []struct {
		square Square
		str    string
	}{
		{0, "a1"},
		{1, "b1"},
		{49, "b7"},
		{27, "d4"},
		{36, "e5"},
		{14, "g2"},
		{63, "h8"},
	}

	for _, tt := range tests {
		str := SquareToStr(tt.square)
		if str != tt.str {
			t.Errorf("Expected string %s for square %d, got %s", tt.str, tt.square, str)
		}
	}
}

func TestStrToSquare(t *testing.T) {
	tests := []struct {
		str    string
		square Square
	}{
		{"a1", 0},
		{"b1", 1},
		{"b7", 49},
		{"d4", 27},
		{"e5", 36},
		{"g2", 14},
		{"h8", 63},
	}

	for _, tt := range tests {
		sq, err := StrToSquare(tt.str)
		assert.Nil(t, err)
		if sq != tt.square {
			t.Errorf("Expected square %d, got %d for string %s", tt.square, sq, tt.str)
		}
	}
}

func TestIsValid(t *testing.T) {
	validSquares := []string{"a1", "h8", "d4", "e5", "b7", "g2"}
	invalidSquares := []string{"a0", "i5", "d9", "z3", "aa", "1b", "", "e"}

	for _, s := range validSquares {
		valid := isValid(s)
		assert.True(t, valid)
	}

	for _, s := range invalidSquares {
		valid := isValid(s)
		assert.False(t, valid)
	}
}

func TestCoordsToSquare(t *testing.T) {
	tests := []struct {
		file   int
		rank   int
		square Square
	}{
		{0, 0, 0},
		{1, 0, 1},
		{1, 6, 49},
		{3, 3, 27},
		{4, 4, 36},
		{6, 1, 14},
		{7, 7, 63},
	}

	for _, tt := range tests {
		sq := coordsToSquare(tt.file, tt.rank)
		if sq != tt.square {
			t.Errorf("Expected square %d for file %d rank %d, got %d", tt.square, tt.file, tt.rank, sq)
		}
	}
}

func TestSquareMask(t *testing.T) {
	tests := []struct {
		square Square
		mask   Bitboard
	}{
		{0, 1 << 0},
		{1, 1 << 1},
		{49, 1 << 49},
		{27, 1 << 27},
		{36, 1 << 36},
		{14, 1 << 14},
		{63, 1 << 63},
	}

	for _, tt := range tests {
		mask := squareMask(tt.square)
		if mask != tt.mask {
			t.Errorf("Expected mask %d for square %d, got %d", tt.mask, tt.square, mask)
		}
	}
}

func TestCoordMask(t *testing.T) {
	tests := []struct {
		file int
		rank int
		mask Bitboard
	}{
		{0, 0, 1 << 0},
		{1, 0, 1 << 1},
		{1, 6, 1 << 49},
		{3, 3, 1 << 27},
		{4, 4, 1 << 36},
		{6, 1, 1 << 14},
		{7, 7, 1 << 63},
	}

	for _, tt := range tests {
		mask := coordMask(tt.file, tt.rank)
		if mask != tt.mask {
			t.Errorf("Expected mask %d for file %d rank %d, got %d", tt.mask, tt.file, tt.rank, mask)
		}
	}
}

// func TestToCoords(t *testing.T) {
// 	tests := []struct {
// 		square Square
// 		file   int
// 		rank   int
// 	}{
// 		{0, 0, 0},
// 		{1, 1, 0},
// 		{49, 1, 6},
// 		{27, 3, 3},
// 		{36, 4, 4},
// 		{14, 6, 1},
// 		{63, 7, 7},
// 	}

// 	for _, tt := range tests {
// 		file, rank := tt.square.ToCoords()
// 		if file != tt.file || rank != tt.rank {
// 			t.Errorf("Expected file %d, rank %d for square %d, got file %d, rank %d", tt.file, tt.rank, tt.square, file, rank)
// 		}
// 	}
// }
