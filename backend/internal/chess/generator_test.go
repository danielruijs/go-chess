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
				{Piece: Pawn, From: strToSquare("a3"), To: strToSquare("a4")}, // a3 -> a4
			},
		},
		{
			name: "black single push",
			pos: &Position{
				BlackPawns: bitboardFromStrs([]string{"a6"}),
			},
			color: Black,
			legalMoves: []Move{
				{Piece: Pawn, From: strToSquare("a6"), To: strToSquare("a5")}, // a6 -> a5
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
				{Piece: Pawn, From: strToSquare("a2"), To: strToSquare("a3")},
				{Piece: Pawn, From: strToSquare("a2"), To: strToSquare("a4")},
			},
		},
		{
			name: "black double push",
			pos: &Position{
				BlackPawns: bitboardFromStrs([]string{"a7"}),
			},
			color: Black,
			legalMoves: []Move{
				{Piece: Pawn, From: strToSquare("a7"), To: strToSquare("a6")},
				{Piece: Pawn, From: strToSquare("a7"), To: strToSquare("a5")},
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
				{Piece: Pawn, From: strToSquare("a2"), To: strToSquare("a3")},
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
				{Piece: Pawn, From: strToSquare("a7"), To: strToSquare("a6")},
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
				{Piece: Pawn, From: strToSquare("d4"), To: strToSquare("d5")},
				{Piece: Pawn, From: strToSquare("d4"), To: strToSquare("c5")},
				{Piece: Pawn, From: strToSquare("d4"), To: strToSquare("e5")},
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
				{Piece: Pawn, From: strToSquare("c5"), To: strToSquare("c4")},
				{Piece: Pawn, From: strToSquare("c5"), To: strToSquare("b4")},
				{Piece: Pawn, From: strToSquare("c5"), To: strToSquare("d4")},
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
				{Piece: Pawn, From: strToSquare("d5"), To: strToSquare("d6")},
				{Piece: Pawn, From: strToSquare("d5"), To: strToSquare("e6")},
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
				{Piece: Pawn, From: strToSquare("e4"), To: strToSquare("e3")},
				{Piece: Pawn, From: strToSquare("e4"), To: strToSquare("d3")},
			},
		},
		{
			name: "white promotion",
			pos: &Position{
				WhitePawns: bitboardFromStrs([]string{"a7"}),
			},
			color: White,
			legalMoves: []Move{
				{Piece: Pawn, From: strToSquare("a7"), To: strToSquare("a8"), Promotion: toPtr(Queen)},
				{Piece: Pawn, From: strToSquare("a7"), To: strToSquare("a8"), Promotion: toPtr(Rook)},
				{Piece: Pawn, From: strToSquare("a7"), To: strToSquare("a8"), Promotion: toPtr(Bishop)},
				{Piece: Pawn, From: strToSquare("a7"), To: strToSquare("a8"), Promotion: toPtr(Knight)},
			},
		},
		{
			name: "black promotion",
			pos: &Position{
				BlackPawns: bitboardFromStrs([]string{"a2"}),
			},
			color: Black,
			legalMoves: []Move{
				{Piece: Pawn, From: strToSquare("a2"), To: strToSquare("a1"), Promotion: toPtr(Queen)},
				{Piece: Pawn, From: strToSquare("a2"), To: strToSquare("a1"), Promotion: toPtr(Rook)},
				{Piece: Pawn, From: strToSquare("a2"), To: strToSquare("a1"), Promotion: toPtr(Bishop)},
				{Piece: Pawn, From: strToSquare("a2"), To: strToSquare("a1"), Promotion: toPtr(Knight)},
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

func TestKnightMoveGenerator_generateMoves(t *testing.T) {
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
				{Piece: Knight, From: strToSquare("e4"), To: strToSquare("f6")},
				{Piece: Knight, From: strToSquare("e4"), To: strToSquare("g5")},
				{Piece: Knight, From: strToSquare("e4"), To: strToSquare("g3")},
				{Piece: Knight, From: strToSquare("e4"), To: strToSquare("f2")},
				{Piece: Knight, From: strToSquare("e4"), To: strToSquare("d2")},
				{Piece: Knight, From: strToSquare("e4"), To: strToSquare("c3")},
				{Piece: Knight, From: strToSquare("e4"), To: strToSquare("c5")},
				{Piece: Knight, From: strToSquare("e4"), To: strToSquare("d6")},
			},
		},
		{
			name: "black clear board",
			pos: &Position{
				BlackKnights: bitboardFromStrs([]string{"d4"}),
			},
			color: Black,
			legalMoves: []Move{
				{Piece: Knight, From: strToSquare("d4"), To: strToSquare("e6")},
				{Piece: Knight, From: strToSquare("d4"), To: strToSquare("f5")},
				{Piece: Knight, From: strToSquare("d4"), To: strToSquare("f3")},
				{Piece: Knight, From: strToSquare("d4"), To: strToSquare("e2")},
				{Piece: Knight, From: strToSquare("d4"), To: strToSquare("c2")},
				{Piece: Knight, From: strToSquare("d4"), To: strToSquare("b3")},
				{Piece: Knight, From: strToSquare("d4"), To: strToSquare("b5")},
				{Piece: Knight, From: strToSquare("d4"), To: strToSquare("c6")},
			},
		},
		{
			name: "white corner",
			pos: &Position{
				WhiteKnights: bitboardFromStrs([]string{"a1"}),
			},
			color: White,
			legalMoves: []Move{
				{Piece: Knight, From: strToSquare("a1"), To: strToSquare("b3")},
				{Piece: Knight, From: strToSquare("a1"), To: strToSquare("c2")},
			},
		},
		{
			name: "black corner",
			pos: &Position{
				BlackKnights: bitboardFromStrs([]string{"a8"}),
			},
			color: Black,
			legalMoves: []Move{
				{Piece: Knight, From: strToSquare("a8"), To: strToSquare("b6")},
				{Piece: Knight, From: strToSquare("a8"), To: strToSquare("c7")},
			},
		},
		{
			name: "white edge",
			pos: &Position{
				WhiteKnights: bitboardFromStrs([]string{"a4"}),
			},
			color: White,
			legalMoves: []Move{
				{Piece: Knight, From: strToSquare("a4"), To: strToSquare("b6")},
				{Piece: Knight, From: strToSquare("a4"), To: strToSquare("c5")},
				{Piece: Knight, From: strToSquare("a4"), To: strToSquare("c3")},
				{Piece: Knight, From: strToSquare("a4"), To: strToSquare("b2")},
			},
		},
		{
			name: "black edge",
			pos: &Position{
				BlackKnights: bitboardFromStrs([]string{"e8"}),
			},
			color: Black,
			legalMoves: []Move{
				{Piece: Knight, From: strToSquare("e8"), To: strToSquare("g7")},
				{Piece: Knight, From: strToSquare("e8"), To: strToSquare("f6")},
				{Piece: Knight, From: strToSquare("e8"), To: strToSquare("d6")},
				{Piece: Knight, From: strToSquare("e8"), To: strToSquare("c7")},
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
				{Piece: Knight, From: strToSquare("e4"), To: strToSquare("g3")},
				{Piece: Knight, From: strToSquare("e4"), To: strToSquare("f2")},
				{Piece: Knight, From: strToSquare("e4"), To: strToSquare("d2")},
				{Piece: Knight, From: strToSquare("e4"), To: strToSquare("c5")},
				{Piece: Knight, From: strToSquare("e4"), To: strToSquare("d6")},
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
				{Piece: Knight, From: strToSquare("d4"), To: strToSquare("f3")},
				{Piece: Knight, From: strToSquare("d4"), To: strToSquare("c2")},
				{Piece: Knight, From: strToSquare("d4"), To: strToSquare("b3")},
				{Piece: Knight, From: strToSquare("d4"), To: strToSquare("b5")},
				{Piece: Knight, From: strToSquare("d4"), To: strToSquare("c6")},
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
				{Piece: Knight, From: strToSquare("e4"), To: strToSquare("f6")},
				{Piece: Knight, From: strToSquare("e4"), To: strToSquare("g5")},
				{Piece: Knight, From: strToSquare("e4"), To: strToSquare("g3")},
				{Piece: Knight, From: strToSquare("e4"), To: strToSquare("f2")},
				{Piece: Knight, From: strToSquare("e4"), To: strToSquare("d2")},
				{Piece: Knight, From: strToSquare("e4"), To: strToSquare("c3")},
				{Piece: Knight, From: strToSquare("e4"), To: strToSquare("c5")},
				{Piece: Knight, From: strToSquare("e4"), To: strToSquare("d6")},
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
				{Piece: Knight, From: strToSquare("d4"), To: strToSquare("e6")},
				{Piece: Knight, From: strToSquare("d4"), To: strToSquare("f5")},
				{Piece: Knight, From: strToSquare("d4"), To: strToSquare("f3")},
				{Piece: Knight, From: strToSquare("d4"), To: strToSquare("e2")},
				{Piece: Knight, From: strToSquare("d4"), To: strToSquare("c2")},
				{Piece: Knight, From: strToSquare("d4"), To: strToSquare("b3")},
				{Piece: Knight, From: strToSquare("d4"), To: strToSquare("b5")},
				{Piece: Knight, From: strToSquare("d4"), To: strToSquare("c6")},
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

func TestBishopMoveGenerator_generateMoves(t *testing.T) {
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
				{Piece: Bishop, From: strToSquare("e4"), To: strToSquare("d5")},
				{Piece: Bishop, From: strToSquare("e4"), To: strToSquare("c6")},
				{Piece: Bishop, From: strToSquare("e4"), To: strToSquare("b7")},
				{Piece: Bishop, From: strToSquare("e4"), To: strToSquare("a8")},
				{Piece: Bishop, From: strToSquare("e4"), To: strToSquare("f5")},
				{Piece: Bishop, From: strToSquare("e4"), To: strToSquare("g6")},
				{Piece: Bishop, From: strToSquare("e4"), To: strToSquare("h7")},
				{Piece: Bishop, From: strToSquare("e4"), To: strToSquare("d3")},
				{Piece: Bishop, From: strToSquare("e4"), To: strToSquare("c2")},
				{Piece: Bishop, From: strToSquare("e4"), To: strToSquare("b1")},
				{Piece: Bishop, From: strToSquare("e4"), To: strToSquare("f3")},
				{Piece: Bishop, From: strToSquare("e4"), To: strToSquare("g2")},
				{Piece: Bishop, From: strToSquare("e4"), To: strToSquare("h1")},
			},
		},
		{
			name: "black clear board",
			pos: &Position{
				BlackBishops: bitboardFromStrs([]string{"d4"}),
			},
			color: Black,
			legalMoves: []Move{
				{Piece: Bishop, From: strToSquare("d4"), To: strToSquare("c5")},
				{Piece: Bishop, From: strToSquare("d4"), To: strToSquare("b6")},
				{Piece: Bishop, From: strToSquare("d4"), To: strToSquare("a7")},
				{Piece: Bishop, From: strToSquare("d4"), To: strToSquare("e5")},
				{Piece: Bishop, From: strToSquare("d4"), To: strToSquare("f6")},
				{Piece: Bishop, From: strToSquare("d4"), To: strToSquare("g7")},
				{Piece: Bishop, From: strToSquare("d4"), To: strToSquare("h8")},
				{Piece: Bishop, From: strToSquare("d4"), To: strToSquare("c3")},
				{Piece: Bishop, From: strToSquare("d4"), To: strToSquare("b2")},
				{Piece: Bishop, From: strToSquare("d4"), To: strToSquare("a1")},
				{Piece: Bishop, From: strToSquare("d4"), To: strToSquare("e3")},
				{Piece: Bishop, From: strToSquare("d4"), To: strToSquare("f2")},
				{Piece: Bishop, From: strToSquare("d4"), To: strToSquare("g1")},
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
				{Piece: Bishop, From: strToSquare("e4"), To: strToSquare("d5")},
				{Piece: Bishop, From: strToSquare("e4"), To: strToSquare("c6")},
				{Piece: Bishop, From: strToSquare("e4"), To: strToSquare("b7")},
				{Piece: Bishop, From: strToSquare("e4"), To: strToSquare("a8")},
				{Piece: Bishop, From: strToSquare("e4"), To: strToSquare("d3")},
				{Piece: Bishop, From: strToSquare("e4"), To: strToSquare("f3")},
				{Piece: Bishop, From: strToSquare("e4"), To: strToSquare("g2")},
				{Piece: Bishop, From: strToSquare("e4"), To: strToSquare("h1")},
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
				{Piece: Bishop, From: strToSquare("d4"), To: strToSquare("c5")},
				{Piece: Bishop, From: strToSquare("d4"), To: strToSquare("b6")},
				{Piece: Bishop, From: strToSquare("d4"), To: strToSquare("a7")},
				{Piece: Bishop, From: strToSquare("d4"), To: strToSquare("e5")},
				{Piece: Bishop, From: strToSquare("d4"), To: strToSquare("c3")},
				{Piece: Bishop, From: strToSquare("d4"), To: strToSquare("b2")},
				{Piece: Bishop, From: strToSquare("d4"), To: strToSquare("a1")},
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
				{Piece: Bishop, From: strToSquare("e4"), To: strToSquare("d5")},
				{Piece: Bishop, From: strToSquare("e4"), To: strToSquare("c6")},
				{Piece: Bishop, From: strToSquare("e4"), To: strToSquare("b7")},
				{Piece: Bishop, From: strToSquare("e4"), To: strToSquare("a8")},
				{Piece: Bishop, From: strToSquare("e4"), To: strToSquare("f5")}, // capture
				{Piece: Bishop, From: strToSquare("e4"), To: strToSquare("d3")},
				{Piece: Bishop, From: strToSquare("e4"), To: strToSquare("c2")}, // capture
				{Piece: Bishop, From: strToSquare("e4"), To: strToSquare("f3")},
				{Piece: Bishop, From: strToSquare("e4"), To: strToSquare("g2")},
				{Piece: Bishop, From: strToSquare("e4"), To: strToSquare("h1")},
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
				{Piece: Bishop, From: strToSquare("d4"), To: strToSquare("c5")},
				{Piece: Bishop, From: strToSquare("d4"), To: strToSquare("b6")},
				{Piece: Bishop, From: strToSquare("d4"), To: strToSquare("a7")},
				{Piece: Bishop, From: strToSquare("d4"), To: strToSquare("e5")},
				{Piece: Bishop, From: strToSquare("d4"), To: strToSquare("f6")}, // capture
				{Piece: Bishop, From: strToSquare("d4"), To: strToSquare("c3")},
				{Piece: Bishop, From: strToSquare("d4"), To: strToSquare("b2")},
				{Piece: Bishop, From: strToSquare("d4"), To: strToSquare("a1")},
				{Piece: Bishop, From: strToSquare("d4"), To: strToSquare("e3")}, // capture
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
				{Piece: Bishop, From: strToSquare("a1"), To: strToSquare("b2")},
				{Piece: Bishop, From: strToSquare("a1"), To: strToSquare("c3")},
				{Piece: Bishop, From: strToSquare("a1"), To: strToSquare("d4")},
				{Piece: Bishop, From: strToSquare("a1"), To: strToSquare("e5")},
				{Piece: Bishop, From: strToSquare("a1"), To: strToSquare("f6")},
				{Piece: Bishop, From: strToSquare("a1"), To: strToSquare("g7")},
				{Piece: Bishop, From: strToSquare("a1"), To: strToSquare("h8")},
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
				{Piece: Bishop, From: strToSquare("a8"), To: strToSquare("b7")},
				{Piece: Bishop, From: strToSquare("a8"), To: strToSquare("c6")},
				{Piece: Bishop, From: strToSquare("a8"), To: strToSquare("d5")},
				{Piece: Bishop, From: strToSquare("a8"), To: strToSquare("e4")},
				{Piece: Bishop, From: strToSquare("a8"), To: strToSquare("f3")},
				{Piece: Bishop, From: strToSquare("a8"), To: strToSquare("g2")},
				{Piece: Bishop, From: strToSquare("a8"), To: strToSquare("h1")},
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
