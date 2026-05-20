package chess

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPiece(t *testing.T) {
	pos := NewInitialPosition()

	assert.Equal(t, Piece{Type: Rook, Color: White}, *pos.GetPiece(strToSquare("a1")))
	assert.Equal(t, Piece{Type: Knight, Color: White}, *pos.GetPiece(strToSquare("b1")))
	assert.Equal(t, Piece{Type: Bishop, Color: White}, *pos.GetPiece(strToSquare("c1")))
	assert.Equal(t, Piece{Type: Queen, Color: White}, *pos.GetPiece(strToSquare("d1")))
	assert.Equal(t, Piece{Type: King, Color: White}, *pos.GetPiece(strToSquare("e1")))

	assert.Equal(t, Piece{Type: Pawn, Color: White}, *pos.GetPiece(strToSquare("a2")))
	assert.Equal(t, Piece{Type: Pawn, Color: White}, *pos.GetPiece(strToSquare("h2")))

	assert.Nil(t, pos.GetPiece(strToSquare("a3")))
	assert.Nil(t, pos.GetPiece(strToSquare("d4")))
	assert.Nil(t, pos.GetPiece(strToSquare("h6")))

	assert.Equal(t, Piece{Type: Pawn, Color: Black}, *pos.GetPiece(strToSquare("a7")))
	assert.Equal(t, Piece{Type: Pawn, Color: Black}, *pos.GetPiece(strToSquare("h7")))

	assert.Equal(t, Piece{Type: Queen, Color: Black}, *pos.GetPiece(strToSquare("d8")))
	assert.Equal(t, Piece{Type: King, Color: Black}, *pos.GetPiece(strToSquare("e8")))
	assert.Equal(t, Piece{Type: Bishop, Color: Black}, *pos.GetPiece(strToSquare("f8")))
	assert.Equal(t, Piece{Type: Knight, Color: Black}, *pos.GetPiece(strToSquare("g8")))
	assert.Equal(t, Piece{Type: Rook, Color: Black}, *pos.GetPiece(strToSquare("h8")))
}

func TestGetBoard(t *testing.T) {
	pos := NewInitialPosition()

	board := pos.GetBoard()

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
	pos := NewInitialPosition()
	pos.Halfmove = 5
	pos.Fullmove = 10
	enPassantSquare := Square(16)
	pos.EnPassant = &enPassantSquare

	copy := pos.GetCopy()

	assert.Equal(t, pos.WhitePawns, copy.WhitePawns)
	assert.Equal(t, pos.BlackPawns, copy.BlackPawns)

	assert.Equal(t, pos.ActiveColor, copy.ActiveColor)
	assert.Equal(t, pos.CastlingRights, copy.CastlingRights)
	assert.Equal(t, pos.Halfmove, copy.Halfmove)
	assert.Equal(t, pos.Fullmove, copy.Fullmove)

	assert.NotNil(t, copy.EnPassant)
	assert.Equal(t, *pos.EnPassant, *copy.EnPassant)
	assert.NotSame(t, pos.EnPassant, copy.EnPassant)

	// Changing en passant should not change original
	*copy.EnPassant = 20
	assert.Equal(t, Square(16), *pos.EnPassant)

	// Changing a piece should not change original
	copy.WhitePawns |= squareMask(24)
	assert.Equal(t, Bitboard(1<<24), copy.WhitePawns&squareMask(24))
	assert.Equal(t, Bitboard(0), pos.WhitePawns&squareMask(24))
}

func TestGetCopyWithNilEnPassant(t *testing.T) {
	pos := NewInitialPosition()

	copy := pos.GetCopy()

	assert.Nil(t, copy.EnPassant)
}

func TestSetPieceWhite(t *testing.T) {
	pos := NewInitialPosition()

	pos.setPiece(strToSquare("e4"), &Piece{Type: Pawn, Color: White})

	assert.Equal(t, Piece{Type: Pawn, Color: White}, *pos.GetPiece(strToSquare("e4")))
}

func TestSetPieceBlack(t *testing.T) {
	pos := NewInitialPosition()

	pos.setPiece(strToSquare("e5"), &Piece{Type: Pawn, Color: Black})

	assert.Equal(t, Piece{Type: Pawn, Color: Black}, *pos.GetPiece(strToSquare("e5")))
}

func TestSetPieceMultiple(t *testing.T) {
	pos := NewInitialPosition()

	pos.setPiece(strToSquare("e4"), &Piece{Type: Pawn, Color: White})
	pos.setPiece(strToSquare("d4"), &Piece{Type: Knight, Color: White})
	pos.setPiece(strToSquare("f4"), &Piece{Type: Bishop, Color: White})

	assert.Equal(t, Piece{Type: Pawn, Color: White}, *pos.GetPiece(strToSquare("e4")))
	assert.Equal(t, Piece{Type: Knight, Color: White}, *pos.GetPiece(strToSquare("d4")))
	assert.Equal(t, Piece{Type: Bishop, Color: White}, *pos.GetPiece(strToSquare("f4")))
}

func TestRemovePiece(t *testing.T) {
	pos := NewInitialPosition()

	pos.removePiece(strToSquare("e2"))

	assert.Nil(t, pos.GetPiece(strToSquare("e2")))
}

func TestRemovePieceMultiplePieces(t *testing.T) {
	pos := NewInitialPosition()

	pos.removePiece(strToSquare("a2"))
	pos.removePiece(strToSquare("b2"))

	assert.Nil(t, pos.GetPiece(strToSquare("a2")))
	assert.Nil(t, pos.GetPiece(strToSquare("b2")))
	assert.Equal(t, Piece{Type: Pawn, Color: White}, *pos.GetPiece(strToSquare("c2")))
}

func TestMakeMoveRegular(t *testing.T) {
	pos := NewInitialPosition()
	move := Move{From: strToSquare("e2"), To: strToSquare("e4")}

	pos.MakeMove(move)

	assert.Nil(t, pos.GetPiece(strToSquare("e2")))
	assert.Equal(t, Piece{Type: Pawn, Color: White}, *pos.GetPiece(strToSquare("e4")))
	assert.Equal(t, uint(0), pos.Halfmove)
	assert.Equal(t, uint(1), pos.Fullmove)
	assert.Equal(t, Black, pos.ActiveColor)
}

func TestMakeMoveCapture(t *testing.T) {
	pos := &Position{
		WhitePawns:  bitboardFromStrs([]string{"e4"}),
		BlackPawns:  bitboardFromStrs([]string{"f5"}),
		Halfmove:    10,
		Fullmove:    3,
		ActiveColor: White,
	}

	move := Move{From: strToSquare("e4"), To: strToSquare("f5")}
	pos.MakeMove(move)

	assert.Nil(t, pos.GetPiece(strToSquare("e4")))
	assert.Equal(t, Piece{Type: Pawn, Color: White}, *pos.GetPiece(strToSquare("f5")))
	assert.Equal(t, uint(0), pos.Halfmove)
	assert.Equal(t, uint(3), pos.Fullmove)
	assert.Equal(t, Black, pos.ActiveColor)
}

func TestMakeMovePromotion(t *testing.T) {
	pos := &Position{
		WhitePawns:  bitboardFromStrs([]string{"e7"}),
		Halfmove:    10,
		Fullmove:    3,
		ActiveColor: White,
	}

	move := Move{From: strToSquare("e7"), To: strToSquare("e8"), Promotion: new(Queen)}
	pos.MakeMove(move)

	assert.Nil(t, pos.GetPiece(strToSquare("e7")))
	assert.Equal(t, Piece{Type: Queen, Color: White}, *pos.GetPiece(strToSquare("e8")))
	assert.Equal(t, uint(0), pos.Halfmove)
	assert.Equal(t, uint(3), pos.Fullmove)
	assert.Equal(t, Black, pos.ActiveColor)
}

func TestMakeMovePromotionCapture(t *testing.T) {
	pos := &Position{
		WhitePawns:  bitboardFromStrs([]string{"e7"}),
		BlackPawns:  bitboardFromStrs([]string{"d8"}),
		Halfmove:    10,
		Fullmove:    3,
		ActiveColor: White,
	}

	move := Move{From: strToSquare("e7"), To: strToSquare("d8"), Promotion: new(Queen)}
	pos.MakeMove(move)

	assert.Nil(t, pos.GetPiece(strToSquare("e7")))
	assert.Equal(t, Piece{Type: Queen, Color: White}, *pos.GetPiece(strToSquare("d8")))
	assert.Equal(t, uint(0), pos.Halfmove)
	assert.Equal(t, uint(3), pos.Fullmove)
	assert.Equal(t, Black, pos.ActiveColor)
}

func TestMakeMoveEnPassant(t *testing.T) {
	pos := &Position{
		WhitePawns:  bitboardFromStrs([]string{"a5"}),
		BlackPawns:  bitboardFromStrs([]string{"b5"}),
		EnPassant:   new(strToSquare("b6")),
		Halfmove:    10,
		Fullmove:    3,
		ActiveColor: White,
	}

	move := Move{From: strToSquare("a5"), To: strToSquare("b6")}
	pos.MakeMove(move)

	assert.Nil(t, pos.GetPiece(strToSquare("a5")))
	assert.Equal(t, Piece{Type: Pawn, Color: White}, *pos.GetPiece(strToSquare("b6")))
	assert.Nil(t, pos.GetPiece(strToSquare("b5")))
	assert.Equal(t, uint(0), pos.Halfmove)
	assert.Equal(t, uint(3), pos.Fullmove)
	assert.Equal(t, Black, pos.ActiveColor)
}

func TestMakeMoveKingSideCastling(t *testing.T) {
	pos := &Position{
		WhiteKing:   bitboardFromStrs([]string{"e1"}),
		WhiteRooks:  bitboardFromStrs([]string{"h1"}),
		Halfmove:    10,
		Fullmove:    3,
		ActiveColor: White,
	}

	move := Move{From: strToSquare("e1"), To: strToSquare("g1")}
	pos.MakeMove(move)

	assert.Equal(t, Piece{Type: King, Color: White}, *pos.GetPiece(strToSquare("g1")))
	assert.Equal(t, Piece{Type: Rook, Color: White}, *pos.GetPiece(strToSquare("f1")))
	assert.Nil(t, pos.GetPiece(strToSquare("e1")))
	assert.Nil(t, pos.GetPiece(strToSquare("h1")))
	assert.Equal(t, uint(11), pos.Halfmove)
	assert.Equal(t, uint(3), pos.Fullmove)
	assert.Equal(t, Black, pos.ActiveColor)
}

func TestMakeMoveQueenSideCastling(t *testing.T) {
	pos := &Position{
		BlackKing:   bitboardFromStrs([]string{"e8"}),
		BlackRooks:  bitboardFromStrs([]string{"a8"}),
		Halfmove:    10,
		Fullmove:    3,
		ActiveColor: Black,
	}

	move := Move{From: strToSquare("e8"), To: strToSquare("c8")}
	pos.MakeMove(move)

	assert.Equal(t, Piece{Type: King, Color: Black}, *pos.GetPiece(strToSquare("c8")))
	assert.Equal(t, Piece{Type: Rook, Color: Black}, *pos.GetPiece(strToSquare("d8")))
	assert.Nil(t, pos.GetPiece(strToSquare("e8")))
	assert.Nil(t, pos.GetPiece(strToSquare("a8")))
	assert.Equal(t, uint(11), pos.Halfmove)
	assert.Equal(t, uint(4), pos.Fullmove)
	assert.Equal(t, White, pos.ActiveColor)
}

func TestMakeMoveUpdatesCastlingRightsKing(t *testing.T) {
	pos := &Position{
		WhiteKing: bitboardFromStrs([]string{"e1"}),
		CastlingRights: CastlingRights{
			WhiteOO:  true,
			WhiteOOO: true,
			BlackOO:  true,
			BlackOOO: true,
		},
		Halfmove:    10,
		Fullmove:    3,
		ActiveColor: White,
	}

	move := Move{From: strToSquare("e1"), To: strToSquare("e2")}
	pos.MakeMove(move)

	assert.False(t, pos.CastlingRights.WhiteOO)
	assert.False(t, pos.CastlingRights.WhiteOOO)
	assert.True(t, pos.CastlingRights.BlackOO)
	assert.True(t, pos.CastlingRights.BlackOOO)
	assert.Equal(t, uint(11), pos.Halfmove)
	assert.Equal(t, uint(3), pos.Fullmove)
	assert.Equal(t, Black, pos.ActiveColor)
}

func TestMakeMoveUpdatesCastlingRightsRook(t *testing.T) {
	pos := &Position{
		WhiteRooks: bitboardFromStrs([]string{"h1"}),
		CastlingRights: CastlingRights{
			WhiteOO:  true,
			WhiteOOO: true,
			BlackOO:  true,
			BlackOOO: true,
		},
		Halfmove:    10,
		Fullmove:    3,
		ActiveColor: White,
	}

	move := Move{From: strToSquare("h1"), To: strToSquare("h2")}
	pos.MakeMove(move)

	assert.False(t, pos.CastlingRights.WhiteOO)
	assert.True(t, pos.CastlingRights.WhiteOOO)
	assert.True(t, pos.CastlingRights.BlackOO)
	assert.True(t, pos.CastlingRights.BlackOOO)
	assert.Equal(t, uint(11), pos.Halfmove)
	assert.Equal(t, uint(3), pos.Fullmove)
	assert.Equal(t, Black, pos.ActiveColor)
}

func TestCastlingNotRestoredWhenRookMovesToOriginalSquareWhite(t *testing.T) {
	pos := &Position{
		WhiteKing:    bitboardFromStrs([]string{"e1"}),
		WhiteRooks:   bitboardFromStrs([]string{"h1", "h5"}),
		BlackKing:    bitboardFromStrs([]string{"e8"}),
		BlackBishops: bitboardFromStrs([]string{"b7"}),
		CastlingRights: CastlingRights{
			WhiteOO: true,
		},
		ActiveColor: Black,
	}

	// Black captures white rook on h1
	pos.MakeMove(Move{From: strToSquare("b7"), To: strToSquare("h1")})
	// White captures black bishop on h1 with second rook
	pos.MakeMove(Move{From: strToSquare("h5"), To: strToSquare("h1")})

	moves := NewGenerator().generateKingMoves(pos, White)
	assert.NotContains(t, moves, Move{From: strToSquare("e1"), To: strToSquare("g1")})
}

func TestCastlingNotRestoredWhenRookMovesToOriginalSquareBlack(t *testing.T) {
	pos := &Position{
		WhiteKing:    bitboardFromStrs([]string{"e1"}),
		WhiteBishops: bitboardFromStrs([]string{"g2"}),
		BlackKing:    bitboardFromStrs([]string{"e8"}),
		BlackRooks:   bitboardFromStrs([]string{"a8", "a5"}),
		CastlingRights: CastlingRights{
			BlackOOO: true,
		},
		ActiveColor: White,
	}

	// White captures black rook on a8
	pos.MakeMove(Move{From: strToSquare("g2"), To: strToSquare("a8")})
	// White captures black bishop on a8 with second rook
	pos.MakeMove(Move{From: strToSquare("a5"), To: strToSquare("a8")})

	moves := NewGenerator().generateKingMoves(pos, Black)
	assert.NotContains(t, moves, Move{From: strToSquare("e8"), To: strToSquare("c8")})
}

func TestMakeMoveUpdatesEnPassantSquare(t *testing.T) {
	pos := &Position{
		WhitePawns:  bitboardFromStrs([]string{"e2"}),
		Halfmove:    10,
		Fullmove:    3,
		ActiveColor: White,
	}

	move := Move{From: strToSquare("e2"), To: strToSquare("e4")}
	pos.MakeMove(move)

	assert.NotNil(t, pos.EnPassant)
	assert.Equal(t, strToSquare("e3"), *pos.EnPassant)
	assert.Equal(t, uint(0), pos.Halfmove)
	assert.Equal(t, uint(3), pos.Fullmove)
	assert.Equal(t, Black, pos.ActiveColor)
}
