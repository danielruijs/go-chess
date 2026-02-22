package chess

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPiece(t *testing.T) {
	position := NewInitialPosition()

	assert.Equal(t, Piece{Type: Rook, Color: White}, position.GetPiece(strToSquare("a1")))
	assert.Equal(t, Piece{Type: Knight, Color: White}, position.GetPiece(strToSquare("b1")))
	assert.Equal(t, Piece{Type: Bishop, Color: White}, position.GetPiece(strToSquare("c1")))
	assert.Equal(t, Piece{Type: Queen, Color: White}, position.GetPiece(strToSquare("d1")))
	assert.Equal(t, Piece{Type: King, Color: White}, position.GetPiece(strToSquare("e1")))

	assert.Equal(t, Piece{Type: Pawn, Color: White}, position.GetPiece(strToSquare("a2")))
	assert.Equal(t, Piece{Type: Pawn, Color: White}, position.GetPiece(strToSquare("h2")))

	assert.Equal(t, Piece{}, position.GetPiece(strToSquare("a3")))
	assert.Equal(t, Piece{}, position.GetPiece(strToSquare("d4")))
	assert.Equal(t, Piece{}, position.GetPiece(strToSquare("h6")))

	assert.Equal(t, Piece{Type: Pawn, Color: Black}, position.GetPiece(strToSquare("a7")))
	assert.Equal(t, Piece{Type: Pawn, Color: Black}, position.GetPiece(strToSquare("h7")))

	assert.Equal(t, Piece{Type: Queen, Color: Black}, position.GetPiece(strToSquare("d8")))
	assert.Equal(t, Piece{Type: King, Color: Black}, position.GetPiece(strToSquare("e8")))
	assert.Equal(t, Piece{Type: Bishop, Color: Black}, position.GetPiece(strToSquare("f8")))
	assert.Equal(t, Piece{Type: Knight, Color: Black}, position.GetPiece(strToSquare("g8")))
	assert.Equal(t, Piece{Type: Rook, Color: Black}, position.GetPiece(strToSquare("h8")))
}

func TestGetBoard(t *testing.T) {
	position := NewInitialPosition()

	board := position.GetBoard()

	assert.Equal(t, 8, len(board))
	for _, rank := range board {
		assert.Equal(t, 8, len(rank))
	}

	assert.Equal(t, Rook, board[0][0].Type)
	assert.Equal(t, White, board[0][0].Color)
	assert.Equal(t, Queen, board[3][0].Type)
	assert.Equal(t, White, board[3][0].Color)
	assert.Equal(t, King, board[4][7].Type)
	assert.Equal(t, Black, board[4][7].Color)
	assert.Equal(t, Pawn, board[2][6].Type)
	assert.Equal(t, Black, board[2][6].Color)
}

func TestGetCopy(t *testing.T) {
	position := NewInitialPosition()
	position.Halfmove = 5
	position.Fullmove = 10
	enPassantSquare := Square(16)
	position.EnPassant = &enPassantSquare

	copy := position.GetCopy()

	assert.Equal(t, position.WhitePawns, copy.WhitePawns)
	assert.Equal(t, position.BlackPawns, copy.BlackPawns)

	assert.Equal(t, position.ActiveColor, copy.ActiveColor)
	assert.Equal(t, position.CastlingRights, copy.CastlingRights)
	assert.Equal(t, position.Halfmove, copy.Halfmove)
	assert.Equal(t, position.Fullmove, copy.Fullmove)

	assert.NotNil(t, copy.EnPassant)
	assert.Equal(t, *position.EnPassant, *copy.EnPassant)
	assert.NotSame(t, position.EnPassant, copy.EnPassant)

	// Changing en passant should not change original
	*copy.EnPassant = 20
	assert.Equal(t, Square(16), *position.EnPassant)

	// Changing a piece should not change original
	copy.WhitePawns |= squareMask(24)
	assert.Equal(t, Bitboard(1<<24), copy.WhitePawns&squareMask(24))
	assert.Equal(t, Bitboard(0), position.WhitePawns&squareMask(24))
}

func TestGetCopyWithNilEnPassant(t *testing.T) {
	position := NewInitialPosition()

	copy := position.GetCopy()

	assert.Nil(t, copy.EnPassant)
}

func TestSetPieceWhite(t *testing.T) {
	position := NewInitialPosition()

	position.setPiece(strToSquare("e4"), Piece{Type: Pawn, Color: White})

	assert.Equal(t, Piece{Type: Pawn, Color: White}, position.GetPiece(strToSquare("e4")))
}

func TestSetPieceBlack(t *testing.T) {
	position := NewInitialPosition()

	position.setPiece(strToSquare("e5"), Piece{Type: Pawn, Color: Black})

	assert.Equal(t, Piece{Type: Pawn, Color: Black}, position.GetPiece(strToSquare("e5")))
}

func TestSetPieceMultiple(t *testing.T) {
	position := NewInitialPosition()

	position.setPiece(strToSquare("e4"), Piece{Type: Pawn, Color: White})
	position.setPiece(strToSquare("d4"), Piece{Type: Knight, Color: White})
	position.setPiece(strToSquare("f4"), Piece{Type: Bishop, Color: White})

	assert.Equal(t, Piece{Type: Pawn, Color: White}, position.GetPiece(strToSquare("e4")))
	assert.Equal(t, Piece{Type: Knight, Color: White}, position.GetPiece(strToSquare("d4")))
	assert.Equal(t, Piece{Type: Bishop, Color: White}, position.GetPiece(strToSquare("f4")))
}

func TestRemovePiece(t *testing.T) {
	position := NewInitialPosition()

	position.removePiece(strToSquare("e2"))

	assert.Equal(t, Piece{}, position.GetPiece(strToSquare("e2")))
}

func TestRemovePieceMultiplePieces(t *testing.T) {
	position := NewInitialPosition()

	position.removePiece(strToSquare("a2"))
	position.removePiece(strToSquare("b2"))

	assert.Equal(t, Piece{}, position.GetPiece(strToSquare("a2")))
	assert.Equal(t, Piece{}, position.GetPiece(strToSquare("b2")))
	assert.Equal(t, Piece{Type: Pawn, Color: White}, position.GetPiece(strToSquare("c2")))
}

func TestMakeMoveRegular(t *testing.T) {
	position := NewInitialPosition()
	move := Move{From: strToSquare("e2"), To: strToSquare("e4")}

	position.MakeMove(move)

	assert.Equal(t, Piece{}, position.GetPiece(strToSquare("e2")))
	assert.Equal(t, Piece{Type: Pawn, Color: White}, position.GetPiece(strToSquare("e4")))
}

func TestMakeMoveCapture(t *testing.T) {
	position := &Position{
		WhitePawns: bitboardFromStrs([]string{"e4"}),
		BlackPawns: bitboardFromStrs([]string{"f5"}),
	}

	move := Move{From: strToSquare("e4"), To: strToSquare("f5")}
	position.MakeMove(move)

	assert.Equal(t, Piece{}, position.GetPiece(strToSquare("e4")))
	assert.Equal(t, Piece{Type: Pawn, Color: White}, position.GetPiece(strToSquare("f5")))
}

func TestMakeMovePromotion(t *testing.T) {
	position := &Position{
		WhitePawns: bitboardFromStrs([]string{"e7"}),
	}
	position.setPiece(strToSquare("e7"), Piece{Type: Pawn, Color: White})

	move := Move{From: strToSquare("e7"), To: strToSquare("e8"), Promotion: new(Queen)}
	position.MakeMove(move)

	assert.Equal(t, Piece{}, position.GetPiece(strToSquare("e7")))
	assert.Equal(t, Piece{Type: Queen, Color: White}, position.GetPiece(strToSquare("e8")))
}

func TestMakeMoveEnPassant(t *testing.T) {
	position := &Position{
		WhitePawns: bitboardFromStrs([]string{"a5"}),
		BlackPawns: bitboardFromStrs([]string{"b5"}),
		EnPassant:  new(strToSquare("b6")),
	}

	move := Move{From: strToSquare("a5"), To: strToSquare("b6")}
	position.MakeMove(move)

	assert.Equal(t, Piece{}, position.GetPiece(strToSquare("a5")))
	assert.Equal(t, Piece{Type: Pawn, Color: White}, position.GetPiece(strToSquare("b6")))
	assert.Equal(t, Piece{}, position.GetPiece(strToSquare("b5")))
}

func TestMakeMoveKingSideCastling(t *testing.T) {
	position := &Position{
		WhiteKing:  bitboardFromStrs([]string{"e1"}),
		WhiteRooks: bitboardFromStrs([]string{"h1"}),
	}

	move := Move{From: strToSquare("e1"), To: strToSquare("g1")}
	position.MakeMove(move)

	assert.Equal(t, Piece{Type: King, Color: White}, position.GetPiece(strToSquare("g1")))
	assert.Equal(t, Piece{Type: Rook, Color: White}, position.GetPiece(strToSquare("f1")))
	assert.Equal(t, Piece{}, position.GetPiece(strToSquare("e1")))
	assert.Equal(t, Piece{}, position.GetPiece(strToSquare("h1")))
}

func TestMakeMoveQueenSideCastling(t *testing.T) {
	position := &Position{
		BlackKing:  bitboardFromStrs([]string{"e8"}),
		BlackRooks: bitboardFromStrs([]string{"a8"}),
	}

	move := Move{From: strToSquare("e8"), To: strToSquare("c8")}
	position.MakeMove(move)

	assert.Equal(t, Piece{Type: King, Color: Black}, position.GetPiece(strToSquare("c8")))
	assert.Equal(t, Piece{Type: Rook, Color: Black}, position.GetPiece(strToSquare("d8")))
	assert.Equal(t, Piece{}, position.GetPiece(strToSquare("e8")))
	assert.Equal(t, Piece{}, position.GetPiece(strToSquare("a8")))
}

func TestMakeMoveUpdatesCastlingRightsKing(t *testing.T) {
	position := &Position{
		WhiteKing: bitboardFromStrs([]string{"e1"}),
		CastlingRights: CastlingRights{
			WhiteOO:  true,
			WhiteOOO: true,
			BlackOO:  true,
			BlackOOO: true,
		},
	}

	move := Move{From: strToSquare("e1"), To: strToSquare("e2")}
	position.MakeMove(move)

	assert.False(t, position.CastlingRights.WhiteOO)
	assert.False(t, position.CastlingRights.WhiteOOO)
	assert.True(t, position.CastlingRights.BlackOO)
	assert.True(t, position.CastlingRights.BlackOOO)
}

func TestMakeMoveUpdatesCastlingRightsRook(t *testing.T) {
	position := &Position{
		WhiteRooks: bitboardFromStrs([]string{"h1"}),
		CastlingRights: CastlingRights{
			WhiteOO:  true,
			WhiteOOO: true,
			BlackOO:  true,
			BlackOOO: true,
		},
	}

	move := Move{From: strToSquare("h1"), To: strToSquare("h2")}
	position.MakeMove(move)

	assert.False(t, position.CastlingRights.WhiteOO)
	assert.True(t, position.CastlingRights.WhiteOOO)
	assert.True(t, position.CastlingRights.BlackOO)
	assert.True(t, position.CastlingRights.BlackOOO)
}

func TestMakeMoveUpdatesEnPassantSquare(t *testing.T) {
	position := &Position{
		WhitePawns: bitboardFromStrs([]string{"e2"}),
	}

	move := Move{From: strToSquare("e2"), To: strToSquare("e4")}
	position.MakeMove(move)

	assert.NotNil(t, position.EnPassant)
	assert.Equal(t, strToSquare("e3"), *position.EnPassant)
}
