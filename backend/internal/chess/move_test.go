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
				EnPassant:  toPtr(strToSquare("e6")),
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
				EnPassant:  toPtr(strToSquare("d3")),
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
				EnPassant:  toPtr(strToSquare("d6")),
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
				EnPassant:    toPtr(strToSquare("c3")),
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
