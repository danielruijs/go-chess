package chess

const (
	rank1Mask Bitboard = 0x00000000000000FF
	rank2Mask Bitboard = 0x000000000000FF00
	rank7Mask Bitboard = 0x00FF000000000000
	rank8Mask Bitboard = 0xFF00000000000000

	fileAMask Bitboard = 0x0101010101010101
	fileHMask Bitboard = 0x8080808080808080
)

type Generator struct {
	knightMoves [64]Bitboard
	// kingMoves   [64]Bitboard

	pieceToGenerator map[PieceType]func(pos *Position, color Color) []Move
}

func NewGenerator() *Generator {
	g := &Generator{}

	g.knightMoves = precomputeKnightMoves()

	g.pieceToGenerator = map[PieceType]func(pos *Position, color Color) []Move{
		Pawn:   g.generatePawnMoves,
		Knight: g.generateKnightMoves,
		Bishop: g.generateBishopMoves,
		Rook:   g.generateRookMoves,
		Queen:  g.generateQueenMoves,
		King:   g.generateKingMoves,
	}

	return g
}

// Generates all legal moves in a position for a given color
func (g *Generator) GenerateMoves(pos *Position, color Color) []Move {
	// TODO
	return []Move{}
}

// Generates all legal moves in a position for a given color and piece type
func (g *Generator) GeneratePieceMoves(pos *Position, color Color, pieceType PieceType) []Move {
	generator := g.pieceToGenerator[pieceType]
	psuedoLegalMoves := generator(pos, color)
	return psuedoLegalMoves // TODO check legality
}

func (g *Generator) generatePawnMoves(pos *Position, color Color) []Move {
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
	for singlePushes != 0 {
		to := popLSB(&singlePushes)
		from := to - forwardOffset
		moves = append(moves, Move{Piece: Pawn, From: Square(from), To: Square(to)})
	}

	// move forward by 2 if on starting rank and not blocked
	doublePushes := shift((pawns&startingRankMask), 2*forwardOffset) & empty
	for doublePushes != 0 {
		to := popLSB(&doublePushes)
		from := to - 2*forwardOffset
		moves = append(moves, Move{Piece: Pawn, From: Square(from), To: Square(to)})
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
	for capturesLeft != 0 {
		to := popLSB(&capturesLeft)
		from := to - forwardOffset - 1
		moves = append(moves, Move{Piece: Pawn, From: Square(from), To: Square(to)})
	}
	for capturesRight != 0 {
		to := popLSB(&capturesRight)
		from := to - forwardOffset + 1
		moves = append(moves, Move{Piece: Pawn, From: Square(from), To: Square(to)})
	}

	// promotion
	promotions := shift(pawns, forwardOffset) & empty & promotionRankMask
	for promotions != 0 {
		to := popLSB(&promotions)
		from := to - forwardOffset
		for _, promotionPiece := range []PieceType{Queen, Rook, Bishop, Knight} {
			moves = append(moves, Move{Piece: Pawn, From: Square(from), To: Square(to), Promotion: &promotionPiece})
		}
	}

	return moves
}

func precomputeKnightMoves() [64]Bitboard {
	var knightMoves [64]Bitboard

	offsets := []struct{ rank, file int }{
		{2, 1}, {2, -1}, {-2, 1}, {-2, -1},
		{1, 2}, {1, -2}, {-1, 2}, {-1, -2},
	}

	for sq := 0; sq < 64; sq++ {
		var moves Bitboard
		rank := sq / 8
		file := sq % 8
		for _, offset := range offsets {
			newRank := rank + offset.rank
			newFile := file + offset.file
			if newRank >= 0 && newRank < 8 && newFile >= 0 && newFile < 8 {
				moves |= coordMask(newFile, newRank)
			}
		}
		knightMoves[sq] = moves
	}

	return knightMoves
}

func (g *Generator) generateKnightMoves(pos *Position, color Color) []Move {
	moves := []Move{}

	var knights, ownPieces Bitboard
	if color == White {
		knights = pos.WhiteKnights
		ownPieces = pos.GetOccupiedWhite()
	} else {
		knights = pos.BlackKnights
		ownPieces = pos.GetOccupiedBlack()
	}

	for knights != 0 {
		from := popLSB(&knights)
		toMask := g.knightMoves[from] & ^ownPieces
		for toMask != 0 {
			to := popLSB(&toMask)
			moves = append(moves, Move{Piece: Knight, From: Square(from), To: Square(to)})
		}
	}
	return moves
}

func (g *Generator) generateSlidingMoves(pos *Position, color Color, pieceType PieceType, directions []int) []Move {
	var pieces, ownPieces Bitboard
	if color == White {
		switch pieceType {
		case Bishop:
			pieces = pos.WhiteBishops
		case Rook:
			pieces = pos.WhiteRooks
		}
		ownPieces = pos.GetOccupiedWhite()
	} else {
		switch pieceType {
		case Bishop:
			pieces = pos.BlackBishops
		case Rook:
			pieces = pos.BlackRooks
		}
		ownPieces = pos.GetOccupiedBlack()
	}

	occupied := pos.GetOccupied()

	moves := []Move{}
	for pieces != 0 {
		from := popLSB(&pieces)
		for _, dir := range directions {
			to := from
			for {
				prevFile := to % 8
				to += dir
				toMask := squareMask(Square(to))
				newFile := to % 8

				// edges of the board
				fileDiff := newFile - prevFile
				if to < 0 || to > 63 || abs(fileDiff) > 1 {
					break
				}

				// blocked by own piece
				if toMask&ownPieces != 0 {
					break
				}
				moves = append(moves, Move{Piece: pieceType, From: Square(from), To: Square(to)})
				// blocked by opponent piece, added capture
				if toMask&occupied != 0 {
					break
				}
			}
		}
	}

	return moves
}

func (g *Generator) generateBishopMoves(pos *Position, color Color) []Move {
	directions := []int{7, 9, -7, -9} // NW, NE, SE, SW
	return g.generateSlidingMoves(pos, color, Bishop, directions)
}

func (g *Generator) generateRookMoves(pos *Position, color Color) []Move {
	directions := []int{1, -1, 8, -8}
	return g.generateSlidingMoves(pos, color, Rook, directions)
}

func (g *Generator) generateQueenMoves(pos *Position, color Color) []Move {
	// TODO
	return []Move{}
}

func (g *Generator) generateKingMoves(pos *Position, color Color) []Move {
	// TODO
	return []Move{}
}
