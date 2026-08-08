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

type moveWithColor struct {
	move  Move
	color Color
}

func TestApplyMoveFullMatchWhiteCheckmate(t *testing.T) {
	e := NewEngine()

	moves := []moveWithColor{
		{move: Move{From: strToSquare("e2"), To: strToSquare("e4")}, color: White},
		{move: Move{From: strToSquare("e7"), To: strToSquare("e5")}, color: Black},
		{move: Move{From: strToSquare("f1"), To: strToSquare("c4")}, color: White},
		{move: Move{From: strToSquare("b8"), To: strToSquare("c6")}, color: Black},
		{move: Move{From: strToSquare("d1"), To: strToSquare("h5")}, color: White},
		{move: Move{From: strToSquare("g8"), To: strToSquare("f6")}, color: Black},
		{move: Move{From: strToSquare("h5"), To: strToSquare("f7")}, color: White},
	}

	for i, m := range moves {
		result, err := e.ApplyMove(m.move, m.color)
		assert.Nil(t, err)
		if i < len(moves)-1 {
			assert.Nil(t, result)
			continue
		}

		assert.Equal(t, &Result{Outcome: WhiteWin, Reason: Checkmate}, result)
		e.ApplyResult(*result)
	}

	assert.Equal(t, PGN("1.e4 e5 2.Bc4 Nc6 3.Qh5 Nf6 4.Qxf7# 1-0"), e.GetPGN())
}

func TestApplyMoveFullMatchBlackCheckmate(t *testing.T) {
	e := NewEngine()

	moves := []moveWithColor{
		{move: Move{From: strToSquare("f2"), To: strToSquare("f3")}, color: White},
		{move: Move{From: strToSquare("e7"), To: strToSquare("e6")}, color: Black},
		{move: Move{From: strToSquare("g2"), To: strToSquare("g4")}, color: White},
		{move: Move{From: strToSquare("d8"), To: strToSquare("h4")}, color: Black},
	}

	for i, m := range moves {
		result, err := e.ApplyMove(m.move, m.color)
		assert.Nil(t, err)
		if i < len(moves)-1 {
			assert.Nil(t, result)
			continue
		}

		assert.Equal(t, &Result{Outcome: BlackWin, Reason: Checkmate}, result)
		e.ApplyResult(*result)
	}

	assert.Equal(t, PGN("1.f3 e6 2.g4 Qh4# 0-1"), e.GetPGN())
}

func TestApplyMoveFullMatchThreefoldRepetitionNonConsecutive(t *testing.T) {
	e := NewEngine()

	moves := []moveWithColor{
		{move: Move{From: strToSquare("e2"), To: strToSquare("e4")}, color: White},
		{move: Move{From: strToSquare("e7"), To: strToSquare("e5")}, color: Black},
		{move: Move{From: strToSquare("e1"), To: strToSquare("e2")}, color: White},
		{move: Move{From: strToSquare("e8"), To: strToSquare("e7")}, color: Black},
		{move: Move{From: strToSquare("e2"), To: strToSquare("e1")}, color: White},
		{move: Move{From: strToSquare("e7"), To: strToSquare("e8")}, color: Black},
		{move: Move{From: strToSquare("b1"), To: strToSquare("c3")}, color: White},
		{move: Move{From: strToSquare("b8"), To: strToSquare("c6")}, color: Black},
		{move: Move{From: strToSquare("e1"), To: strToSquare("e2")}, color: White},
		{move: Move{From: strToSquare("e8"), To: strToSquare("e7")}, color: Black},
		{move: Move{From: strToSquare("e2"), To: strToSquare("e1")}, color: White},
		{move: Move{From: strToSquare("e7"), To: strToSquare("e8")}, color: Black},
		{move: Move{From: strToSquare("c3"), To: strToSquare("b1")}, color: White},
		{move: Move{From: strToSquare("c6"), To: strToSquare("b8")}, color: Black},
		{move: Move{From: strToSquare("e1"), To: strToSquare("e2")}, color: White},
		{move: Move{From: strToSquare("e8"), To: strToSquare("e7")}, color: Black},
		{move: Move{From: strToSquare("e2"), To: strToSquare("e1")}, color: White},
		{move: Move{From: strToSquare("e7"), To: strToSquare("e8")}, color: Black},
	}

	for i, m := range moves {
		result, err := e.ApplyMove(m.move, m.color)
		assert.Nil(t, err)
		if i < len(moves)-1 {
			assert.Nil(t, result)
			continue
		}

		assert.Equal(t, &Result{Outcome: Draw, Reason: ThreefoldRepetition}, result)
		e.ApplyResult(*result)
	}

	assert.Equal(
		t,
		PGN("1.e4 e5 2.Ke2 Ke7 3.Ke1 Ke8 4.Nc3 Nc6 5.Ke2 Ke7 6.Ke1 Ke8 7.Nb1 Nb8 8.Ke2 Ke7 9.Ke1 Ke8 1/2-1/2"),
		e.GetPGN(),
	)
}
