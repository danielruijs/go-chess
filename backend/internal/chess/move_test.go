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

func TestToSAN(t *testing.T) {
	tests := []struct {
		name     string
		move     Move
		pos      *Position
		expected string
	}{
		{
			name: "pawn move",
			move: Move{
				From: strToSquare("e2"),
				To:   strToSquare("e4"),
			},
			pos: &Position{
				WhitePawns: bitboardFromStrs([]string{"e2"}),
			},
			expected: "e4",
		},
		{
			name: "knight move",
			move: Move{
				From: strToSquare("b1"),
				To:   strToSquare("c3"),
			},
			pos: &Position{
				WhiteKnights: bitboardFromStrs([]string{"b1"}),
			},
			expected: "Nc3",
		},
		{
			name: "bishop move",
			move: Move{
				From: strToSquare("f1"),
				To:   strToSquare("c4"),
			},
			pos: &Position{
				WhiteBishops: bitboardFromStrs([]string{"f1"}),
			},
			expected: "Bc4",
		},
		{
			name: "rook move",
			move: Move{
				From: strToSquare("a1"),
				To:   strToSquare("a4"),
			},
			pos: &Position{
				WhiteRooks: bitboardFromStrs([]string{"a1"}),
			},
			expected: "Ra4",
		},
		{
			name: "queen move",
			move: Move{
				From: strToSquare("d1"),
				To:   strToSquare("d4"),
			},
			pos: &Position{
				WhiteQueens: bitboardFromStrs([]string{"d1"}),
			},
			expected: "Qd4",
		},
		{
			name: "king move",
			move: Move{
				From: strToSquare("e1"),
				To:   strToSquare("e2"),
			},
			pos: &Position{
				WhiteKing: bitboardFromStrs([]string{"e1"}),
			},
			expected: "Ke2",
		},
		{
			name: "pawn capture",
			move: Move{
				From: strToSquare("e4"),
				To:   strToSquare("d5"),
			},
			pos: &Position{
				WhitePawns: bitboardFromStrs([]string{"e4"}),
				BlackPawns: bitboardFromStrs([]string{"d5"}),
			},
			expected: "exd5",
		},
		{
			name: "knight capture",
			move: Move{
				From: strToSquare("c3"),
				To:   strToSquare("d5"),
			},
			pos: &Position{
				WhiteKnights: bitboardFromStrs([]string{"c3"}),
				BlackPawns:   bitboardFromStrs([]string{"d5"}),
			},
			expected: "Nxd5",
		},
		{
			name: "en passant",
			move: Move{
				From: strToSquare("e5"),
				To:   strToSquare("d6"),
			},
			pos: &Position{
				WhitePawns: bitboardFromStrs([]string{"e5"}),
				BlackPawns: bitboardFromStrs([]string{"d5"}),
				EnPassant:  new(strToSquare("d6")),
			},
			expected: "exd6",
		},
		{
			name: "pawn promotion to queen",
			move: Move{
				From:      strToSquare("a7"),
				To:        strToSquare("a8"),
				Promotion: new(Queen),
			},
			pos: &Position{
				WhitePawns: bitboardFromStrs([]string{"a7"}),
			},
			expected: "a8=Q",
		},
		{
			name: "pawn promotion to knight",
			move: Move{
				From:      strToSquare("a2"),
				To:        strToSquare("a1"),
				Promotion: new(Knight),
			},
			pos: &Position{
				BlackPawns: bitboardFromStrs([]string{"a2"}),
			},
			expected: "a1=N",
		},
		{
			name: "kingside castling",
			move: Move{
				From: strToSquare("e1"),
				To:   strToSquare("g1"),
			},
			pos: &Position{
				WhiteKing:  bitboardFromStrs([]string{"e1"}),
				WhiteRooks: bitboardFromStrs([]string{"h1"}),
				CastlingRights: CastlingRights{
					WhiteOO: true,
				},
			},
			expected: "O-O",
		},
		{
			name: "queenside castling",
			move: Move{
				From: strToSquare("e8"),
				To:   strToSquare("c8"),
			},
			pos: &Position{
				BlackKing:  bitboardFromStrs([]string{"e8"}),
				BlackRooks: bitboardFromStrs([]string{"a8"}),
				CastlingRights: CastlingRights{
					BlackOOO: true,
				},
			},
			expected: "O-O-O",
		},
		{
			name: "knight file disambiguation",
			move: Move{
				From: strToSquare("g1"),
				To:   strToSquare("f3"),
			},
			pos: &Position{
				WhiteKnights: bitboardFromStrs([]string{"g1", "d2"}),
			},
			expected: "Ngf3",
		},
		{
			name: "rook file disambiguation",
			move: Move{
				From: strToSquare("d8"),
				To:   strToSquare("f8"),
			},
			pos: &Position{
				BlackRooks: bitboardFromStrs([]string{"d8", "h8"}),
			},
			expected: "Rdf8",
		},
		{
			name: "rook rank disambiguation",
			move: Move{
				From: strToSquare("a1"),
				To:   strToSquare("a3"),
			},
			pos: &Position{
				WhiteRooks: bitboardFromStrs([]string{"a1", "a5"}),
			},
			expected: "R1a3",
		},
		{
			name: "queen double disambiguation",
			move: Move{
				From: strToSquare("h4"),
				To:   strToSquare("e1"),
			},
			pos: &Position{
				WhiteQueens: bitboardFromStrs([]string{"h1", "h4", "e4"}),
			},
			expected: "Qh4e1",
		},
		{
			name: "check",
			move: Move{
				From: strToSquare("d1"),
				To:   strToSquare("h5"),
			},
			pos: &Position{
				WhiteQueens: bitboardFromStrs([]string{"d1"}),
				BlackKing:   bitboardFromStrs([]string{"h8"}),
			},
			expected: "Qh5+",
		},
		{
			name: "checkmate",
			move: Move{
				From:      strToSquare("f7"),
				To:        strToSquare("f8"),
				Promotion: new(Queen),
			},
			pos: &Position{
				WhitePawns: bitboardFromStrs([]string{"f7"}),
				WhiteKing:  bitboardFromStrs([]string{"g6"}),
				BlackKing:  bitboardFromStrs([]string{"h8"}),
			},
			expected: "f8=Q#",
		},
		{
			name: "pawn capture promotion",
			move: Move{
				From:      strToSquare("g7"),
				To:        strToSquare("h8"),
				Promotion: new(Bishop),
			},
			pos: &Position{
				WhitePawns: bitboardFromStrs([]string{"g7"}),
				BlackRooks: bitboardFromStrs([]string{"h8"}),
			},
			expected: "gxh8=B",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.move.ToSAN(tt.pos)
			assert.Equal(t, tt.expected, result)
		})
	}
}
