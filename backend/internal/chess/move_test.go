package chess

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsEnPassant(t *testing.T) {
	tests := []struct {
		name     string
		move     Move
		pos      *Position
		expected bool
	}{
		{
			name: "white en passant",
			move: Move{
				From: strToSquare("d5"),
				To:   strToSquare("e6"),
			},
			pos: &Position{
				WhitePawns: bitboardFromStrs([]string{"d5"}),
				BlackPawns: bitboardFromStrs([]string{"e5"}),
				EnPassant:  new(strToSquare("e6")),
			},
			expected: true,
		},
		{
			name: "black en passant",
			move: Move{
				From: strToSquare("e4"),
				To:   strToSquare("d3"),
			},
			pos: &Position{
				WhitePawns: bitboardFromStrs([]string{"d4"}),
				BlackPawns: bitboardFromStrs([]string{"e4"}),
				EnPassant:  new(strToSquare("d3")),
			},
			expected: true,
		},
		{
			name: "not en passant - regular pawn move",
			move: Move{
				From: strToSquare("d4"),
				To:   strToSquare("d5"),
			},
			pos: &Position{
				WhitePawns: bitboardFromStrs([]string{"d4"}),
			},
			expected: false,
		},
		{
			name: "not en passant - no en passant square",
			move: Move{
				From: strToSquare("d5"),
				To:   strToSquare("e6"),
			},
			pos: &Position{
				WhitePawns: bitboardFromStrs([]string{"d5"}),
				BlackPawns: bitboardFromStrs([]string{"e5"}),
				EnPassant:  nil,
			},
			expected: false,
		},
		{
			name: "not en passant - wrong en passant square",
			move: Move{
				From: strToSquare("d5"),
				To:   strToSquare("e6"),
			},
			pos: &Position{
				WhitePawns: bitboardFromStrs([]string{"d5"}),
				BlackPawns: bitboardFromStrs([]string{"e5"}),
				EnPassant:  new(strToSquare("d6")),
			},
			expected: false,
		},
		{
			name: "not en passant - knight move",
			move: Move{
				From: strToSquare("b1"),
				To:   strToSquare("c3"),
			},
			pos: &Position{
				WhiteKnights: bitboardFromStrs([]string{"b1"}),
				EnPassant:    new(strToSquare("c3")),
			},
			expected: false,
		},
		{
			name: "not en passant - regular capture",
			move: Move{
				From: strToSquare("d4"),
				To:   strToSquare("e5"),
			},
			pos: &Position{
				WhitePawns: bitboardFromStrs([]string{"d4"}),
				BlackPawns: bitboardFromStrs([]string{"e5"}),
				EnPassant:  nil,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.move.IsEnPassant(tt.pos)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsCastling(t *testing.T) {
	tests := []struct {
		name     string
		move     Move
		pos      *Position
		expected bool
	}{
		{
			name: "white kingside castling",
			move: Move{
				From: strToSquare("e1"),
				To:   strToSquare("g1"),
			},
			pos: &Position{
				WhiteKing: bitboardFromStrs([]string{"e1"}),
			},
			expected: true,
		},
		{
			name: "white queenside castling",
			move: Move{
				From: strToSquare("e1"),
				To:   strToSquare("c1"),
			},
			pos: &Position{
				WhiteKing: bitboardFromStrs([]string{"e1"}),
			},
			expected: true,
		},
		{
			name: "black kingside castling",
			move: Move{
				From: strToSquare("e8"),
				To:   strToSquare("g8"),
			},
			pos: &Position{
				BlackKing: bitboardFromStrs([]string{"e8"}),
			},
			expected: true,
		},
		{
			name: "black queenside castling",
			move: Move{
				From: strToSquare("e8"),
				To:   strToSquare("c8"),
			},
			pos: &Position{
				BlackKing: bitboardFromStrs([]string{"e8"}),
			},
			expected: true,
		},
		{
			name: "not castling - king short move",
			move: Move{
				From: strToSquare("e1"),
				To:   strToSquare("e2"),
			},
			pos: &Position{
				WhiteKing: bitboardFromStrs([]string{"e1"}),
			},
			expected: false,
		},
		{
			name: "not castling - pawn move",
			move: Move{
				From: strToSquare("e2"),
				To:   strToSquare("e4"),
			},
			pos: &Position{
				WhitePawns: bitboardFromStrs([]string{"e2"}),
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.move.IsCastling(tt.pos)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestContainsMove(t *testing.T) {
	tests := []struct {
		name     string
		moves    []Move
		m        Move
		expected bool
	}{
		{
			name: "move found in list",
			moves: []Move{
				{From: strToSquare("e2"), To: strToSquare("e4")},
				{From: strToSquare("d2"), To: strToSquare("d4")},
			},
			m:        Move{From: strToSquare("e2"), To: strToSquare("e4")},
			expected: true,
		},
		{
			name: "move not found in list",
			moves: []Move{
				{From: strToSquare("e2"), To: strToSquare("e4")},
				{From: strToSquare("d2"), To: strToSquare("d4")},
			},
			m:        Move{From: strToSquare("c2"), To: strToSquare("c4")},
			expected: false,
		},
		{
			name: "move with promotion found",
			moves: []Move{
				{From: strToSquare("a7"), To: strToSquare("a8"), Promotion: new(Queen)},
				{From: strToSquare("e2"), To: strToSquare("e4")},
			},
			m:        Move{From: strToSquare("a7"), To: strToSquare("a8"), Promotion: new(Queen)},
			expected: true,
		},
		{
			name: "move with different promotion not found",
			moves: []Move{
				{From: strToSquare("a7"), To: strToSquare("a8"), Promotion: new(Queen)},
			},
			m:        Move{From: strToSquare("a7"), To: strToSquare("a8"), Promotion: new(Rook)},
			expected: false,
		},
		{
			name: "move without promotion not found when list has promotion",
			moves: []Move{
				{From: strToSquare("a7"), To: strToSquare("a8"), Promotion: new(Queen)},
			},
			m:        Move{From: strToSquare("a7"), To: strToSquare("a8")},
			expected: false,
		},
		{
			name:     "empty moves list",
			moves:    []Move{},
			m:        Move{From: strToSquare("e2"), To: strToSquare("e4")},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ContainsMove(tt.moves, tt.m)
			assert.Equal(t, tt.expected, result)
		})
	}
}
