package chess

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func toPtr[T any](v T) *T {
	return &v
}

func strToSquare(str string) Square {
	file := int(str[0] - 'a')
	rank := int(str[1] - '1')
	return Square(file + rank*8)
}

func bitboardFromStrs(squareStrings []string) Bitboard {
	var bb Bitboard
	for _, sqStr := range squareStrings {
		sq := strToSquare(sqStr)
		bb |= squareMask(sq)
	}
	return bb
}

func TestPawnMoveGenerator_generateMoves(t *testing.T) {
	tests := []struct {
		name       string
		pos        *Position
		color      Color
		legalMoves []Move
	}{
		{
			name: "white single push",
			pos: &Position{
				WhitePawns: bitboardFromStrs([]string{"a3"}),
			},
			color: White,
			legalMoves: []Move{
				{From: strToSquare("a3"), To: strToSquare("a4")}, // a3 -> a4
			},
		},
		{
			name: "black single push",
			pos: &Position{
				BlackPawns: bitboardFromStrs([]string{"a6"}),
			},
			color: Black,
			legalMoves: []Move{
				{From: strToSquare("a6"), To: strToSquare("a5")}, // a6 -> a5
			},
		},
		{
			name: "white single push blocked",
			pos: &Position{
				WhitePawns: bitboardFromStrs([]string{"a3"}),
				BlackPawns: bitboardFromStrs([]string{"a4"}),
			},
			color:      White,
			legalMoves: []Move{},
		},
		{
			name: "black single push blocked",
			pos: &Position{
				WhitePawns: bitboardFromStrs([]string{"a5"}),
				BlackPawns: bitboardFromStrs([]string{"a6"}),
			},
			color:      Black,
			legalMoves: []Move{},
		},
		{
			name: "white double push",
			pos: &Position{
				WhitePawns: bitboardFromStrs([]string{"a2"}),
			},
			color: White,
			legalMoves: []Move{
				{From: strToSquare("a2"), To: strToSquare("a3")},
				{From: strToSquare("a2"), To: strToSquare("a4")},
			},
		},
		{
			name: "black double push",
			pos: &Position{
				BlackPawns: bitboardFromStrs([]string{"a7"}),
			},
			color: Black,
			legalMoves: []Move{
				{From: strToSquare("a7"), To: strToSquare("a6")},
				{From: strToSquare("a7"), To: strToSquare("a5")},
			},
		},
		{
			name: "white double push blocked",
			pos: &Position{
				WhitePawns: bitboardFromStrs([]string{"a2"}),
				BlackPawns: bitboardFromStrs([]string{"a4"}),
			},
			color: White,
			legalMoves: []Move{
				{From: strToSquare("a2"), To: strToSquare("a3")},
			},
		},
		{
			name: "black double push blocked",
			pos: &Position{
				WhitePawns: bitboardFromStrs([]string{"a5"}),
				BlackPawns: bitboardFromStrs([]string{"a7"}),
			},
			color: Black,
			legalMoves: []Move{
				{From: strToSquare("a7"), To: strToSquare("a6")},
			},
		},
		{
			name: "white capture",
			pos: &Position{
				WhitePawns: bitboardFromStrs([]string{"d4"}),
				BlackPawns: bitboardFromStrs([]string{"c5", "e5"}),
			},
			color: White,
			legalMoves: []Move{
				{From: strToSquare("d4"), To: strToSquare("d5")},
				{From: strToSquare("d4"), To: strToSquare("c5")},
				{From: strToSquare("d4"), To: strToSquare("e5")},
			},
		},
		{
			name: "black capture",
			pos: &Position{
				WhitePawns: bitboardFromStrs([]string{"b4", "d4"}),
				BlackPawns: bitboardFromStrs([]string{"c5"}),
			},
			color: Black,
			legalMoves: []Move{
				{From: strToSquare("c5"), To: strToSquare("c4")},
				{From: strToSquare("c5"), To: strToSquare("b4")},
				{From: strToSquare("c5"), To: strToSquare("d4")},
			},
		},
		{
			name: "white en passant",
			pos: &Position{
				WhitePawns: bitboardFromStrs([]string{"d5"}),
				BlackPawns: bitboardFromStrs([]string{"e5"}),
				EnPassant:  toPtr(strToSquare("e6")),
			},
			color: White,
			legalMoves: []Move{
				{From: strToSquare("d5"), To: strToSquare("d6")},
				{From: strToSquare("d5"), To: strToSquare("e6")},
			},
		},
		{
			name: "black en passant",
			pos: &Position{
				WhitePawns: bitboardFromStrs([]string{"d4"}),
				BlackPawns: bitboardFromStrs([]string{"e4"}),
				EnPassant:  toPtr(strToSquare("d3")),
			},
			color: Black,
			legalMoves: []Move{
				{From: strToSquare("e4"), To: strToSquare("e3")},
				{From: strToSquare("e4"), To: strToSquare("d3")},
			},
		},
		{
			name: "white promotion",
			pos: &Position{
				WhitePawns: bitboardFromStrs([]string{"a7"}),
			},
			color: White,
			legalMoves: []Move{
				{From: strToSquare("a7"), To: strToSquare("a8"), Promotion: toPtr(Queen)},
				{From: strToSquare("a7"), To: strToSquare("a8"), Promotion: toPtr(Rook)},
				{From: strToSquare("a7"), To: strToSquare("a8"), Promotion: toPtr(Bishop)},
				{From: strToSquare("a7"), To: strToSquare("a8"), Promotion: toPtr(Knight)},
			},
		},
		{
			name: "black promotion",
			pos: &Position{
				BlackPawns: bitboardFromStrs([]string{"a2"}),
			},
			color: Black,
			legalMoves: []Move{
				{From: strToSquare("a2"), To: strToSquare("a1"), Promotion: toPtr(Queen)},
				{From: strToSquare("a2"), To: strToSquare("a1"), Promotion: toPtr(Rook)},
				{From: strToSquare("a2"), To: strToSquare("a1"), Promotion: toPtr(Bishop)},
				{From: strToSquare("a2"), To: strToSquare("a1"), Promotion: toPtr(Knight)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := PawnMoveGenerator{}
			moves := g.generateMoves(tt.pos, tt.color)
			assert.ElementsMatch(t, tt.legalMoves, moves)
		})
	}
}
