package chess

import "math/bits"

const (
	rank1Mask Bitboard = 0x00000000000000FF
	rank2Mask Bitboard = 0x000000000000FF00
	rank7Mask Bitboard = 0x00FF000000000000
	rank8Mask Bitboard = 0xFF00000000000000

	fileAMask Bitboard = 0x0101010101010101
	fileHMask Bitboard = 0x8080808080808080
)

type Generator interface {
	generateMoves(pos *Position, color Color) []Move
}

var pieceToGenerator = map[PieceType]Generator{
	Pawn:   pawnMoveGenerator{},
	Knight: knightMoveGenerator{},
	Bishop: bishopMoveGenerator{},
	Rook:   rookMoveGenerator{},
	Queen:  queenMoveGenerator{},
	King:   kingMoveGenerator{},
}

// generates all legal moves in a position for a given color
func GenerateMoves(pos *Position, color Color) []Move {
	// TODO
	return []Move{}
}

// generates all legal moves in a position for a given color and piece type
func GeneratePieceMoves(pos *Position, color Color, pieceType PieceType) []Move {
	generator := pieceToGenerator[pieceType]
	return generator.generateMoves(pos, color)
}

type pawnMoveGenerator struct{}

func (g pawnMoveGenerator) generateMoves(pos *Position, color Color) []Move {
	moves := []Move{}

	var pawns, startingRankMask, promotionRankMask Bitboard
	var forwardOffset int
	if color == White {
		pawns = pos.WhitePawns
		startingRankMask = rank2Mask
		promotionRankMask = rank8Mask
		forwardOffset = 8
	} else {
		pawns = pos.BlackPawns
		startingRankMask = rank7Mask
		promotionRankMask = rank1Mask
		forwardOffset = -8
	}

	occupied := pos.GetOccupied()
	empty := ^occupied

	// regular forward move by 1 if not blocked
	singlePushes := shift(pawns, forwardOffset) & empty & ^promotionRankMask
	for sq := singlePushes; sq != 0; {
		to := bits.TrailingZeros64(uint64(sq)) // index of LSB
		from := to - forwardOffset
		moves = append(moves, Move{Piece: Pawn, From: Square(from), To: Square(to)})
		sq &= sq - 1 // clear LSB
	}

	// move forward by 2 if on starting rank and not blocked
	doublePushes := shift((pawns&startingRankMask), 2*forwardOffset) & empty
	for sq := doublePushes; sq != 0; {
		to := bits.TrailingZeros64(uint64(sq)) // index of LSB
		from := to - 2*forwardOffset
		moves = append(moves, Move{Piece: Pawn, From: Square(from), To: Square(to)})
		sq &= sq - 1 // clear LSB
	}

	// capture diagonally by 1 if occupied by opponent piece
	capturesLeft := shift(pawns, forwardOffset+1) & occupied & ^fileAMask
	capturesRight := shift(pawns, forwardOffset-1) & occupied & ^fileHMask

	// en passant
	if pos.EnPassant != nil {
		enPassantMask := squareMask(*pos.EnPassant)
		epCapturesLeft := shift(pawns, forwardOffset+1) & enPassantMask & ^fileAMask
		epCapturesRight := shift(pawns, forwardOffset-1) & enPassantMask & ^fileHMask

		capturesLeft |= epCapturesLeft
		capturesRight |= epCapturesRight
	}
	for sq := capturesLeft; sq != 0; {
		to := bits.TrailingZeros64(uint64(sq))
		from := to - forwardOffset - 1
		moves = append(moves, Move{Piece: Pawn, From: Square(from), To: Square(to)})
		sq &= sq - 1
	}
	for sq := capturesRight; sq != 0; {
		to := bits.TrailingZeros64(uint64(sq))
		from := to - forwardOffset + 1
		moves = append(moves, Move{Piece: Pawn, From: Square(from), To: Square(to)})
		sq &= sq - 1
	}

	// promotion
	promotions := shift(pawns, forwardOffset) & empty & promotionRankMask
	for sq := promotions; sq != 0; {
		to := bits.TrailingZeros64(uint64(sq))
		from := to - forwardOffset
		for _, promotionPiece := range []PieceType{Queen, Rook, Bishop, Knight} {
			moves = append(moves, Move{Piece: Pawn, From: Square(from), To: Square(to), Promotion: &promotionPiece})
		}
		sq &= sq - 1
	}

	return moves
}

func shift(b Bitboard, n int) Bitboard {
	if n > 0 {
		return b << n
	} else {
		return b >> (-n)
	}
}

type knightMoveGenerator struct{}

func (g knightMoveGenerator) generateMoves(pos *Position, color Color) []Move {
	// TODO
	return []Move{}
}

type bishopMoveGenerator struct{}

func (g bishopMoveGenerator) generateMoves(pos *Position, color Color) []Move {
	directions := []int{7, 9, -7, -9} // NW, NE, SE, SW

	var bishops, ownPieces Bitboard
	if color == White {
		bishops = pos.WhiteBishops
		ownPieces = pos.GetOccupiedWhite()
	} else {
		bishops = pos.BlackBishops
		ownPieces = pos.GetOccupiedBlack()
	}

	occupied := pos.GetOccupied()

	moves := []Move{}
	for sq := bishops; sq != 0; {
		from := bits.TrailingZeros64(uint64(sq)) // index of LSB
		for _, dir := range directions {
			to := from
			for {
				prevFile := to % 8
				to += dir
				toMask := squareMask(Square(to))
				newFile := to % 8

				// edges of the board
				fileDiff := newFile - prevFile
				if to < 0 || to > 63 || (fileDiff != 1 && fileDiff != -1) {
					break
				}

				// blocked by own piece
				if toMask&ownPieces != 0 {
					break
				}
				moves = append(moves, Move{Piece: Bishop, From: Square(from), To: Square(to)})
				// blocked by opponent piece, added capture
				if toMask&occupied != 0 {
					break
				}
			}
		}
		sq &= sq - 1 // clear LSB
	}

	return moves
}

type rookMoveGenerator struct{}

func (g rookMoveGenerator) generateMoves(pos *Position, color Color) []Move {
	// TODO
	return []Move{}
}

type queenMoveGenerator struct{}

func (g queenMoveGenerator) generateMoves(pos *Position, color Color) []Move {
	// TODO
	return []Move{}
}

type kingMoveGenerator struct{}

func (g kingMoveGenerator) generateMoves(pos *Position, color Color) []Move {
	// TODO
	return []Move{}
}
