package chess

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppendMoveFromInitialPositionQueensGambit(t *testing.T) {
	gen := NewGenerator()
	position := NewInitialPosition()
	pgn := PGN("")

	moves := []Move{
		{From: strToSquare("d2"), To: strToSquare("d4")},
		{From: strToSquare("d7"), To: strToSquare("d5")},
		{From: strToSquare("c2"), To: strToSquare("c4")},
		{From: strToSquare("g8"), To: strToSquare("f6")},
		{From: strToSquare("c4"), To: strToSquare("d5")},
		{From: strToSquare("f6"), To: strToSquare("d5")},
		{From: strToSquare("e2"), To: strToSquare("e4")},
		{From: strToSquare("d5"), To: strToSquare("f6")},
	}

	for _, move := range moves {
		pgn = pgn.AppendMove(move, position, gen)
		position.MakeMove(move)
	}

	expected := PGN("1.d4 d5 2.c4 Nf6 3.cxd5 Nxd5 4.e4 Nf6 ")
	assert.Equal(t, expected, pgn)
}

func TestAppendMoveFromInitialPositionCaroKann(t *testing.T) {
	gen := NewGenerator()
	position := NewInitialPosition()
	pgn := PGN("")

	moves := []Move{
		{From: strToSquare("e2"), To: strToSquare("e4")},
		{From: strToSquare("c7"), To: strToSquare("c6")},
		{From: strToSquare("d2"), To: strToSquare("d4")},
		{From: strToSquare("d7"), To: strToSquare("d5")},
		{From: strToSquare("b1"), To: strToSquare("c3")},
		{From: strToSquare("d5"), To: strToSquare("e4")},
		{From: strToSquare("c3"), To: strToSquare("e4")},
		{From: strToSquare("c8"), To: strToSquare("f5")},
		{From: strToSquare("e4"), To: strToSquare("g3")},
		{From: strToSquare("f5"), To: strToSquare("g6")},
		{From: strToSquare("g1"), To: strToSquare("f3")},
		{From: strToSquare("e7"), To: strToSquare("e6")},
		{From: strToSquare("h2"), To: strToSquare("h4")},
		{From: strToSquare("h7"), To: strToSquare("h6")},
		{From: strToSquare("f1"), To: strToSquare("d3")},
		{From: strToSquare("g6"), To: strToSquare("d3")},
	}

	for _, move := range moves {
		pgn = pgn.AppendMove(move, position, gen)
		position.MakeMove(move)
	}

	expected := PGN("1.e4 c6 2.d4 d5 3.Nc3 dxe4 4.Nxe4 Bf5 5.Ng3 Bg6 6.Nf3 e6 7.h4 h6 8.Bd3 Bxd3 ")
	assert.Equal(t, expected, pgn)
}

func TestAppendMoveFromInitialPositionBlackFirstMove(t *testing.T) {
	gen := NewGenerator()
	position := NewInitialPosition()
	position.ActiveColor = Black
	pgn := PGN("")

	moves := []Move{
		{From: strToSquare("d7"), To: strToSquare("d5")},
		{From: strToSquare("c2"), To: strToSquare("c4")},
		{From: strToSquare("d5"), To: strToSquare("c4")},
		{From: strToSquare("b1"), To: strToSquare("c3")},
	}

	for _, move := range moves {
		pgn = pgn.AppendMove(move, position, gen)
		position.MakeMove(move)
	}

	expected := PGN("1... d5 2.c4 dxc4 3.Nc3 ")
	assert.Equal(t, expected, pgn)
}

func TestAppendResult(t *testing.T) {
	tests := []struct {
		name     string
		pgn      PGN
		result   *Result
		expected PGN
	}{
		{
			name:     "white win",
			pgn:      PGN("1.e4 e5 2.Nf3 Nc6 "),
			result:   &Result{Outcome: WhiteWin},
			expected: PGN("1.e4 e5 2.Nf3 Nc6 1-0"),
		},
		{
			name:     "black win",
			pgn:      PGN("1.e4 e5 2.Nf3 Nc6 "),
			result:   &Result{Outcome: BlackWin},
			expected: PGN("1.e4 e5 2.Nf3 Nc6 0-1"),
		},
		{
			name:     "draw",
			pgn:      PGN("1.e4 e5 2.Nf3 Nc6 "),
			result:   &Result{Outcome: Draw},
			expected: PGN("1.e4 e5 2.Nf3 Nc6 1/2-1/2"),
		},
		{
			name:     "unknown result",
			pgn:      PGN("1.e4 e5 "),
			result:   &Result{},
			expected: PGN("1.e4 e5 *"),
		},
		{
			name:     "white win no trailing space",
			pgn:      PGN("1.e4 e5 2.Nf3 Nc6"),
			result:   &Result{Outcome: WhiteWin},
			expected: PGN("1.e4 e5 2.Nf3 Nc6 1-0"),
		},
		{
			name:     "unknown result no trailing space",
			pgn:      PGN("1.e4 e5"),
			result:   &Result{},
			expected: PGN("1.e4 e5 *"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.pgn.AppendResult(tt.result))
		})
	}
}
