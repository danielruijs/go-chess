package chess

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateMoveNoPiece(t *testing.T) {
	e := NewEngine()
	move := Move{From: 20, To: 28} // e3 to e4
	isLegal := e.isMoveLegal(move, White)
	assert.False(t, isLegal)
}

func TestValidateMoveWrongActiveColor(t *testing.T) {
	e := NewEngine()
	move := Move{From: 52, To: 44} // e7 to e6
	isLegal := e.isMoveLegal(move, Black)
	assert.False(t, isLegal)
}

func TestValidateMoveNotPlayersPiece(t *testing.T) {
	e := NewEngine()
	move := Move{From: 52, To: 44} // e7 to e6
	isLegal := e.isMoveLegal(move, White)
	assert.False(t, isLegal)
}
