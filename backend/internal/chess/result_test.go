package chess

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetResult(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() *Engine
		expected *Result
	}{
		{
			name: "no result - game started",
			setup: func() *Engine {
				e := NewEngine()
				return e
			},
			expected: nil,
		},
		{
			name: "checkmate - black wins",
			setup: func() *Engine {
				e := NewEngine()
				e.ApplyMove(Move{From: strToSquare("f2"), To: strToSquare("f3"), Promotion: nil}, White)
				e.ApplyMove(Move{From: strToSquare("e7"), To: strToSquare("e6"), Promotion: nil}, Black)
				e.ApplyMove(Move{From: strToSquare("g2"), To: strToSquare("g4"), Promotion: nil}, White)
				e.ApplyMove(Move{From: strToSquare("d8"), To: strToSquare("h4"), Promotion: nil}, Black)
				return e
			},
			expected: &Result{Outcome: BlackWin, Reason: Checkmate},
		},
		{
			name: "checkmate - white wins",
			setup: func() *Engine {
				e := NewEngine()
				e.ApplyMove(Move{From: strToSquare("e2"), To: strToSquare("e4"), Promotion: nil}, White)
				e.ApplyMove(Move{From: strToSquare("e7"), To: strToSquare("e5"), Promotion: nil}, Black)
				e.ApplyMove(Move{From: strToSquare("f1"), To: strToSquare("c4"), Promotion: nil}, White)
				e.ApplyMove(Move{From: strToSquare("b8"), To: strToSquare("c6"), Promotion: nil}, Black)
				e.ApplyMove(Move{From: strToSquare("d1"), To: strToSquare("h5"), Promotion: nil}, White)
				e.ApplyMove(Move{From: strToSquare("g8"), To: strToSquare("f6"), Promotion: nil}, Black)
				e.ApplyMove(Move{From: strToSquare("h5"), To: strToSquare("f7"), Promotion: nil}, White)
				return e
			},
			expected: &Result{Outcome: WhiteWin, Reason: Checkmate},
		},
		{
			name: "stalemate",
			setup: func() *Engine {
				e := NewEngine()
				e.position = &Position{
					WhiteKing:   bitboardFromStrs([]string{"c7"}),
					BlackKing:   bitboardFromStrs([]string{"a8"}),
					WhiteQueens: bitboardFromStrs([]string{"c7"}),
					ActiveColor: Black,
				}
				return e
			},
			expected: &Result{Outcome: Draw, Reason: Stalemate},
		},
		{
			name: "not stalemate - white to move",
			setup: func() *Engine {
				e := NewEngine()
				e.position = &Position{
					WhiteKing:   bitboardFromStrs([]string{"c7"}),
					BlackKing:   bitboardFromStrs([]string{"a8"}),
					WhiteQueens: bitboardFromStrs([]string{"c7"}),
					ActiveColor: White,
				}
				return e
			},
			expected: nil,
		},
		{
			name: "threefold repetition - simple repetition",
			setup: func() *Engine {
				e := NewEngine()
				e.ApplyMove(Move{From: strToSquare("g1"), To: strToSquare("f3")}, White)
				e.ApplyMove(Move{From: strToSquare("g8"), To: strToSquare("f6")}, Black)
				e.ApplyMove(Move{From: strToSquare("f3"), To: strToSquare("g1")}, White)
				e.ApplyMove(Move{From: strToSquare("f6"), To: strToSquare("g8")}, Black)
				e.ApplyMove(Move{From: strToSquare("g1"), To: strToSquare("f3")}, White)
				e.ApplyMove(Move{From: strToSquare("g8"), To: strToSquare("f6")}, Black)
				e.ApplyMove(Move{From: strToSquare("f3"), To: strToSquare("g1")}, White)
				e.ApplyMove(Move{From: strToSquare("f6"), To: strToSquare("g8")}, Black)
				return e
			},
			expected: &Result{Outcome: Draw, Reason: ThreefoldRepetition},
		},
		{
			name: "threefold repetition - non-consecutive repetition",
			setup: func() *Engine {
				e := NewEngine()
				e.ApplyMove(Move{From: strToSquare("e2"), To: strToSquare("e4")}, White)
				e.ApplyMove(Move{From: strToSquare("e7"), To: strToSquare("e5")}, Black)
				e.ApplyMove(Move{From: strToSquare("e1"), To: strToSquare("e2")}, White)
				e.ApplyMove(Move{From: strToSquare("e8"), To: strToSquare("e7")}, Black)
				e.ApplyMove(Move{From: strToSquare("e2"), To: strToSquare("e1")}, White)
				e.ApplyMove(Move{From: strToSquare("e7"), To: strToSquare("e8")}, Black)
				e.ApplyMove(Move{From: strToSquare("b1"), To: strToSquare("c3")}, White)
				e.ApplyMove(Move{From: strToSquare("b8"), To: strToSquare("c6")}, Black)
				e.ApplyMove(Move{From: strToSquare("e1"), To: strToSquare("e2")}, White)
				e.ApplyMove(Move{From: strToSquare("e8"), To: strToSquare("e7")}, Black)
				e.ApplyMove(Move{From: strToSquare("e2"), To: strToSquare("e1")}, White)
				e.ApplyMove(Move{From: strToSquare("e7"), To: strToSquare("e8")}, Black)
				e.ApplyMove(Move{From: strToSquare("c3"), To: strToSquare("b1")}, White)
				e.ApplyMove(Move{From: strToSquare("c6"), To: strToSquare("b8")}, Black)
				e.ApplyMove(Move{From: strToSquare("e1"), To: strToSquare("e2")}, White)
				e.ApplyMove(Move{From: strToSquare("e8"), To: strToSquare("e7")}, Black)
				e.ApplyMove(Move{From: strToSquare("e2"), To: strToSquare("e1")}, White)
				e.ApplyMove(Move{From: strToSquare("e7"), To: strToSquare("e8")}, Black)
				return e
			},
			expected: &Result{Outcome: Draw, Reason: ThreefoldRepetition},
		},
		{
			name: "not repetition - castling rights changed",
			setup: func() *Engine {
				e := NewEngine()
				e.ApplyMove(Move{From: strToSquare("e2"), To: strToSquare("e4")}, White)
				e.ApplyMove(Move{From: strToSquare("e7"), To: strToSquare("e5")}, Black)
				e.ApplyMove(Move{From: strToSquare("e1"), To: strToSquare("e2")}, White)
				e.ApplyMove(Move{From: strToSquare("e8"), To: strToSquare("e7")}, Black)
				e.ApplyMove(Move{From: strToSquare("e2"), To: strToSquare("e1")}, White)
				e.ApplyMove(Move{From: strToSquare("e7"), To: strToSquare("e8")}, Black)
				e.ApplyMove(Move{From: strToSquare("b1"), To: strToSquare("c3")}, White)
				e.ApplyMove(Move{From: strToSquare("b8"), To: strToSquare("c6")}, Black)
				e.ApplyMove(Move{From: strToSquare("e1"), To: strToSquare("e2")}, White)
				e.ApplyMove(Move{From: strToSquare("e8"), To: strToSquare("e7")}, Black)
				e.ApplyMove(Move{From: strToSquare("e2"), To: strToSquare("e1")}, White)
				e.ApplyMove(Move{From: strToSquare("e7"), To: strToSquare("e8")}, Black)
				e.ApplyMove(Move{From: strToSquare("c3"), To: strToSquare("b1")}, White)
				e.ApplyMove(Move{From: strToSquare("c6"), To: strToSquare("b8")}, Black)
				return e
			},
			expected: nil,
		},
		{
			name: "fifty move rule",
			setup: func() *Engine {
				e := NewEngine()
				e.position.Halfmove = 100
				return e
			},
			expected: &Result{Outcome: Draw, Reason: FiftyMoveRule},
		},
		{
			name: "insufficient material - K vs K",
			setup: func() *Engine {
				e := NewEngine()
				e.position = &Position{
					WhiteKing:   bitboardFromStrs([]string{"e1"}),
					BlackKing:   bitboardFromStrs([]string{"e8"}),
					ActiveColor: White,
				}
				return e
			},
			expected: &Result{Outcome: Draw, Reason: InsufficientMaterial},
		},
		{
			name: "insufficient material - K vs B",
			setup: func() *Engine {
				e := NewEngine()
				e.position = &Position{
					WhiteKing:    bitboardFromStrs([]string{"e1"}),
					BlackKing:    bitboardFromStrs([]string{"e8"}),
					BlackBishops: bitboardFromStrs([]string{"f8"}),
					ActiveColor:  White,
				}
				return e
			},
			expected: &Result{Outcome: Draw, Reason: InsufficientMaterial},
		},
		{
			name: "insufficient material - K vs N",
			setup: func() *Engine {
				e := NewEngine()
				e.position = &Position{
					WhiteKing:    bitboardFromStrs([]string{"e1"}),
					BlackKing:    bitboardFromStrs([]string{"e8"}),
					WhiteKnights: bitboardFromStrs([]string{"g1"}),
					ActiveColor:  White,
				}
				return e
			},
			expected: &Result{Outcome: Draw, Reason: InsufficientMaterial},
		},
		{
			name: "insufficient material - B vs B (same color)",
			setup: func() *Engine {
				e := NewEngine()
				e.position = &Position{
					WhiteKing:    bitboardFromStrs([]string{"e1"}),
					BlackKing:    bitboardFromStrs([]string{"e8"}),
					WhiteBishops: bitboardFromStrs([]string{"f1"}),
					BlackBishops: bitboardFromStrs([]string{"c8"}),
					ActiveColor:  White,
				}
				return e
			},
			expected: &Result{Outcome: Draw, Reason: InsufficientMaterial},
		},
		{
			name: "not insufficient material - B vs B (different colors)",
			setup: func() *Engine {
				e := NewEngine()
				e.position = &Position{
					WhiteKing:    bitboardFromStrs([]string{"e1"}),
					BlackKing:    bitboardFromStrs([]string{"e8"}),
					WhiteBishops: bitboardFromStrs([]string{"f1"}),
					BlackBishops: bitboardFromStrs([]string{"f8"}),
					ActiveColor:  White,
				}
				return e
			},
			expected: nil,
		},
		{
			name: "not insufficient material - one pawn remains",
			setup: func() *Engine {
				e := NewEngine()
				e.position = &Position{
					WhiteKing:    bitboardFromStrs([]string{"e1"}),
					BlackKing:    bitboardFromStrs([]string{"e8"}),
					WhiteBishops: bitboardFromStrs([]string{"f1"}),
					WhitePawns:   bitboardFromStrs([]string{"a2"}),
					ActiveColor:  White,
				}
				return e
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine := tt.setup()
			result := engine.GetResult()
			if tt.expected == nil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
