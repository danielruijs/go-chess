package chess

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateMoveNoPiece(t *testing.T) {
	e := NewEngine()
	move := Move{From: strToSquare("e3"), To: strToSquare("e4")}
	err := e.validateMove(move, White)
	assert.NotNil(t, err)
}

func TestValidateMoveWrongActiveColor(t *testing.T) {
	e := NewEngine()
	move := Move{From: strToSquare("e7"), To: strToSquare("e6")}
	err := e.validateMove(move, Black)
	assert.NotNil(t, err)
}

func TestValidateMoveNotPlayersPiece(t *testing.T) {
	e := NewEngine()
	move := Move{From: strToSquare("e7"), To: strToSquare("e6")}
	err := e.validateMove(move, White)
	assert.NotNil(t, err)
}
