package chess

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStrToSquare(t *testing.T) {
	validSquares := []string{"a1", "h8", "d4", "e5", "b7", "g2"}
	invalidSquares := []string{"a0", "i5", "d9", "z3", "aa", "1b", "", "e"}

	for _, s := range validSquares {
		sq, err := StrToSquare(s)
		assert.Nil(t, err)
		assert.Equal(t, Square(s), sq)
	}

	for _, s := range invalidSquares {
		_, err := StrToSquare(s)
		assert.NotNil(t, err)
	}
}

func TestIsValid(t *testing.T) {
	validSquares := []string{"a1", "h8", "d4", "e5", "b7", "g2"}
	invalidSquares := []string{"a0", "i5", "d9", "z3", "aa", "1b", "", "e"}

	for _, s := range validSquares {
		sq, err := StrToSquare(s)
		assert.Nil(t, err)
		assert.Equal(t, true, sq.IsValid())
	}

	for _, s := range invalidSquares {
		_, err := StrToSquare(s)
		assert.NotNil(t, err)
		assert.Equal(t, false, Square(s).IsValid())
	}
}
