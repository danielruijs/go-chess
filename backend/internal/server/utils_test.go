package server

import (
	"go-chess/internal/chess"
	"testing"

	"github.com/stretchr/testify/assert"
)

// unsafe version of StrToSquare, only for tests
func strToSquare(str string) chess.Square {
	file := int(str[0] - 'a')
	rank := int(str[1] - '1')
	return chess.Square(file + rank*8)
}

func TestMoveDataToMove(t *testing.T) {
	tests := []struct {
		name      string
		data      MoveData
		expected  chess.Move
		expectErr bool
	}{
		{
			name: "valid move without promotion",
			data: MoveData{
				From:      "e2",
				To:        "e4",
				Promotion: nil,
			},
			expected: chess.Move{
				From:      strToSquare("e2"),
				To:        strToSquare("e4"),
				Promotion: nil,
			},
			expectErr: false,
		},
		{
			name: "valid move with promotion",
			data: MoveData{
				From:      "e7",
				To:        "e8",
				Promotion: new(chess.Queen),
			},
			expected: chess.Move{
				From:      strToSquare("e7"),
				To:        strToSquare("e8"),
				Promotion: new(chess.Queen),
			},
			expectErr: false,
		},
		{
			name: "invalid from square",
			data: MoveData{
				From:      "z9",
				To:        "e4",
				Promotion: nil,
			},
			expectErr: true,
		},
		{
			name: "invalid to square",
			data: MoveData{
				From:      "e2",
				To:        "z9",
				Promotion: nil,
			},
			expectErr: true,
		},
		{
			name: "both squares invalid",
			data: MoveData{
				From:      "invalid",
				To:        "invalid",
				Promotion: nil,
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := moveDataToMove(tt.data)
			if tt.expectErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestMoveListToLegalMoves(t *testing.T) {
	tests := []struct {
		name     string
		moves    []chess.Move
		expected map[string][]LegalMove
	}{
		{
			name:     "empty move list",
			moves:    []chess.Move{},
			expected: map[string][]LegalMove{},
		},
		{
			name: "single move",
			moves: []chess.Move{
				{
					From:      strToSquare("a2"),
					To:        strToSquare("a4"),
					Promotion: nil,
				},
			},
			expected: map[string][]LegalMove{
				"a2": {
					{
						To:        "a4",
						Promotion: nil,
					},
				},
			},
		},
		{
			name: "multiple moves from same square",
			moves: []chess.Move{
				{
					From:      strToSquare("e2"),
					To:        strToSquare("e3"),
					Promotion: nil,
				},
				{
					From:      strToSquare("e2"),
					To:        strToSquare("e4"),
					Promotion: nil,
				},
			},
			expected: map[string][]LegalMove{
				"e2": {
					{
						To:        "e3",
						Promotion: nil,
					},
					{
						To:        "e4",
						Promotion: nil,
					},
				},
			},
		},
		{
			name: "moves from different squares with promotion",
			moves: []chess.Move{
				{
					From:      strToSquare("a7"),
					To:        strToSquare("a8"),
					Promotion: new(chess.Queen),
				},
				{
					From:      strToSquare("h7"),
					To:        strToSquare("h8"),
					Promotion: new(chess.Knight),
				},
			},
			expected: map[string][]LegalMove{
				"a7": {
					{
						To:        "a8",
						Promotion: new(chess.Queen),
					},
				},
				"h7": {
					{
						To:        "h8",
						Promotion: new(chess.Knight),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := moveListToLegalMoves(tt.moves)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGenerateMatchID(t *testing.T) {
	id1, err := GenerateMatchID()
	assert.NoError(t, err)
	assert.Len(t, id1, 12)

	id2, err := GenerateMatchID()
	assert.NoError(t, err)
	assert.NotEqual(t, id1, id2)
}
