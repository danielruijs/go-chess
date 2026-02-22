package chess

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// unsafe version of StrToSquare, only for tests
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

func TestGeneratePawnMoves(t *testing.T) {
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
			name: "white double push blocked, single allowed",
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
			name: "black double push blocked, single allowed",
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
			name: "white double push blocked",
			pos: &Position{
				WhitePawns: bitboardFromStrs([]string{"a2"}),
				BlackPawns: bitboardFromStrs([]string{"a3"}),
			},
			color:      White,
			legalMoves: []Move{},
		},
		{
			name: "black double push blocked",
			pos: &Position{
				WhitePawns: bitboardFromStrs([]string{"a5"}),
				BlackPawns: bitboardFromStrs([]string{"a6"}),
			},
			color:      Black,
			legalMoves: []Move{},
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
				EnPassant:  new(strToSquare("e6")),
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
				EnPassant:  new(strToSquare("d3")),
			},
			color: Black,
			legalMoves: []Move{
				{From: strToSquare("e4"), To: strToSquare("e3")},
				{From: strToSquare("e4"), To: strToSquare("d3")},
			},
		},
		{
			name: "white en passant cant capture own pawn",
			pos: &Position{
				WhitePawns: bitboardFromStrs([]string{"d4", "e2"}),
				EnPassant:  new(strToSquare("d3")),
			},
			color: White,
			legalMoves: []Move{
				{From: strToSquare("d4"), To: strToSquare("d5")},
				{From: strToSquare("e2"), To: strToSquare("e3")},
				{From: strToSquare("e2"), To: strToSquare("e4")},
			},
		},
		{
			name: "white promotion",
			pos: &Position{
				WhitePawns: bitboardFromStrs([]string{"a7"}),
			},
			color: White,
			legalMoves: []Move{
				{From: strToSquare("a7"), To: strToSquare("a8"), Promotion: new(Queen)},
				{From: strToSquare("a7"), To: strToSquare("a8"), Promotion: new(Rook)},
				{From: strToSquare("a7"), To: strToSquare("a8"), Promotion: new(Bishop)},
				{From: strToSquare("a7"), To: strToSquare("a8"), Promotion: new(Knight)},
			},
		},
		{
			name: "black promotion",
			pos: &Position{
				BlackPawns: bitboardFromStrs([]string{"a2"}),
			},
			color: Black,
			legalMoves: []Move{
				{From: strToSquare("a2"), To: strToSquare("a1"), Promotion: new(Queen)},
				{From: strToSquare("a2"), To: strToSquare("a1"), Promotion: new(Rook)},
				{From: strToSquare("a2"), To: strToSquare("a1"), Promotion: new(Bishop)},
				{From: strToSquare("a2"), To: strToSquare("a1"), Promotion: new(Knight)},
			},
		},
		{
			name: "white multiple pawns",
			pos: &Position{
				WhitePawns: bitboardFromStrs([]string{"a2", "b3", "c4"}),
			},
			color: White,
			legalMoves: []Move{
				{From: strToSquare("a2"), To: strToSquare("a3")},
				{From: strToSquare("a2"), To: strToSquare("a4")},
				{From: strToSquare("b3"), To: strToSquare("b4")},
				{From: strToSquare("c4"), To: strToSquare("c5")},
			},
		},
		{
			name: "black multiple pawns",
			pos: &Position{
				BlackPawns: bitboardFromStrs([]string{"a7", "b6", "c5"}),
			},
			color: Black,
			legalMoves: []Move{
				{From: strToSquare("a7"), To: strToSquare("a6")},
				{From: strToSquare("a7"), To: strToSquare("a5")},
				{From: strToSquare("b6"), To: strToSquare("b5")},
				{From: strToSquare("c5"), To: strToSquare("c4")},
			},
		},
	}

	g := NewGenerator()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			moves := g.generatePawnMoves(tt.pos, tt.color)
			assert.ElementsMatch(t, tt.legalMoves, moves)
		})
	}
}

func TestGenerateKnightMoves(t *testing.T) {
	tests := []struct {
		name       string
		pos        *Position
		color      Color
		legalMoves []Move
	}{
		{
			name: "white clear board",
			pos: &Position{
				WhiteKnights: bitboardFromStrs([]string{"e4"}),
			},
			color: White,
			legalMoves: []Move{
				{From: strToSquare("e4"), To: strToSquare("f6")},
				{From: strToSquare("e4"), To: strToSquare("g5")},
				{From: strToSquare("e4"), To: strToSquare("g3")},
				{From: strToSquare("e4"), To: strToSquare("f2")},
				{From: strToSquare("e4"), To: strToSquare("d2")},
				{From: strToSquare("e4"), To: strToSquare("c3")},
				{From: strToSquare("e4"), To: strToSquare("c5")},
				{From: strToSquare("e4"), To: strToSquare("d6")},
			},
		},
		{
			name: "black clear board",
			pos: &Position{
				BlackKnights: bitboardFromStrs([]string{"d4"}),
			},
			color: Black,
			legalMoves: []Move{
				{From: strToSquare("d4"), To: strToSquare("e6")},
				{From: strToSquare("d4"), To: strToSquare("f5")},
				{From: strToSquare("d4"), To: strToSquare("f3")},
				{From: strToSquare("d4"), To: strToSquare("e2")},
				{From: strToSquare("d4"), To: strToSquare("c2")},
				{From: strToSquare("d4"), To: strToSquare("b3")},
				{From: strToSquare("d4"), To: strToSquare("b5")},
				{From: strToSquare("d4"), To: strToSquare("c6")},
			},
		},
		{
			name: "white corner",
			pos: &Position{
				WhiteKnights: bitboardFromStrs([]string{"a1"}),
			},
			color: White,
			legalMoves: []Move{
				{From: strToSquare("a1"), To: strToSquare("b3")},
				{From: strToSquare("a1"), To: strToSquare("c2")},
			},
		},
		{
			name: "black corner",
			pos: &Position{
				BlackKnights: bitboardFromStrs([]string{"a8"}),
			},
			color: Black,
			legalMoves: []Move{
				{From: strToSquare("a8"), To: strToSquare("b6")},
				{From: strToSquare("a8"), To: strToSquare("c7")},
			},
		},
		{
			name: "white edge",
			pos: &Position{
				WhiteKnights: bitboardFromStrs([]string{"a4"}),
			},
			color: White,
			legalMoves: []Move{
				{From: strToSquare("a4"), To: strToSquare("b6")},
				{From: strToSquare("a4"), To: strToSquare("c5")},
				{From: strToSquare("a4"), To: strToSquare("c3")},
				{From: strToSquare("a4"), To: strToSquare("b2")},
			},
		},
		{
			name: "black edge",
			pos: &Position{
				BlackKnights: bitboardFromStrs([]string{"e8"}),
			},
			color: Black,
			legalMoves: []Move{
				{From: strToSquare("e8"), To: strToSquare("g7")},
				{From: strToSquare("e8"), To: strToSquare("f6")},
				{From: strToSquare("e8"), To: strToSquare("d6")},
				{From: strToSquare("e8"), To: strToSquare("c7")},
			},
		},
		{
			name: "white blocked own pieces",
			pos: &Position{
				WhiteKnights: bitboardFromStrs([]string{"e4"}),
				WhitePawns:   bitboardFromStrs([]string{"f6", "g5", "c3"}),
			},
			color: White,
			legalMoves: []Move{
				{From: strToSquare("e4"), To: strToSquare("g3")},
				{From: strToSquare("e4"), To: strToSquare("f2")},
				{From: strToSquare("e4"), To: strToSquare("d2")},
				{From: strToSquare("e4"), To: strToSquare("c5")},
				{From: strToSquare("e4"), To: strToSquare("d6")},
			},
		},
		{
			name: "black blocked own pieces",
			pos: &Position{
				BlackKnights: bitboardFromStrs([]string{"d4"}),
				BlackPawns:   bitboardFromStrs([]string{"e6", "f5", "e2"}),
			},
			color: Black,
			legalMoves: []Move{
				{From: strToSquare("d4"), To: strToSquare("f3")},
				{From: strToSquare("d4"), To: strToSquare("c2")},
				{From: strToSquare("d4"), To: strToSquare("b3")},
				{From: strToSquare("d4"), To: strToSquare("b5")},
				{From: strToSquare("d4"), To: strToSquare("c6")},
			},
		},
		{
			name: "white blocked enemy pieces",
			pos: &Position{
				WhiteKnights: bitboardFromStrs([]string{"e4"}),
				BlackPawns:   bitboardFromStrs([]string{"f6", "g5"}),
			},
			color: White,
			legalMoves: []Move{
				{From: strToSquare("e4"), To: strToSquare("f6")},
				{From: strToSquare("e4"), To: strToSquare("g5")},
				{From: strToSquare("e4"), To: strToSquare("g3")},
				{From: strToSquare("e4"), To: strToSquare("f2")},
				{From: strToSquare("e4"), To: strToSquare("d2")},
				{From: strToSquare("e4"), To: strToSquare("c3")},
				{From: strToSquare("e4"), To: strToSquare("c5")},
				{From: strToSquare("e4"), To: strToSquare("d6")},
			},
		},
		{
			name: "black blocked enemy pieces",
			pos: &Position{
				BlackKnights: bitboardFromStrs([]string{"d4"}),
				WhitePawns:   bitboardFromStrs([]string{"e6", "f5"}),
			},
			color: Black,
			legalMoves: []Move{
				{From: strToSquare("d4"), To: strToSquare("e6")},
				{From: strToSquare("d4"), To: strToSquare("f5")},
				{From: strToSquare("d4"), To: strToSquare("f3")},
				{From: strToSquare("d4"), To: strToSquare("e2")},
				{From: strToSquare("d4"), To: strToSquare("c2")},
				{From: strToSquare("d4"), To: strToSquare("b3")},
				{From: strToSquare("d4"), To: strToSquare("b5")},
				{From: strToSquare("d4"), To: strToSquare("c6")},
			},
		},
		{
			name: "white multiple knights",
			pos: &Position{
				WhiteKnights: bitboardFromStrs([]string{"e4", "d6"}),
			},
			color: White,
			legalMoves: []Move{
				{From: strToSquare("e4"), To: strToSquare("f6")},
				{From: strToSquare("e4"), To: strToSquare("g5")},
				{From: strToSquare("e4"), To: strToSquare("g3")},
				{From: strToSquare("e4"), To: strToSquare("f2")},
				{From: strToSquare("e4"), To: strToSquare("d2")},
				{From: strToSquare("e4"), To: strToSquare("c3")},
				{From: strToSquare("e4"), To: strToSquare("c5")},
				{From: strToSquare("d6"), To: strToSquare("e8")},
				{From: strToSquare("d6"), To: strToSquare("f7")},
				{From: strToSquare("d6"), To: strToSquare("f5")},
				{From: strToSquare("d6"), To: strToSquare("c4")},
				{From: strToSquare("d6"), To: strToSquare("b5")},
				{From: strToSquare("d6"), To: strToSquare("b7")},
				{From: strToSquare("d6"), To: strToSquare("c8")},
			},
		},
		{
			name: "black multiple knights",
			pos: &Position{
				BlackKnights: bitboardFromStrs([]string{"d4", "c2"}),
				WhitePawns:   bitboardFromStrs([]string{"e6", "f5"}),
			},
			color: Black,
			legalMoves: []Move{
				{From: strToSquare("d4"), To: strToSquare("e6")},
				{From: strToSquare("d4"), To: strToSquare("f5")},
				{From: strToSquare("d4"), To: strToSquare("f3")},
				{From: strToSquare("d4"), To: strToSquare("e2")},
				{From: strToSquare("d4"), To: strToSquare("b3")},
				{From: strToSquare("d4"), To: strToSquare("b5")},
				{From: strToSquare("d4"), To: strToSquare("c6")},
				{From: strToSquare("c2"), To: strToSquare("e3")},
				{From: strToSquare("c2"), To: strToSquare("e1")},
				{From: strToSquare("c2"), To: strToSquare("a1")},
				{From: strToSquare("c2"), To: strToSquare("a3")},
				{From: strToSquare("c2"), To: strToSquare("b4")},
			},
		},
	}

	g := NewGenerator()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			moves := g.generateKnightMoves(tt.pos, tt.color)
			assert.ElementsMatch(t, tt.legalMoves, moves)
		})
	}
}

func TestGenerateBishopMoves(t *testing.T) {
	tests := []struct {
		name       string
		pos        *Position
		color      Color
		legalMoves []Move
	}{
		{
			name: "white clear board",
			pos: &Position{
				WhiteBishops: bitboardFromStrs([]string{"e4"}),
			},
			color: White,
			legalMoves: []Move{
				{From: strToSquare("e4"), To: strToSquare("d5")},
				{From: strToSquare("e4"), To: strToSquare("c6")},
				{From: strToSquare("e4"), To: strToSquare("b7")},
				{From: strToSquare("e4"), To: strToSquare("a8")},
				{From: strToSquare("e4"), To: strToSquare("f5")},
				{From: strToSquare("e4"), To: strToSquare("g6")},
				{From: strToSquare("e4"), To: strToSquare("h7")},
				{From: strToSquare("e4"), To: strToSquare("d3")},
				{From: strToSquare("e4"), To: strToSquare("c2")},
				{From: strToSquare("e4"), To: strToSquare("b1")},
				{From: strToSquare("e4"), To: strToSquare("f3")},
				{From: strToSquare("e4"), To: strToSquare("g2")},
				{From: strToSquare("e4"), To: strToSquare("h1")},
			},
		},
		{
			name: "black clear board",
			pos: &Position{
				BlackBishops: bitboardFromStrs([]string{"d4"}),
			},
			color: Black,
			legalMoves: []Move{
				{From: strToSquare("d4"), To: strToSquare("c5")},
				{From: strToSquare("d4"), To: strToSquare("b6")},
				{From: strToSquare("d4"), To: strToSquare("a7")},
				{From: strToSquare("d4"), To: strToSquare("e5")},
				{From: strToSquare("d4"), To: strToSquare("f6")},
				{From: strToSquare("d4"), To: strToSquare("g7")},
				{From: strToSquare("d4"), To: strToSquare("h8")},
				{From: strToSquare("d4"), To: strToSquare("c3")},
				{From: strToSquare("d4"), To: strToSquare("b2")},
				{From: strToSquare("d4"), To: strToSquare("a1")},
				{From: strToSquare("d4"), To: strToSquare("e3")},
				{From: strToSquare("d4"), To: strToSquare("f2")},
				{From: strToSquare("d4"), To: strToSquare("g1")},
			},
		},
		{
			name: "white blocked own pieces",
			pos: &Position{
				WhiteBishops: bitboardFromStrs([]string{"e4"}),
				WhitePawns:   bitboardFromStrs([]string{"c2", "f5"}),
			},
			color: White,
			legalMoves: []Move{
				{From: strToSquare("e4"), To: strToSquare("d5")},
				{From: strToSquare("e4"), To: strToSquare("c6")},
				{From: strToSquare("e4"), To: strToSquare("b7")},
				{From: strToSquare("e4"), To: strToSquare("a8")},
				{From: strToSquare("e4"), To: strToSquare("d3")},
				{From: strToSquare("e4"), To: strToSquare("f3")},
				{From: strToSquare("e4"), To: strToSquare("g2")},
				{From: strToSquare("e4"), To: strToSquare("h1")},
			},
		},
		{
			name: "black blocked own pieces",
			pos: &Position{
				BlackBishops: bitboardFromStrs([]string{"d4"}),
				BlackPawns:   bitboardFromStrs([]string{"e3", "f6"}),
			},
			color: Black,
			legalMoves: []Move{
				{From: strToSquare("d4"), To: strToSquare("c5")},
				{From: strToSquare("d4"), To: strToSquare("b6")},
				{From: strToSquare("d4"), To: strToSquare("a7")},
				{From: strToSquare("d4"), To: strToSquare("e5")},
				{From: strToSquare("d4"), To: strToSquare("c3")},
				{From: strToSquare("d4"), To: strToSquare("b2")},
				{From: strToSquare("d4"), To: strToSquare("a1")},
			},
		},
		{
			name: "white blocked enemy pieces",
			pos: &Position{
				WhiteBishops: bitboardFromStrs([]string{"e4"}),
				BlackPawns:   bitboardFromStrs([]string{"c2", "f5"}),
			},
			color: White,
			legalMoves: []Move{
				{From: strToSquare("e4"), To: strToSquare("d5")},
				{From: strToSquare("e4"), To: strToSquare("c6")},
				{From: strToSquare("e4"), To: strToSquare("b7")},
				{From: strToSquare("e4"), To: strToSquare("a8")},
				{From: strToSquare("e4"), To: strToSquare("f5")}, // capture
				{From: strToSquare("e4"), To: strToSquare("d3")},
				{From: strToSquare("e4"), To: strToSquare("c2")}, // capture
				{From: strToSquare("e4"), To: strToSquare("f3")},
				{From: strToSquare("e4"), To: strToSquare("g2")},
				{From: strToSquare("e4"), To: strToSquare("h1")},
			},
		},
		{
			name: "black blocked enemy pieces",
			pos: &Position{
				BlackBishops: bitboardFromStrs([]string{"d4"}),
				WhitePawns:   bitboardFromStrs([]string{"e3", "f6"}),
			},
			color: Black,
			legalMoves: []Move{
				{From: strToSquare("d4"), To: strToSquare("c5")},
				{From: strToSquare("d4"), To: strToSquare("b6")},
				{From: strToSquare("d4"), To: strToSquare("a7")},
				{From: strToSquare("d4"), To: strToSquare("e5")},
				{From: strToSquare("d4"), To: strToSquare("f6")}, // capture
				{From: strToSquare("d4"), To: strToSquare("c3")},
				{From: strToSquare("d4"), To: strToSquare("b2")},
				{From: strToSquare("d4"), To: strToSquare("a1")},
				{From: strToSquare("d4"), To: strToSquare("e3")}, // capture
			},
		},
		{
			name: "white blocked in corner",
			pos: &Position{
				WhiteBishops: bitboardFromStrs([]string{"a1"}),
				WhitePawns:   bitboardFromStrs([]string{"b2"}),
			},
			color:      White,
			legalMoves: []Move{},
		},
		{
			name: "black blocked in corner",
			pos: &Position{
				BlackBishops: bitboardFromStrs([]string{"a8"}),
				BlackPawns:   bitboardFromStrs([]string{"b7"}),
			},
			color:      Black,
			legalMoves: []Move{},
		},
		{
			name: "white corner clear",
			pos: &Position{
				WhiteBishops: bitboardFromStrs([]string{"a1"}),
				WhitePawns:   bitboardFromStrs([]string{"a2", "b1"}),
			},
			color: White,
			legalMoves: []Move{
				{From: strToSquare("a1"), To: strToSquare("b2")},
				{From: strToSquare("a1"), To: strToSquare("c3")},
				{From: strToSquare("a1"), To: strToSquare("d4")},
				{From: strToSquare("a1"), To: strToSquare("e5")},
				{From: strToSquare("a1"), To: strToSquare("f6")},
				{From: strToSquare("a1"), To: strToSquare("g7")},
				{From: strToSquare("a1"), To: strToSquare("h8")},
			},
		},
		{
			name: "black corner clear",
			pos: &Position{
				BlackBishops: bitboardFromStrs([]string{"a8"}),
				BlackPawns:   bitboardFromStrs([]string{"a7", "b8"}),
			},
			color: Black,
			legalMoves: []Move{
				{From: strToSquare("a8"), To: strToSquare("b7")},
				{From: strToSquare("a8"), To: strToSquare("c6")},
				{From: strToSquare("a8"), To: strToSquare("d5")},
				{From: strToSquare("a8"), To: strToSquare("e4")},
				{From: strToSquare("a8"), To: strToSquare("f3")},
				{From: strToSquare("a8"), To: strToSquare("g2")},
				{From: strToSquare("a8"), To: strToSquare("h1")},
			},
		},
		{
			name: "white mutiple bishops same color",
			pos: &Position{
				WhiteBishops: bitboardFromStrs([]string{"e4", "d5"}),
			},
			color: White,
			legalMoves: []Move{
				{From: strToSquare("e4"), To: strToSquare("f5")},
				{From: strToSquare("e4"), To: strToSquare("g6")},
				{From: strToSquare("e4"), To: strToSquare("h7")},
				{From: strToSquare("e4"), To: strToSquare("d3")},
				{From: strToSquare("e4"), To: strToSquare("c2")},
				{From: strToSquare("e4"), To: strToSquare("b1")},
				{From: strToSquare("e4"), To: strToSquare("f3")},
				{From: strToSquare("e4"), To: strToSquare("g2")},
				{From: strToSquare("e4"), To: strToSquare("h1")},
				{From: strToSquare("d5"), To: strToSquare("c6")},
				{From: strToSquare("d5"), To: strToSquare("b7")},
				{From: strToSquare("d5"), To: strToSquare("a8")},
				{From: strToSquare("d5"), To: strToSquare("e6")},
				{From: strToSquare("d5"), To: strToSquare("f7")},
				{From: strToSquare("d5"), To: strToSquare("g8")},
				{From: strToSquare("d5"), To: strToSquare("c4")},
				{From: strToSquare("d5"), To: strToSquare("b3")},
				{From: strToSquare("d5"), To: strToSquare("a2")},
			},
		},
		{
			name: "black multiple bishops same color",
			pos: &Position{
				BlackBishops: bitboardFromStrs([]string{"d4", "c5"}),
			},
			color: Black,
			legalMoves: []Move{
				{From: strToSquare("d4"), To: strToSquare("e5")},
				{From: strToSquare("d4"), To: strToSquare("f6")},
				{From: strToSquare("d4"), To: strToSquare("g7")},
				{From: strToSquare("d4"), To: strToSquare("h8")},
				{From: strToSquare("d4"), To: strToSquare("c3")},
				{From: strToSquare("d4"), To: strToSquare("b2")},
				{From: strToSquare("d4"), To: strToSquare("a1")},
				{From: strToSquare("d4"), To: strToSquare("e3")},
				{From: strToSquare("d4"), To: strToSquare("f2")},
				{From: strToSquare("d4"), To: strToSquare("g1")},
				{From: strToSquare("c5"), To: strToSquare("b6")},
				{From: strToSquare("c5"), To: strToSquare("a7")},
				{From: strToSquare("c5"), To: strToSquare("b4")},
				{From: strToSquare("c5"), To: strToSquare("a3")},
				{From: strToSquare("c5"), To: strToSquare("d6")},
				{From: strToSquare("c5"), To: strToSquare("e7")},
				{From: strToSquare("c5"), To: strToSquare("f8")},
			},
		},
		{
			name: "white mutiple bishops different colors",
			pos: &Position{
				WhiteBishops: bitboardFromStrs([]string{"e4", "g1"}),
			},
			color: White,
			legalMoves: []Move{
				{From: strToSquare("e4"), To: strToSquare("d5")},
				{From: strToSquare("e4"), To: strToSquare("c6")},
				{From: strToSquare("e4"), To: strToSquare("b7")},
				{From: strToSquare("e4"), To: strToSquare("a8")},
				{From: strToSquare("e4"), To: strToSquare("f5")},
				{From: strToSquare("e4"), To: strToSquare("g6")},
				{From: strToSquare("e4"), To: strToSquare("h7")},
				{From: strToSquare("e4"), To: strToSquare("d3")},
				{From: strToSquare("e4"), To: strToSquare("c2")},
				{From: strToSquare("e4"), To: strToSquare("b1")},
				{From: strToSquare("e4"), To: strToSquare("f3")},
				{From: strToSquare("e4"), To: strToSquare("g2")},
				{From: strToSquare("e4"), To: strToSquare("h1")},
				{From: strToSquare("g1"), To: strToSquare("h2")},
				{From: strToSquare("g1"), To: strToSquare("f2")},
				{From: strToSquare("g1"), To: strToSquare("e3")},
				{From: strToSquare("g1"), To: strToSquare("d4")},
				{From: strToSquare("g1"), To: strToSquare("c5")},
				{From: strToSquare("g1"), To: strToSquare("b6")},
				{From: strToSquare("g1"), To: strToSquare("a7")},
			},
		},
		{
			name: "black multiple bishops different colors",
			pos: &Position{
				BlackBishops: bitboardFromStrs([]string{"d4", "d1"}),
			},
			color: Black,
			legalMoves: []Move{
				{From: strToSquare("d4"), To: strToSquare("c5")},
				{From: strToSquare("d4"), To: strToSquare("b6")},
				{From: strToSquare("d4"), To: strToSquare("a7")},
				{From: strToSquare("d4"), To: strToSquare("e5")},
				{From: strToSquare("d4"), To: strToSquare("f6")},
				{From: strToSquare("d4"), To: strToSquare("g7")},
				{From: strToSquare("d4"), To: strToSquare("h8")},
				{From: strToSquare("d4"), To: strToSquare("c3")},
				{From: strToSquare("d4"), To: strToSquare("b2")},
				{From: strToSquare("d4"), To: strToSquare("a1")},
				{From: strToSquare("d4"), To: strToSquare("e3")},
				{From: strToSquare("d4"), To: strToSquare("f2")},
				{From: strToSquare("d4"), To: strToSquare("g1")},
				{From: strToSquare("d1"), To: strToSquare("e2")},
				{From: strToSquare("d1"), To: strToSquare("f3")},
				{From: strToSquare("d1"), To: strToSquare("g4")},
				{From: strToSquare("d1"), To: strToSquare("h5")},
				{From: strToSquare("d1"), To: strToSquare("c2")},
				{From: strToSquare("d1"), To: strToSquare("b3")},
				{From: strToSquare("d1"), To: strToSquare("a4")},
			},
		},
	}

	g := NewGenerator()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			moves := g.generateBishopMoves(tt.pos, tt.color)
			assert.ElementsMatch(t, tt.legalMoves, moves)
		})
	}
}

func TestGenerateRookMoves(t *testing.T) {
	tests := []struct {
		name       string
		pos        *Position
		color      Color
		legalMoves []Move
	}{
		{
			name: "white clear board",
			pos: &Position{
				WhiteRooks: bitboardFromStrs([]string{"e4"}),
			},
			color: White,
			legalMoves: []Move{
				{From: strToSquare("e4"), To: strToSquare("a4")},
				{From: strToSquare("e4"), To: strToSquare("b4")},
				{From: strToSquare("e4"), To: strToSquare("c4")},
				{From: strToSquare("e4"), To: strToSquare("d4")},
				{From: strToSquare("e4"), To: strToSquare("f4")},
				{From: strToSquare("e4"), To: strToSquare("g4")},
				{From: strToSquare("e4"), To: strToSquare("h4")},
				{From: strToSquare("e4"), To: strToSquare("e1")},
				{From: strToSquare("e4"), To: strToSquare("e2")},
				{From: strToSquare("e4"), To: strToSquare("e3")},
				{From: strToSquare("e4"), To: strToSquare("e5")},
				{From: strToSquare("e4"), To: strToSquare("e6")},
				{From: strToSquare("e4"), To: strToSquare("e7")},
				{From: strToSquare("e4"), To: strToSquare("e8")},
			},
		},
		{
			name: "black clear board",
			pos: &Position{
				BlackRooks: bitboardFromStrs([]string{"d4"}),
			},
			color: Black,
			legalMoves: []Move{
				{From: strToSquare("d4"), To: strToSquare("a4")},
				{From: strToSquare("d4"), To: strToSquare("b4")},
				{From: strToSquare("d4"), To: strToSquare("c4")},
				{From: strToSquare("d4"), To: strToSquare("e4")},
				{From: strToSquare("d4"), To: strToSquare("f4")},
				{From: strToSquare("d4"), To: strToSquare("g4")},
				{From: strToSquare("d4"), To: strToSquare("h4")},
				{From: strToSquare("d4"), To: strToSquare("d1")},
				{From: strToSquare("d4"), To: strToSquare("d2")},
				{From: strToSquare("d4"), To: strToSquare("d3")},
				{From: strToSquare("d4"), To: strToSquare("d5")},
				{From: strToSquare("d4"), To: strToSquare("d6")},
				{From: strToSquare("d4"), To: strToSquare("d7")},
				{From: strToSquare("d4"), To: strToSquare("d8")},
			},
		},
		{
			name: "white blocked own pieces",
			pos: &Position{
				WhiteRooks: bitboardFromStrs([]string{"e4"}),
				WhitePawns: bitboardFromStrs([]string{"e5", "f4"}),
			},
			color: White,
			legalMoves: []Move{
				{From: strToSquare("e4"), To: strToSquare("a4")},
				{From: strToSquare("e4"), To: strToSquare("b4")},
				{From: strToSquare("e4"), To: strToSquare("c4")},
				{From: strToSquare("e4"), To: strToSquare("d4")},
				{From: strToSquare("e4"), To: strToSquare("e1")},
				{From: strToSquare("e4"), To: strToSquare("e2")},
				{From: strToSquare("e4"), To: strToSquare("e3")},
			},
		},
		{
			name: "black blocked own pieces",
			pos: &Position{
				BlackRooks: bitboardFromStrs([]string{"d4"}),
				BlackPawns: bitboardFromStrs([]string{"d3", "c4"}),
			},
			color: Black,
			legalMoves: []Move{
				{From: strToSquare("d4"), To: strToSquare("e4")},
				{From: strToSquare("d4"), To: strToSquare("f4")},
				{From: strToSquare("d4"), To: strToSquare("g4")},
				{From: strToSquare("d4"), To: strToSquare("h4")},
				{From: strToSquare("d4"), To: strToSquare("d5")},
				{From: strToSquare("d4"), To: strToSquare("d6")},
				{From: strToSquare("d4"), To: strToSquare("d7")},
				{From: strToSquare("d4"), To: strToSquare("d8")},
			},
		},
		{
			name: "white blocked enemy pieces",
			pos: &Position{
				WhiteRooks: bitboardFromStrs([]string{"e4"}),
				BlackPawns: bitboardFromStrs([]string{"e5", "f4"}),
			},
			color: White,
			legalMoves: []Move{
				{From: strToSquare("e4"), To: strToSquare("a4")},
				{From: strToSquare("e4"), To: strToSquare("b4")},
				{From: strToSquare("e4"), To: strToSquare("c4")},
				{From: strToSquare("e4"), To: strToSquare("d4")},
				{From: strToSquare("e4"), To: strToSquare("f4")}, // capture
				{From: strToSquare("e4"), To: strToSquare("e1")},
				{From: strToSquare("e4"), To: strToSquare("e2")},
				{From: strToSquare("e4"), To: strToSquare("e3")},
				{From: strToSquare("e4"), To: strToSquare("e5")}, // capture
			},
		},
		{
			name: "black blocked enemy pieces",
			pos: &Position{
				BlackRooks: bitboardFromStrs([]string{"d4"}),
				WhitePawns: bitboardFromStrs([]string{"d3", "c4"}),
			},
			color: Black,
			legalMoves: []Move{
				{From: strToSquare("d4"), To: strToSquare("c4")}, // capture
				{From: strToSquare("d4"), To: strToSquare("e4")},
				{From: strToSquare("d4"), To: strToSquare("f4")},
				{From: strToSquare("d4"), To: strToSquare("g4")},
				{From: strToSquare("d4"), To: strToSquare("h4")},
				{From: strToSquare("d4"), To: strToSquare("d3")}, // capture
				{From: strToSquare("d4"), To: strToSquare("d5")},
				{From: strToSquare("d4"), To: strToSquare("d6")},
				{From: strToSquare("d4"), To: strToSquare("d7")},
				{From: strToSquare("d4"), To: strToSquare("d8")},
			},
		},
		{
			name: "white blocked in corner",
			pos: &Position{
				WhiteRooks: bitboardFromStrs([]string{"a1"}),
				WhitePawns: bitboardFromStrs([]string{"a2", "b1"}),
			},
			color:      White,
			legalMoves: []Move{},
		},
		{
			name: "black blocked in corner",
			pos: &Position{
				BlackRooks: bitboardFromStrs([]string{"a8"}),
				BlackPawns: bitboardFromStrs([]string{"a7", "b8"}),
			},
			color:      Black,
			legalMoves: []Move{},
		},
		{
			name: "white corner clear",
			pos: &Position{
				WhiteRooks: bitboardFromStrs([]string{"a1"}),
				WhitePawns: bitboardFromStrs([]string{"b2"}),
			},
			color: White,
			legalMoves: []Move{
				{From: strToSquare("a1"), To: strToSquare("b1")},
				{From: strToSquare("a1"), To: strToSquare("c1")},
				{From: strToSquare("a1"), To: strToSquare("d1")},
				{From: strToSquare("a1"), To: strToSquare("e1")},
				{From: strToSquare("a1"), To: strToSquare("f1")},
				{From: strToSquare("a1"), To: strToSquare("g1")},
				{From: strToSquare("a1"), To: strToSquare("h1")},
				{From: strToSquare("a1"), To: strToSquare("a2")},
				{From: strToSquare("a1"), To: strToSquare("a3")},
				{From: strToSquare("a1"), To: strToSquare("a4")},
				{From: strToSquare("a1"), To: strToSquare("a5")},
				{From: strToSquare("a1"), To: strToSquare("a6")},
				{From: strToSquare("a1"), To: strToSquare("a7")},
				{From: strToSquare("a1"), To: strToSquare("a8")},
			},
		},
		{
			name: "black corner clear",
			pos: &Position{
				BlackRooks: bitboardFromStrs([]string{"a8"}),
				BlackPawns: bitboardFromStrs([]string{"b7"}),
			},
			color: Black,
			legalMoves: []Move{
				{From: strToSquare("a8"), To: strToSquare("b8")},
				{From: strToSquare("a8"), To: strToSquare("c8")},
				{From: strToSquare("a8"), To: strToSquare("d8")},
				{From: strToSquare("a8"), To: strToSquare("e8")},
				{From: strToSquare("a8"), To: strToSquare("f8")},
				{From: strToSquare("a8"), To: strToSquare("g8")},
				{From: strToSquare("a8"), To: strToSquare("h8")},
				{From: strToSquare("a8"), To: strToSquare("a7")},
				{From: strToSquare("a8"), To: strToSquare("a6")},
				{From: strToSquare("a8"), To: strToSquare("a5")},
				{From: strToSquare("a8"), To: strToSquare("a4")},
				{From: strToSquare("a8"), To: strToSquare("a3")},
				{From: strToSquare("a8"), To: strToSquare("a2")},
				{From: strToSquare("a8"), To: strToSquare("a1")},
			},
		},
		{
			name: "white multiple rooks",
			pos: &Position{
				WhiteRooks: bitboardFromStrs([]string{"e4", "e5"}),
			},
			color: White,
			legalMoves: []Move{
				{From: strToSquare("e4"), To: strToSquare("a4")},
				{From: strToSquare("e4"), To: strToSquare("b4")},
				{From: strToSquare("e4"), To: strToSquare("c4")},
				{From: strToSquare("e4"), To: strToSquare("d4")},
				{From: strToSquare("e4"), To: strToSquare("f4")},
				{From: strToSquare("e4"), To: strToSquare("g4")},
				{From: strToSquare("e4"), To: strToSquare("h4")},
				{From: strToSquare("e4"), To: strToSquare("e1")},
				{From: strToSquare("e4"), To: strToSquare("e2")},
				{From: strToSquare("e4"), To: strToSquare("e3")},
				{From: strToSquare("e5"), To: strToSquare("e6")},
				{From: strToSquare("e5"), To: strToSquare("e7")},
				{From: strToSquare("e5"), To: strToSquare("e8")},
				{From: strToSquare("e5"), To: strToSquare("a5")},
				{From: strToSquare("e5"), To: strToSquare("b5")},
				{From: strToSquare("e5"), To: strToSquare("c5")},
				{From: strToSquare("e5"), To: strToSquare("d5")},
				{From: strToSquare("e5"), To: strToSquare("f5")},
				{From: strToSquare("e5"), To: strToSquare("g5")},
				{From: strToSquare("e5"), To: strToSquare("h5")},
			},
		},
		{
			name: "black multiple rooks",
			pos: &Position{
				BlackRooks: bitboardFromStrs([]string{"d4", "f4"}),
			},
			color: Black,
			legalMoves: []Move{
				{From: strToSquare("d4"), To: strToSquare("a4")},
				{From: strToSquare("d4"), To: strToSquare("b4")},
				{From: strToSquare("d4"), To: strToSquare("c4")},
				{From: strToSquare("d4"), To: strToSquare("e4")},
				{From: strToSquare("d4"), To: strToSquare("d1")},
				{From: strToSquare("d4"), To: strToSquare("d2")},
				{From: strToSquare("d4"), To: strToSquare("d3")},
				{From: strToSquare("d4"), To: strToSquare("d5")},
				{From: strToSquare("d4"), To: strToSquare("d6")},
				{From: strToSquare("d4"), To: strToSquare("d7")},
				{From: strToSquare("d4"), To: strToSquare("d8")},
				{From: strToSquare("f4"), To: strToSquare("e4")},
				{From: strToSquare("f4"), To: strToSquare("g4")},
				{From: strToSquare("f4"), To: strToSquare("h4")},
				{From: strToSquare("f4"), To: strToSquare("f1")},
				{From: strToSquare("f4"), To: strToSquare("f2")},
				{From: strToSquare("f4"), To: strToSquare("f3")},
				{From: strToSquare("f4"), To: strToSquare("f5")},
				{From: strToSquare("f4"), To: strToSquare("f6")},
				{From: strToSquare("f4"), To: strToSquare("f7")},
				{From: strToSquare("f4"), To: strToSquare("f8")},
			},
		},
	}

	g := NewGenerator()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			moves := g.generateRookMoves(tt.pos, tt.color)
			assert.ElementsMatch(t, tt.legalMoves, moves)
		})
	}
}

func TestGenerateQueenMoves(t *testing.T) {
	tests := []struct {
		name       string
		pos        *Position
		color      Color
		legalMoves []Move
	}{
		{
			name: "white clear board",
			pos: &Position{
				WhiteQueens: bitboardFromStrs([]string{"e4"}),
			},
			color: White,
			legalMoves: []Move{
				{From: strToSquare("e4"), To: strToSquare("a4")},
				{From: strToSquare("e4"), To: strToSquare("b4")},
				{From: strToSquare("e4"), To: strToSquare("c4")},
				{From: strToSquare("e4"), To: strToSquare("d4")},
				{From: strToSquare("e4"), To: strToSquare("f4")},
				{From: strToSquare("e4"), To: strToSquare("g4")},
				{From: strToSquare("e4"), To: strToSquare("h4")},
				{From: strToSquare("e4"), To: strToSquare("e1")},
				{From: strToSquare("e4"), To: strToSquare("e2")},
				{From: strToSquare("e4"), To: strToSquare("e3")},
				{From: strToSquare("e4"), To: strToSquare("e5")},
				{From: strToSquare("e4"), To: strToSquare("e6")},
				{From: strToSquare("e4"), To: strToSquare("e7")},
				{From: strToSquare("e4"), To: strToSquare("e8")},
				{From: strToSquare("e4"), To: strToSquare("d5")},
				{From: strToSquare("e4"), To: strToSquare("c6")},
				{From: strToSquare("e4"), To: strToSquare("b7")},
				{From: strToSquare("e4"), To: strToSquare("a8")},
				{From: strToSquare("e4"), To: strToSquare("f5")},
				{From: strToSquare("e4"), To: strToSquare("g6")},
				{From: strToSquare("e4"), To: strToSquare("h7")},
				{From: strToSquare("e4"), To: strToSquare("d3")},
				{From: strToSquare("e4"), To: strToSquare("c2")},
				{From: strToSquare("e4"), To: strToSquare("b1")},
				{From: strToSquare("e4"), To: strToSquare("f3")},
				{From: strToSquare("e4"), To: strToSquare("g2")},
				{From: strToSquare("e4"), To: strToSquare("h1")},
			},
		},
		{
			name: "black clear board",
			pos: &Position{
				BlackQueens: bitboardFromStrs([]string{"d4"}),
			},
			color: Black,
			legalMoves: []Move{
				{From: strToSquare("d4"), To: strToSquare("a4")},
				{From: strToSquare("d4"), To: strToSquare("b4")},
				{From: strToSquare("d4"), To: strToSquare("c4")},
				{From: strToSquare("d4"), To: strToSquare("e4")},
				{From: strToSquare("d4"), To: strToSquare("f4")},
				{From: strToSquare("d4"), To: strToSquare("g4")},
				{From: strToSquare("d4"), To: strToSquare("h4")},
				{From: strToSquare("d4"), To: strToSquare("d1")},
				{From: strToSquare("d4"), To: strToSquare("d2")},
				{From: strToSquare("d4"), To: strToSquare("d3")},
				{From: strToSquare("d4"), To: strToSquare("d5")},
				{From: strToSquare("d4"), To: strToSquare("d6")},
				{From: strToSquare("d4"), To: strToSquare("d7")},
				{From: strToSquare("d4"), To: strToSquare("d8")},
				{From: strToSquare("d4"), To: strToSquare("c5")},
				{From: strToSquare("d4"), To: strToSquare("b6")},
				{From: strToSquare("d4"), To: strToSquare("a7")},
				{From: strToSquare("d4"), To: strToSquare("e5")},
				{From: strToSquare("d4"), To: strToSquare("f6")},
				{From: strToSquare("d4"), To: strToSquare("g7")},
				{From: strToSquare("d4"), To: strToSquare("h8")},
				{From: strToSquare("d4"), To: strToSquare("c3")},
				{From: strToSquare("d4"), To: strToSquare("b2")},
				{From: strToSquare("d4"), To: strToSquare("a1")},
				{From: strToSquare("d4"), To: strToSquare("e3")},
				{From: strToSquare("d4"), To: strToSquare("f2")},
				{From: strToSquare("d4"), To: strToSquare("g1")},
			},
		},
		{
			name: "white blocked own pieces",
			pos: &Position{
				WhiteQueens: bitboardFromStrs([]string{"e4"}),
				WhitePawns:  bitboardFromStrs([]string{"e5", "f4", "c2", "f5"}),
			},
			color: White,
			legalMoves: []Move{
				{From: strToSquare("e4"), To: strToSquare("a4")},
				{From: strToSquare("e4"), To: strToSquare("b4")},
				{From: strToSquare("e4"), To: strToSquare("c4")},
				{From: strToSquare("e4"), To: strToSquare("d4")},
				{From: strToSquare("e4"), To: strToSquare("e1")},
				{From: strToSquare("e4"), To: strToSquare("e2")},
				{From: strToSquare("e4"), To: strToSquare("e3")},
				{From: strToSquare("e4"), To: strToSquare("d5")},
				{From: strToSquare("e4"), To: strToSquare("c6")},
				{From: strToSquare("e4"), To: strToSquare("b7")},
				{From: strToSquare("e4"), To: strToSquare("a8")},
				{From: strToSquare("e4"), To: strToSquare("d3")},
				{From: strToSquare("e4"), To: strToSquare("f3")},
				{From: strToSquare("e4"), To: strToSquare("g2")},
				{From: strToSquare("e4"), To: strToSquare("h1")},
			},
		},
		{
			name: "black blocked own pieces",
			pos: &Position{
				BlackQueens: bitboardFromStrs([]string{"d4"}),
				BlackPawns:  bitboardFromStrs([]string{"d3", "c4", "e3", "f6"}),
			},
			color: Black,
			legalMoves: []Move{
				{From: strToSquare("d4"), To: strToSquare("e4")},
				{From: strToSquare("d4"), To: strToSquare("f4")},
				{From: strToSquare("d4"), To: strToSquare("g4")},
				{From: strToSquare("d4"), To: strToSquare("h4")},
				{From: strToSquare("d4"), To: strToSquare("d5")},
				{From: strToSquare("d4"), To: strToSquare("d6")},
				{From: strToSquare("d4"), To: strToSquare("d7")},
				{From: strToSquare("d4"), To: strToSquare("d8")},
				{From: strToSquare("d4"), To: strToSquare("c5")},
				{From: strToSquare("d4"), To: strToSquare("b6")},
				{From: strToSquare("d4"), To: strToSquare("a7")},
				{From: strToSquare("d4"), To: strToSquare("e5")},
				{From: strToSquare("d4"), To: strToSquare("c3")},
				{From: strToSquare("d4"), To: strToSquare("b2")},
				{From: strToSquare("d4"), To: strToSquare("a1")},
			},
		},
		{
			name: "white blocked enemy pieces",
			pos: &Position{
				WhiteQueens: bitboardFromStrs([]string{"e4"}),
				BlackPawns:  bitboardFromStrs([]string{"e5", "f4", "c2", "f5"}),
			},
			color: White,
			legalMoves: []Move{
				{From: strToSquare("e4"), To: strToSquare("a4")},
				{From: strToSquare("e4"), To: strToSquare("b4")},
				{From: strToSquare("e4"), To: strToSquare("c4")},
				{From: strToSquare("e4"), To: strToSquare("d4")},
				{From: strToSquare("e4"), To: strToSquare("f4")}, // capture
				{From: strToSquare("e4"), To: strToSquare("e1")},
				{From: strToSquare("e4"), To: strToSquare("e2")},
				{From: strToSquare("e4"), To: strToSquare("e3")},
				{From: strToSquare("e4"), To: strToSquare("e5")}, // capture
				{From: strToSquare("e4"), To: strToSquare("d5")},
				{From: strToSquare("e4"), To: strToSquare("c6")},
				{From: strToSquare("e4"), To: strToSquare("b7")},
				{From: strToSquare("e4"), To: strToSquare("a8")},
				{From: strToSquare("e4"), To: strToSquare("f5")}, // capture
				{From: strToSquare("e4"), To: strToSquare("d3")},
				{From: strToSquare("e4"), To: strToSquare("c2")}, // capture
				{From: strToSquare("e4"), To: strToSquare("f3")},
				{From: strToSquare("e4"), To: strToSquare("g2")},
				{From: strToSquare("e4"), To: strToSquare("h1")},
			},
		},
		{
			name: "black blocked enemy pieces",
			pos: &Position{
				BlackQueens: bitboardFromStrs([]string{"d4"}),
				WhitePawns:  bitboardFromStrs([]string{"d3", "c4", "e3", "f6"}),
			},
			color: Black,
			legalMoves: []Move{
				{From: strToSquare("d4"), To: strToSquare("c4")}, // capture
				{From: strToSquare("d4"), To: strToSquare("e4")},
				{From: strToSquare("d4"), To: strToSquare("f4")},
				{From: strToSquare("d4"), To: strToSquare("g4")},
				{From: strToSquare("d4"), To: strToSquare("h4")},
				{From: strToSquare("d4"), To: strToSquare("d3")}, // capture
				{From: strToSquare("d4"), To: strToSquare("d5")},
				{From: strToSquare("d4"), To: strToSquare("d6")},
				{From: strToSquare("d4"), To: strToSquare("d7")},
				{From: strToSquare("d4"), To: strToSquare("d8")},
				{From: strToSquare("d4"), To: strToSquare("c5")},
				{From: strToSquare("d4"), To: strToSquare("b6")},
				{From: strToSquare("d4"), To: strToSquare("a7")},
				{From: strToSquare("d4"), To: strToSquare("e5")},
				{From: strToSquare("d4"), To: strToSquare("f6")}, // capture
				{From: strToSquare("d4"), To: strToSquare("c3")},
				{From: strToSquare("d4"), To: strToSquare("b2")},
				{From: strToSquare("d4"), To: strToSquare("a1")},
				{From: strToSquare("d4"), To: strToSquare("e3")}, // capture
			},
		},
		{
			name: "white blocked in corner",
			pos: &Position{
				WhiteQueens: bitboardFromStrs([]string{"a1"}),
				WhitePawns:  bitboardFromStrs([]string{"a2", "b1", "b2"}),
			},
			color:      White,
			legalMoves: []Move{},
		},
		{
			name: "black blocked in corner",
			pos: &Position{
				BlackQueens: bitboardFromStrs([]string{"a8"}),
				BlackPawns:  bitboardFromStrs([]string{"a7", "b7", "b8"}),
			},
			color:      Black,
			legalMoves: []Move{},
		},
		{
			name: "white multiple queens blocked own pieces",
			pos: &Position{
				WhiteQueens: bitboardFromStrs([]string{"e4", "g2"}),
				WhitePawns:  bitboardFromStrs([]string{"e5", "f4", "c2", "f5", "g5"}),
			},
			color: White,
			legalMoves: []Move{
				{From: strToSquare("e4"), To: strToSquare("a4")},
				{From: strToSquare("e4"), To: strToSquare("b4")},
				{From: strToSquare("e4"), To: strToSquare("c4")},
				{From: strToSquare("e4"), To: strToSquare("d4")},
				{From: strToSquare("e4"), To: strToSquare("e1")},
				{From: strToSquare("e4"), To: strToSquare("e2")},
				{From: strToSquare("e4"), To: strToSquare("e3")},
				{From: strToSquare("e4"), To: strToSquare("d5")},
				{From: strToSquare("e4"), To: strToSquare("c6")},
				{From: strToSquare("e4"), To: strToSquare("b7")},
				{From: strToSquare("e4"), To: strToSquare("a8")},
				{From: strToSquare("e4"), To: strToSquare("d3")},
				{From: strToSquare("e4"), To: strToSquare("f3")},
				{From: strToSquare("g2"), To: strToSquare("d2")},
				{From: strToSquare("g2"), To: strToSquare("e2")},
				{From: strToSquare("g2"), To: strToSquare("f2")},
				{From: strToSquare("g2"), To: strToSquare("h2")},
				{From: strToSquare("g2"), To: strToSquare("f1")},
				{From: strToSquare("g2"), To: strToSquare("g1")},
				{From: strToSquare("g2"), To: strToSquare("h1")},
				{From: strToSquare("g2"), To: strToSquare("f3")},
				{From: strToSquare("g2"), To: strToSquare("g3")},
				{From: strToSquare("g2"), To: strToSquare("h3")},
				{From: strToSquare("g2"), To: strToSquare("g4")},
			},
		},
		{
			name: "black multiple queens blocked own pieces",
			pos: &Position{
				BlackQueens: bitboardFromStrs([]string{"d4", "c3"}),
				BlackPawns:  bitboardFromStrs([]string{"d3", "c4", "e3", "f6"}),
			},
			color: Black,
			legalMoves: []Move{
				{From: strToSquare("d4"), To: strToSquare("e4")},
				{From: strToSquare("d4"), To: strToSquare("f4")},
				{From: strToSquare("d4"), To: strToSquare("g4")},
				{From: strToSquare("d4"), To: strToSquare("h4")},
				{From: strToSquare("d4"), To: strToSquare("d5")},
				{From: strToSquare("d4"), To: strToSquare("d6")},
				{From: strToSquare("d4"), To: strToSquare("d7")},
				{From: strToSquare("d4"), To: strToSquare("d8")},
				{From: strToSquare("d4"), To: strToSquare("c5")},
				{From: strToSquare("d4"), To: strToSquare("b6")},
				{From: strToSquare("d4"), To: strToSquare("a7")},
				{From: strToSquare("d4"), To: strToSquare("e5")},
				{From: strToSquare("c3"), To: strToSquare("b2")},
				{From: strToSquare("c3"), To: strToSquare("a1")},
				{From: strToSquare("c3"), To: strToSquare("b3")},
				{From: strToSquare("c3"), To: strToSquare("a3")},
				{From: strToSquare("c3"), To: strToSquare("b4")},
				{From: strToSquare("c3"), To: strToSquare("a5")},
				{From: strToSquare("c3"), To: strToSquare("d2")},
				{From: strToSquare("c3"), To: strToSquare("e1")},
				{From: strToSquare("c3"), To: strToSquare("c2")},
				{From: strToSquare("c3"), To: strToSquare("c1")},
			},
		},
	}

	g := NewGenerator()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			moves := g.generateQueenMoves(tt.pos, tt.color)
			assert.ElementsMatch(t, tt.legalMoves, moves)
		})
	}
}

func TestGenerateKingMoves(t *testing.T) {
	tests := []struct {
		name       string
		pos        *Position
		color      Color
		legalMoves []Move
	}{
		{
			name: "white clear board",
			pos: &Position{
				WhiteKing: bitboardFromStrs([]string{"e4"}),
			},
			color: White,
			legalMoves: []Move{
				{From: strToSquare("e4"), To: strToSquare("e5")},
				{From: strToSquare("e4"), To: strToSquare("f5")},
				{From: strToSquare("e4"), To: strToSquare("f4")},
				{From: strToSquare("e4"), To: strToSquare("f3")},
				{From: strToSquare("e4"), To: strToSquare("e3")},
				{From: strToSquare("e4"), To: strToSquare("d3")},
				{From: strToSquare("e4"), To: strToSquare("d4")},
				{From: strToSquare("e4"), To: strToSquare("d5")},
			},
		},
		{
			name: "black clear board",
			pos: &Position{
				BlackKing: bitboardFromStrs([]string{"d4"}),
			},
			color: Black,
			legalMoves: []Move{
				{From: strToSquare("d4"), To: strToSquare("d5")},
				{From: strToSquare("d4"), To: strToSquare("e5")},
				{From: strToSquare("d4"), To: strToSquare("e4")},
				{From: strToSquare("d4"), To: strToSquare("e3")},
				{From: strToSquare("d4"), To: strToSquare("d3")},
				{From: strToSquare("d4"), To: strToSquare("c3")},
				{From: strToSquare("d4"), To: strToSquare("c4")},
				{From: strToSquare("d4"), To: strToSquare("c5")},
			},
		},
		{
			name: "white blocked own pieces",
			pos: &Position{
				WhiteKing:  bitboardFromStrs([]string{"e4"}),
				WhitePawns: bitboardFromStrs([]string{"e5", "f4", "f5"}),
			},
			color: White,
			legalMoves: []Move{
				{From: strToSquare("e4"), To: strToSquare("f3")},
				{From: strToSquare("e4"), To: strToSquare("e3")},
				{From: strToSquare("e4"), To: strToSquare("d3")},
				{From: strToSquare("e4"), To: strToSquare("d4")},
				{From: strToSquare("e4"), To: strToSquare("d5")},
			},
		},
		{
			name: "black blocked own pieces",
			pos: &Position{
				BlackKing:  bitboardFromStrs([]string{"d4"}),
				BlackPawns: bitboardFromStrs([]string{"d3", "c4", "e3"}),
			},
			color: Black,
			legalMoves: []Move{
				{From: strToSquare("d4"), To: strToSquare("d5")},
				{From: strToSquare("d4"), To: strToSquare("e5")},
				{From: strToSquare("d4"), To: strToSquare("e4")},
				{From: strToSquare("d4"), To: strToSquare("c3")},
				{From: strToSquare("d4"), To: strToSquare("c5")},
			},
		},
		{
			name: "white blocked enemy pieces",
			pos: &Position{
				WhiteKing:  bitboardFromStrs([]string{"e4"}),
				BlackPawns: bitboardFromStrs([]string{"e5", "f4", "f5"}),
			},
			color: White,
			legalMoves: []Move{
				{From: strToSquare("e4"), To: strToSquare("e5")}, // capture
				{From: strToSquare("e4"), To: strToSquare("f5")}, // capture
				{From: strToSquare("e4"), To: strToSquare("f4")}, // capture
				{From: strToSquare("e4"), To: strToSquare("f3")},
				{From: strToSquare("e4"), To: strToSquare("e3")},
				{From: strToSquare("e4"), To: strToSquare("d3")},
				{From: strToSquare("e4"), To: strToSquare("d4")},
				{From: strToSquare("e4"), To: strToSquare("d5")},
			},
		},
		{
			name: "black blocked enemy pieces",
			pos: &Position{
				BlackKing:  bitboardFromStrs([]string{"d4"}),
				WhitePawns: bitboardFromStrs([]string{"d3", "c4", "e3"}),
			},
			color: Black,
			legalMoves: []Move{
				{From: strToSquare("d4"), To: strToSquare("d5")},
				{From: strToSquare("d4"), To: strToSquare("e5")},
				{From: strToSquare("d4"), To: strToSquare("e4")},
				{From: strToSquare("d4"), To: strToSquare("e3")}, // capture
				{From: strToSquare("d4"), To: strToSquare("d3")}, // capture
				{From: strToSquare("d4"), To: strToSquare("c3")},
				{From: strToSquare("d4"), To: strToSquare("c4")}, // capture
				{From: strToSquare("d4"), To: strToSquare("c5")},
			},
		},
		{
			name: "white blocked in corner",
			pos: &Position{
				WhiteKing:  bitboardFromStrs([]string{"a1"}),
				WhitePawns: bitboardFromStrs([]string{"a2", "b1", "b2"}),
			},
			color:      White,
			legalMoves: []Move{},
		},
		{
			name: "black blocked in corner",
			pos: &Position{
				BlackKing:  bitboardFromStrs([]string{"a8"}),
				BlackPawns: bitboardFromStrs([]string{"a7", "b7", "b8"}),
			},
			color:      Black,
			legalMoves: []Move{},
		},
		{
			name: "white kingside castling allowed",
			pos: &Position{
				WhiteKing:  bitboardFromStrs([]string{"e1"}),
				WhiteRooks: bitboardFromStrs([]string{"h1"}),
				CastlingRights: CastlingRights{
					WhiteOO: true,
				},
			},
			color: White,
			legalMoves: []Move{
				{From: strToSquare("e1"), To: strToSquare("d1")},
				{From: strToSquare("e1"), To: strToSquare("d2")},
				{From: strToSquare("e1"), To: strToSquare("e2")},
				{From: strToSquare("e1"), To: strToSquare("f2")},
				{From: strToSquare("e1"), To: strToSquare("f1")},
				{From: strToSquare("e1"), To: strToSquare("g1")}, // castling
			},
		},
		{
			name: "white queenside castling allowed",
			pos: &Position{
				WhiteKing:  bitboardFromStrs([]string{"e1"}),
				WhiteRooks: bitboardFromStrs([]string{"a1"}),
				CastlingRights: CastlingRights{
					WhiteOOO: true,
				},
			},
			color: White,
			legalMoves: []Move{
				{From: strToSquare("e1"), To: strToSquare("d1")},
				{From: strToSquare("e1"), To: strToSquare("d2")},
				{From: strToSquare("e1"), To: strToSquare("e2")},
				{From: strToSquare("e1"), To: strToSquare("f2")},
				{From: strToSquare("e1"), To: strToSquare("f1")},
				{From: strToSquare("e1"), To: strToSquare("c1")}, // castling
			},
		},
		{
			name: "black kingside castling allowed",
			pos: &Position{
				BlackKing:  bitboardFromStrs([]string{"e8"}),
				BlackRooks: bitboardFromStrs([]string{"h8"}),
				CastlingRights: CastlingRights{
					BlackOO: true,
				},
			},
			color: Black,
			legalMoves: []Move{
				{From: strToSquare("e8"), To: strToSquare("d8")},
				{From: strToSquare("e8"), To: strToSquare("d7")},
				{From: strToSquare("e8"), To: strToSquare("e7")},
				{From: strToSquare("e8"), To: strToSquare("f7")},
				{From: strToSquare("e8"), To: strToSquare("f8")},
				{From: strToSquare("e8"), To: strToSquare("g8")}, // castling
			},
		},
		{
			name: "black queenside castling allowed",
			pos: &Position{
				BlackKing:  bitboardFromStrs([]string{"e8"}),
				BlackRooks: bitboardFromStrs([]string{"a8"}),
				CastlingRights: CastlingRights{
					BlackOOO: true,
				},
			},
			color: Black,
			legalMoves: []Move{
				{From: strToSquare("e8"), To: strToSquare("d8")},
				{From: strToSquare("e8"), To: strToSquare("d7")},
				{From: strToSquare("e8"), To: strToSquare("e7")},
				{From: strToSquare("e8"), To: strToSquare("f7")},
				{From: strToSquare("e8"), To: strToSquare("f8")},
				{From: strToSquare("e8"), To: strToSquare("c8")}, // castling
			},
		},
		{
			name: "white kingside castling blocked by piece",
			pos: &Position{
				WhiteKing:    bitboardFromStrs([]string{"e1"}),
				WhiteRooks:   bitboardFromStrs([]string{"h1"}),
				WhiteKnights: bitboardFromStrs([]string{"g1"}),
				CastlingRights: CastlingRights{
					WhiteOO: true,
				},
			},
			color: White,
			legalMoves: []Move{
				{From: strToSquare("e1"), To: strToSquare("d1")},
				{From: strToSquare("e1"), To: strToSquare("d2")},
				{From: strToSquare("e1"), To: strToSquare("e2")},
				{From: strToSquare("e1"), To: strToSquare("f2")},
				{From: strToSquare("e1"), To: strToSquare("f1")},
			},
		},
		{
			name: "black queenside castling blocked by piece",
			pos: &Position{
				BlackKing:    bitboardFromStrs([]string{"e8"}),
				BlackRooks:   bitboardFromStrs([]string{"a8"}),
				BlackBishops: bitboardFromStrs([]string{"c8"}),
				CastlingRights: CastlingRights{
					BlackOOO: true,
				},
			},
			color: Black,
			legalMoves: []Move{
				{From: strToSquare("e8"), To: strToSquare("d8")},
				{From: strToSquare("e8"), To: strToSquare("d7")},
				{From: strToSquare("e8"), To: strToSquare("e7")},
				{From: strToSquare("e8"), To: strToSquare("f7")},
				{From: strToSquare("e8"), To: strToSquare("f8")},
			},
		},
		{
			name: "white both castling allowed",
			pos: &Position{
				WhiteKing:  bitboardFromStrs([]string{"e1"}),
				WhiteRooks: bitboardFromStrs([]string{"a1", "h1"}),
				CastlingRights: CastlingRights{
					WhiteOO:  true,
					WhiteOOO: true,
				},
			},
			color: White,
			legalMoves: []Move{
				{From: strToSquare("e1"), To: strToSquare("d1")},
				{From: strToSquare("e1"), To: strToSquare("d2")},
				{From: strToSquare("e1"), To: strToSquare("e2")},
				{From: strToSquare("e1"), To: strToSquare("f2")},
				{From: strToSquare("e1"), To: strToSquare("f1")},
				{From: strToSquare("e1"), To: strToSquare("g1")}, // castling
				{From: strToSquare("e1"), To: strToSquare("c1")}, // castling
			},
		},
	}

	g := NewGenerator()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			moves := g.generateKingMoves(tt.pos, tt.color)
			assert.ElementsMatch(t, tt.legalMoves, moves)
		})
	}
}

func TestIsSquareAttacked(t *testing.T) {
	tests := []struct {
		name       string
		pos        *Position
		sq         Square
		color      Color
		isAttacked bool
	}{
		{
			name: "white pawn attacks black square",
			pos: &Position{
				WhitePawns: bitboardFromStrs([]string{"d4"}),
			},
			sq:         strToSquare("e5"),
			color:      Black,
			isAttacked: true,
		},
		{
			name: "white pawn does not attack square forward",
			pos: &Position{
				WhitePawns: bitboardFromStrs([]string{"d4"}),
			},
			sq:         strToSquare("d5"),
			color:      Black,
			isAttacked: false,
		},
		{
			name: "black pawn attacks white square",
			pos: &Position{
				BlackPawns: bitboardFromStrs([]string{"e5"}),
			},
			sq:         strToSquare("f4"),
			color:      White,
			isAttacked: true,
		},
		{
			name: "black pawn does not attack backwards",
			pos: &Position{
				BlackPawns: bitboardFromStrs([]string{"d4"}),
			},
			sq:         strToSquare("e5"),
			color:      White,
			isAttacked: false,
		},
		{
			name: "knight attacks square",
			pos: &Position{
				WhiteKnights: bitboardFromStrs([]string{"e4"}),
			},
			sq:         strToSquare("f6"),
			color:      Black,
			isAttacked: true,
		},
		{
			name: "knight does not attack adjacent square",
			pos: &Position{
				WhiteKnights: bitboardFromStrs([]string{"e4"}),
			},
			sq:         strToSquare("e5"),
			color:      Black,
			isAttacked: false,
		},
		{
			name: "bishop attacks square on diagonal",
			pos: &Position{
				BlackBishops: bitboardFromStrs([]string{"a1"}),
			},
			sq:         strToSquare("h8"),
			color:      White,
			isAttacked: true,
		},
		{
			name: "bishop blocked by own piece",
			pos: &Position{
				BlackBishops: bitboardFromStrs([]string{"a1"}),
				BlackPawns:   bitboardFromStrs([]string{"d4"}),
			},
			sq:         strToSquare("h8"),
			color:      White,
			isAttacked: false,
		},
		{
			name: "bishop blocked by opponent piece",
			pos: &Position{
				BlackBishops: bitboardFromStrs([]string{"a1"}),
				WhitePawns:   bitboardFromStrs([]string{"d4"}),
			},
			sq:         strToSquare("h8"),
			color:      White,
			isAttacked: false,
		},
		{
			name: "rook attacks square on rank",
			pos: &Position{
				WhiteRooks: bitboardFromStrs([]string{"a4"}),
			},
			sq:         strToSquare("h4"),
			color:      Black,
			isAttacked: true,
		},
		{
			name: "rook attacks square on file",
			pos: &Position{
				WhiteRooks: bitboardFromStrs([]string{"d1"}),
			},
			sq:         strToSquare("d8"),
			color:      Black,
			isAttacked: true,
		},
		{
			name: "rook blocked by piece",
			pos: &Position{
				WhiteRooks: bitboardFromStrs([]string{"a4"}),
				BlackPawns: bitboardFromStrs([]string{"e4"}),
			},
			sq:         strToSquare("h4"),
			color:      Black,
			isAttacked: false,
		},
		{
			name: "queen attacks square on diagonal",
			pos: &Position{
				WhiteQueens: bitboardFromStrs([]string{"a1"}),
			},
			sq:         strToSquare("h8"),
			color:      Black,
			isAttacked: true,
		},
		{
			name: "queen attacks square on rank",
			pos: &Position{
				WhiteQueens: bitboardFromStrs([]string{"a4"}),
			},
			sq:         strToSquare("h4"),
			color:      Black,
			isAttacked: true,
		},
		{
			name: "queen attacks square on file",
			pos: &Position{
				WhiteQueens: bitboardFromStrs([]string{"d1"}),
			},
			sq:         strToSquare("d8"),
			color:      Black,
			isAttacked: true,
		},
		{
			name: "king attacks adjacent square",
			pos: &Position{
				BlackKing: bitboardFromStrs([]string{"e4"}),
			},
			sq:         strToSquare("f5"),
			color:      White,
			isAttacked: true,
		},
		{
			name: "king does not attack square two away",
			pos: &Position{
				BlackKing: bitboardFromStrs([]string{"e4"}),
			},
			sq:         strToSquare("e6"),
			color:      White,
			isAttacked: false,
		},
		{
			name: "multiple pieces attack same square",
			pos: &Position{
				WhitePawns:   bitboardFromStrs([]string{"d4"}),
				WhiteKnights: bitboardFromStrs([]string{"f3"}),
			},
			sq:         strToSquare("e5"),
			color:      Black,
			isAttacked: true,
		},
		{
			name: "no pieces attack square",
			pos: &Position{
				WhiteKnights: bitboardFromStrs([]string{"a1"}),
			},
			sq:         strToSquare("h8"),
			color:      Black,
			isAttacked: false,
		},
		{
			name: "white not attacked by own piece",
			pos: &Position{
				WhiteKnights: bitboardFromStrs([]string{"e4"}),
			},
			sq:         strToSquare("f6"),
			color:      White,
			isAttacked: false,
		},
		{
			name: "black not attacked by own piece",
			pos: &Position{
				BlackQueens: bitboardFromStrs([]string{"a4"}),
			},
			sq:         strToSquare("h4"),
			color:      Black,
			isAttacked: false,
		},
	}

	g := NewGenerator()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := g.IsSquareAttacked(tt.sq, tt.pos, tt.color)
			assert.Equal(t, tt.isAttacked, result)
		})
	}
}

func TestGenerateMoves(t *testing.T) {
	tests := []struct {
		name       string
		pos        *Position
		color      Color
		legalMoves []Move
	}{
		{
			name: "white king cant move into check",
			pos: &Position{
				WhiteKing:   bitboardFromStrs([]string{"e4"}),
				BlackQueens: bitboardFromStrs([]string{"d6"}),
				BlackPawns:  bitboardFromStrs([]string{"g6"}),
			},
			color: White,
			legalMoves: []Move{
				{From: strToSquare("e4"), To: strToSquare("f3")},
				{From: strToSquare("e4"), To: strToSquare("e3")},
			},
		},
		{
			name: "black king cant move into check",
			pos: &Position{
				BlackKing:    bitboardFromStrs([]string{"d4"}),
				WhiteRooks:   bitboardFromStrs([]string{"c5"}),
				WhiteBishops: bitboardFromStrs([]string{"g5"}),
			},
			color: Black,
			legalMoves: []Move{
				{From: strToSquare("d4"), To: strToSquare("c5")},
				{From: strToSquare("d4"), To: strToSquare("e4")},
				{From: strToSquare("d4"), To: strToSquare("d3")},
			},
		},
		{
			name: "white cant move other piece when in check",
			pos: &Position{
				WhiteKing:   bitboardFromStrs([]string{"e4"}),
				WhitePawns:  bitboardFromStrs([]string{"c2"}),
				BlackQueens: bitboardFromStrs([]string{"e6"}),
			},
			color: White,
			legalMoves: []Move{
				{From: strToSquare("e4"), To: strToSquare("d3")},
				{From: strToSquare("e4"), To: strToSquare("d4")},
				{From: strToSquare("e4"), To: strToSquare("f3")},
				{From: strToSquare("e4"), To: strToSquare("f4")},
			},
		},
		{
			name: "black cant move other piece when in check",
			pos: &Position{
				BlackKing:    bitboardFromStrs([]string{"d4"}),
				BlackRooks:   bitboardFromStrs([]string{"a1"}),
				WhiteBishops: bitboardFromStrs([]string{"f2"}),
			},
			color: Black,
			legalMoves: []Move{
				{From: strToSquare("d4"), To: strToSquare("d3")},
				{From: strToSquare("d4"), To: strToSquare("c3")},
				{From: strToSquare("d4"), To: strToSquare("c4")},
				{From: strToSquare("d4"), To: strToSquare("d5")},
				{From: strToSquare("d4"), To: strToSquare("e5")},
				{From: strToSquare("d4"), To: strToSquare("e4")},
			},
		},
		{
			name: "white can move other piece when in check to capture checking piece",
			pos: &Position{
				WhiteKing:    bitboardFromStrs([]string{"d4"}),
				WhiteBishops: bitboardFromStrs([]string{"c3"}),
				BlackQueens:  bitboardFromStrs([]string{"b4"}),
			},
			color: White,
			legalMoves: []Move{
				{From: strToSquare("c3"), To: strToSquare("b4")}, // capture
				{From: strToSquare("d4"), To: strToSquare("d3")},
				{From: strToSquare("d4"), To: strToSquare("e3")},
				{From: strToSquare("d4"), To: strToSquare("d5")},
				{From: strToSquare("d4"), To: strToSquare("e5")},
			},
		},
		{
			name: "black can move other piece when in check to capture checking piece",
			pos: &Position{
				BlackKing:    bitboardFromStrs([]string{"d4"}),
				BlackRooks:   bitboardFromStrs([]string{"a1"}),
				WhiteBishops: bitboardFromStrs([]string{"g1"}),
			},
			color: Black,
			legalMoves: []Move{
				{From: strToSquare("a1"), To: strToSquare("g1")}, // capture
				{From: strToSquare("d4"), To: strToSquare("d3")},
				{From: strToSquare("d4"), To: strToSquare("c3")},
				{From: strToSquare("d4"), To: strToSquare("c4")},
				{From: strToSquare("d4"), To: strToSquare("d5")},
				{From: strToSquare("d4"), To: strToSquare("e5")},
				{From: strToSquare("d4"), To: strToSquare("e4")},
			},
		},
		{
			name: "white can move other piece when in check to block checking piece",
			pos: &Position{
				WhiteKing:   bitboardFromStrs([]string{"e4"}),
				WhitePawns:  bitboardFromStrs([]string{"c2"}),
				BlackQueens: bitboardFromStrs([]string{"a4"}),
			},
			color: White,
			legalMoves: []Move{
				{From: strToSquare("c2"), To: strToSquare("c4")}, // block
				{From: strToSquare("e4"), To: strToSquare("d3")},
				{From: strToSquare("e4"), To: strToSquare("e3")},
				{From: strToSquare("e4"), To: strToSquare("f3")},
				{From: strToSquare("e4"), To: strToSquare("d5")},
				{From: strToSquare("e4"), To: strToSquare("e5")},
				{From: strToSquare("e4"), To: strToSquare("f5")},
			},
		},
		{
			name: "black can move other piece when in check to block checking piece",
			pos: &Position{
				BlackKing:   bitboardFromStrs([]string{"e7"}),
				BlackQueens: bitboardFromStrs([]string{"c4"}),
				WhiteRooks:  bitboardFromStrs([]string{"e1"}),
			},
			color: Black,
			legalMoves: []Move{
				{From: strToSquare("c4"), To: strToSquare("e2")}, // block
				{From: strToSquare("c4"), To: strToSquare("e4")}, // block
				{From: strToSquare("c4"), To: strToSquare("e6")}, // block
				{From: strToSquare("e7"), To: strToSquare("d6")},
				{From: strToSquare("e7"), To: strToSquare("d7")},
				{From: strToSquare("e7"), To: strToSquare("d8")},
				{From: strToSquare("e7"), To: strToSquare("f6")},
				{From: strToSquare("e7"), To: strToSquare("f7")},
				{From: strToSquare("e7"), To: strToSquare("f8")},
			},
		},
		{
			name: "white can only move king when double checked",
			pos: &Position{
				WhiteKing:   bitboardFromStrs([]string{"d5"}),
				WhiteQueens: bitboardFromStrs([]string{"b1"}),
				BlackQueens: bitboardFromStrs([]string{"b7"}),
				BlackRooks:  bitboardFromStrs([]string{"f5"}),
			},
			color: White,
			legalMoves: []Move{
				{From: strToSquare("d5"), To: strToSquare("d6")},
				{From: strToSquare("d5"), To: strToSquare("e6")},
				{From: strToSquare("d5"), To: strToSquare("d4")},
				{From: strToSquare("d5"), To: strToSquare("c4")},
			},
		},
		{
			name: "black can only move king when double checked",
			pos: &Position{
				BlackKing:    bitboardFromStrs([]string{"e8"}),
				BlackBishops: bitboardFromStrs([]string{"f6"}),
				WhiteRooks:   bitboardFromStrs([]string{"e5"}),
				WhiteKnights: bitboardFromStrs([]string{"g7"}),
			},
			color: Black,
			legalMoves: []Move{
				{From: strToSquare("e8"), To: strToSquare("d7")},
				{From: strToSquare("e8"), To: strToSquare("d8")},
				{From: strToSquare("e8"), To: strToSquare("f7")},
				{From: strToSquare("e8"), To: strToSquare("f8")},
			},
		},
		{
			name: "white evades check with en passant",
			pos: &Position{
				WhiteKing:  bitboardFromStrs([]string{"f4"}),
				WhitePawns: bitboardFromStrs([]string{"d5"}),
				BlackPawns: bitboardFromStrs([]string{"e5"}),
				EnPassant:  new(strToSquare("e6")),
			},
			color: White,
			legalMoves: []Move{
				{From: strToSquare("d5"), To: strToSquare("e6")},
				{From: strToSquare("f4"), To: strToSquare("e5")},
				{From: strToSquare("f4"), To: strToSquare("f5")},
				{From: strToSquare("f4"), To: strToSquare("g5")},
				{From: strToSquare("f4"), To: strToSquare("g4")},
				{From: strToSquare("f4"), To: strToSquare("g3")},
				{From: strToSquare("f4"), To: strToSquare("f3")},
				{From: strToSquare("f4"), To: strToSquare("e3")},
				{From: strToSquare("f4"), To: strToSquare("e4")},
			},
		},
		{
			name: "black evades check with en passant",
			pos: &Position{
				BlackKing:   bitboardFromStrs([]string{"b5"}),
				BlackPawns:  bitboardFromStrs([]string{"e4"}),
				WhiteQueens: bitboardFromStrs([]string{"f1"}),
				WhitePawns:  bitboardFromStrs([]string{"d4"}),
				EnPassant:   new(strToSquare("d3")),
			},
			color: Black,
			legalMoves: []Move{
				{From: strToSquare("e4"), To: strToSquare("d3")},
				{From: strToSquare("b5"), To: strToSquare("b6")},
				{From: strToSquare("b5"), To: strToSquare("c6")},
				{From: strToSquare("b5"), To: strToSquare("b4")},
				{From: strToSquare("b5"), To: strToSquare("a4")},
				{From: strToSquare("b5"), To: strToSquare("a5")},
			},
		},
		{
			name: "white cant move pinned piece",
			pos: &Position{
				WhiteKing:    bitboardFromStrs([]string{"e3"}),
				WhiteBishops: bitboardFromStrs([]string{"d3"}),
				BlackRooks:   bitboardFromStrs([]string{"f8"}),
				BlackQueens:  bitboardFromStrs([]string{"b3"}),
			},
			color: White,
			legalMoves: []Move{
				{From: strToSquare("e3"), To: strToSquare("d4")},
				{From: strToSquare("e3"), To: strToSquare("e4")},
				{From: strToSquare("e3"), To: strToSquare("d2")},
				{From: strToSquare("e3"), To: strToSquare("e2")},
			},
		},
		{
			name: "black cant move pinned piece",
			pos: &Position{
				BlackKing:    bitboardFromStrs([]string{"d5"}),
				BlackKnights: bitboardFromStrs([]string{"d4"}),
				BlackRooks:   bitboardFromStrs([]string{"e4"}),
				WhiteRooks:   bitboardFromStrs([]string{"d1"}),
				WhiteBishops: bitboardFromStrs([]string{"g2"}),
			},
			color: Black,
			legalMoves: []Move{
				{From: strToSquare("d5"), To: strToSquare("c4")},
				{From: strToSquare("d5"), To: strToSquare("c5")},
				{From: strToSquare("d5"), To: strToSquare("c6")},
				{From: strToSquare("d5"), To: strToSquare("d6")},
				{From: strToSquare("d5"), To: strToSquare("e6")},
				{From: strToSquare("d5"), To: strToSquare("e5")},
			},
		},
		{
			name: "white cant take en passant to leave king in check",
			pos: &Position{
				WhiteKing:   bitboardFromStrs([]string{"a5"}),
				WhitePawns:  bitboardFromStrs([]string{"d5"}),
				BlackQueens: bitboardFromStrs([]string{"h5"}),
				BlackPawns:  bitboardFromStrs([]string{"e5"}),
				EnPassant:   new(strToSquare("e6")),
			},
			color: White,
			legalMoves: []Move{
				{From: strToSquare("d5"), To: strToSquare("d6")},
				{From: strToSquare("a5"), To: strToSquare("a6")},
				{From: strToSquare("a5"), To: strToSquare("b6")},
				{From: strToSquare("a5"), To: strToSquare("b5")},
				{From: strToSquare("a5"), To: strToSquare("b4")},
				{From: strToSquare("a5"), To: strToSquare("a4")},
			},
		},
		{
			name: "white kingside castling blocked by king in check",
			pos: &Position{
				WhiteKing:  bitboardFromStrs([]string{"e1"}),
				WhiteRooks: bitboardFromStrs([]string{"h1"}),
				BlackRooks: bitboardFromStrs([]string{"e8"}),
				CastlingRights: CastlingRights{
					WhiteOO: true,
				},
			},
			color: White,
			legalMoves: []Move{
				{From: strToSquare("e1"), To: strToSquare("d1")},
				{From: strToSquare("e1"), To: strToSquare("d2")},
				{From: strToSquare("e1"), To: strToSquare("f2")},
				{From: strToSquare("e1"), To: strToSquare("f1")},
			},
		},
		{
			name: "white kingside castling blocked by attacked square",
			pos: &Position{
				WhiteKing:  bitboardFromStrs([]string{"e1"}),
				WhiteRooks: bitboardFromStrs([]string{"h1"}),
				BlackRooks: bitboardFromStrs([]string{"f8"}),
				CastlingRights: CastlingRights{
					WhiteOO: true,
				},
			},
			color: White,
			legalMoves: []Move{
				{From: strToSquare("e1"), To: strToSquare("d1")},
				{From: strToSquare("e1"), To: strToSquare("d2")},
				{From: strToSquare("e1"), To: strToSquare("e2")},
				{From: strToSquare("h1"), To: strToSquare("g1")},
				{From: strToSquare("h1"), To: strToSquare("f1")},
				{From: strToSquare("h1"), To: strToSquare("h2")},
				{From: strToSquare("h1"), To: strToSquare("h3")},
				{From: strToSquare("h1"), To: strToSquare("h4")},
				{From: strToSquare("h1"), To: strToSquare("h5")},
				{From: strToSquare("h1"), To: strToSquare("h6")},
				{From: strToSquare("h1"), To: strToSquare("h7")},
				{From: strToSquare("h1"), To: strToSquare("h8")},
			},
		},
		{
			name: "black queenside castling puts king in check",
			pos: &Position{
				BlackKing:  bitboardFromStrs([]string{"e8"}),
				BlackRooks: bitboardFromStrs([]string{"a8"}),
				WhiteRooks: bitboardFromStrs([]string{"c1"}),
				CastlingRights: CastlingRights{
					BlackOOO: true,
				},
			},
			color: Black,
			legalMoves: []Move{
				{From: strToSquare("e8"), To: strToSquare("d8")},
				{From: strToSquare("e8"), To: strToSquare("d7")},
				{From: strToSquare("e8"), To: strToSquare("e7")},
				{From: strToSquare("e8"), To: strToSquare("f7")},
				{From: strToSquare("e8"), To: strToSquare("f8")},
				{From: strToSquare("a8"), To: strToSquare("b8")},
				{From: strToSquare("a8"), To: strToSquare("c8")},
				{From: strToSquare("a8"), To: strToSquare("d8")},
				{From: strToSquare("a8"), To: strToSquare("a7")},
				{From: strToSquare("a8"), To: strToSquare("a6")},
				{From: strToSquare("a8"), To: strToSquare("a5")},
				{From: strToSquare("a8"), To: strToSquare("a4")},
				{From: strToSquare("a8"), To: strToSquare("a3")},
				{From: strToSquare("a8"), To: strToSquare("a2")},
				{From: strToSquare("a8"), To: strToSquare("a1")},
			},
		},
	}

	g := NewGenerator()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test full move generation
			moves := g.GenerateMoves(tt.pos, tt.color)
			assert.ElementsMatch(t, tt.legalMoves, moves, "full move generation")

			// Test GeneratePieceMoves for each piece type
			for _, pieceType := range []PieceType{Pawn, Knight, Bishop, Rook, Queen, King} {
				var expectedPieceMoves []Move
				for _, move := range tt.legalMoves {
					piece := tt.pos.GetPiece(move.From)
					if piece.Type == pieceType {
						expectedPieceMoves = append(expectedPieceMoves, move)
					}
				}

				pieceMoves := g.GeneratePieceMoves(tt.pos, tt.color, pieceType)
				assert.ElementsMatch(t, expectedPieceMoves, pieceMoves, "piece type: %v", pieceType)
			}
		})
	}
}

func TestGenerateMovesFullPosition(t *testing.T) {
	position := NewInitialPosition()

	g := NewGenerator()
	moves := g.GenerateMoves(position, White)
	assert.Len(t, moves, 20)
	moves = g.GenerateMoves(position, Black)
	assert.Len(t, moves, 20)

	position.MakeMove(Move{From: strToSquare("d2"), To: strToSquare("d4")})

	moves = g.GenerateMoves(position, White)
	assert.Len(t, moves, 28)
	moves = g.GenerateMoves(position, Black)
	assert.Len(t, moves, 20)
}
