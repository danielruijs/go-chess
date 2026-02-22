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
	kingMoves   [64]Bitboard

	bishopMoveDirections []int
	rookMoveDirections   []int

	pieceToGenerator map[PieceType]func(pos *Position, color Color) []Move
}

func NewGenerator() *Generator {
	g := &Generator{}

	g.knightMoves = precomputeKnightMoves()
	g.kingMoves = precomputeKingMoves()

	g.bishopMoveDirections = []int{7, 9, -7, -9} // NW, NE, SE, SW
	g.rookMoveDirections = []int{1, -1, 8, -8}   // E, W, N, S

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
	psuedoLegalMoves := []Move{}
	for _, generator := range g.pieceToGenerator {
		psuedoLegalMoves = append(psuedoLegalMoves, generator(pos, color)...)
	}
	return g.filterLegalMoves(pos, color, psuedoLegalMoves)
}

// Generates all legal moves in a position for a given color and piece type
func (g *Generator) GeneratePieceMoves(pos *Position, color Color, pieceType PieceType) []Move {
	generator := g.pieceToGenerator[pieceType]
	psuedoLegalMoves := generator(pos, color)
	return g.filterLegalMoves(pos, color, psuedoLegalMoves)
}

// Filters pseudo legal moves to only legal moves
func (g *Generator) filterLegalMoves(pos *Position, color Color, moves []Move) []Move {
	legalMoves := []Move{}
	for _, move := range moves {
		posCopy := pos.GetCopy()
		posCopy.MakeMove(move)
		if !posCopy.IsInCheck(g, color) {
			// TODO: Special case for castling: Cant castle out of, through, or into check.
			legalMoves = append(legalMoves, move)
		}
	}
	return legalMoves
}

func (g *Generator) generatePawnMoves(pos *Position, color Color) []Move {
	moves := []Move{}

	var pawns, startingRankMask, promotionRankMask, opponentPieces Bitboard
	var forwardOffset int
	if color == White {
		pawns = pos.WhitePawns
		startingRankMask = rank2Mask
		promotionRankMask = rank8Mask
		opponentPieces = pos.GetOccupiedBlack()
		forwardOffset = 8
	} else {
		pawns = pos.BlackPawns
		startingRankMask = rank7Mask
		promotionRankMask = rank1Mask
		opponentPieces = pos.GetOccupiedWhite()
		forwardOffset = -8
	}

	occupied := pos.GetOccupied()

	// regular forward move by 1 if not blocked
	singlePushes := shift(pawns, forwardOffset) & ^occupied & ^promotionRankMask
	for singlePushes != 0 {
		to := popLSB(&singlePushes)
		from := to - forwardOffset
		moves = append(moves, Move{From: Square(from), To: Square(to)})
	}

	// move forward by 2 if on starting rank and not blocked
	oneStep := shift((pawns&startingRankMask), forwardOffset) & ^occupied
	doublePushes := shift(oneStep, forwardOffset) & ^occupied
	for doublePushes != 0 {
		to := popLSB(&doublePushes)
		from := to - 2*forwardOffset
		moves = append(moves, Move{From: Square(from), To: Square(to)})
	}

	// capture diagonally by 1 if occupied by opponent piece
	capturesLeft := shift(pawns, forwardOffset-1) & opponentPieces & ^fileAMask
	capturesRight := shift(pawns, forwardOffset+1) & opponentPieces & ^fileHMask

	// en passant
	if pos.EnPassant != nil {
		enPassantMask := squareMask(*pos.EnPassant)
		epCapturesLeft := shift(pawns, forwardOffset-1) & enPassantMask & ^fileAMask
		epCapturesRight := shift(pawns, forwardOffset+1) & enPassantMask & ^fileHMask

		capturesLeft |= epCapturesLeft
		capturesRight |= epCapturesRight
	}
	for capturesLeft != 0 {
		to := popLSB(&capturesLeft)
		from := to - forwardOffset + 1
		moves = append(moves, Move{From: Square(from), To: Square(to)})
	}
	for capturesRight != 0 {
		to := popLSB(&capturesRight)
		from := to - forwardOffset - 1
		moves = append(moves, Move{From: Square(from), To: Square(to)})
	}

	// promotion
	promotions := shift(pawns, forwardOffset) & ^occupied & promotionRankMask
	for promotions != 0 {
		to := popLSB(&promotions)
		from := to - forwardOffset
		for _, promotionPiece := range []PieceType{Queen, Rook, Bishop, Knight} {
			moves = append(moves, Move{From: Square(from), To: Square(to), Promotion: &promotionPiece})
		}
	}

	return moves
}

func precomputeNonSlidingMoves(offsets []struct{ file, rank int }) [64]Bitboard {
	var moves [64]Bitboard

	for sq := range 64 {
		var moveMask Bitboard
		file := sq % 8
		rank := sq / 8
		for _, offset := range offsets {
			newFile := file + offset.file
			newRank := rank + offset.rank
			if newRank >= 0 && newRank < 8 && newFile >= 0 && newFile < 8 {
				moveMask |= coordMask(newFile, newRank)
			}
		}
		moves[sq] = moveMask
	}

	return moves
}

func precomputeKnightMoves() [64]Bitboard {
	offsets := []struct{ file, rank int }{
		{2, 1}, {2, -1}, {-2, 1}, {-2, -1},
		{1, 2}, {1, -2}, {-1, 2}, {-1, -2},
	}
	return precomputeNonSlidingMoves(offsets)
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
			moves = append(moves, Move{From: Square(from), To: Square(to)})
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
		case Queen:
			pieces = pos.WhiteQueens
		}
		ownPieces = pos.GetOccupiedWhite()
	} else {
		switch pieceType {
		case Bishop:
			pieces = pos.BlackBishops
		case Rook:
			pieces = pos.BlackRooks
		case Queen:
			pieces = pos.BlackQueens
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
				if to < 0 || to > 63 || abs(newFile-prevFile) > 1 {
					break
				}

				// blocked by own piece
				if toMask&ownPieces != 0 {
					break
				}
				moves = append(moves, Move{From: Square(from), To: Square(to)})
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
	directions := g.bishopMoveDirections
	return g.generateSlidingMoves(pos, color, Bishop, directions)
}

func (g *Generator) generateRookMoves(pos *Position, color Color) []Move {
	directions := g.rookMoveDirections
	return g.generateSlidingMoves(pos, color, Rook, directions)
}

func (g *Generator) generateQueenMoves(pos *Position, color Color) []Move {
	directions := append(g.bishopMoveDirections, g.rookMoveDirections...)
	return g.generateSlidingMoves(pos, color, Queen, directions)
}

func precomputeKingMoves() [64]Bitboard {
	offsets := []struct{ file, rank int }{
		{1, 0}, {-1, 0}, {0, 1}, {0, -1},
		{1, 1}, {1, -1}, {-1, 1}, {-1, -1},
	}
	return precomputeNonSlidingMoves(offsets)
}

func (g *Generator) generateKingMoves(pos *Position, color Color) []Move {
	moves := []Move{}

	var king, ownPieces Bitboard
	var kingSideRight, queenSideRight bool
	var kingSideMask, queenSideMask Bitboard
	var kingSideTo, queenSideTo Square
	if color == White {
		king = pos.WhiteKing
		ownPieces = pos.GetOccupiedWhite()

		kingSideRight = pos.CastlingRights.WhiteOO
		queenSideRight = pos.CastlingRights.WhiteOOO
		kingSideMask = 0x60 // f1,g1
		queenSideMask = 0xE // b1,c1,d1
		kingSideTo = Square(6)
		queenSideTo = Square(2)
	} else {
		king = pos.BlackKing
		ownPieces = pos.GetOccupiedBlack()

		kingSideRight = pos.CastlingRights.BlackOO
		queenSideRight = pos.CastlingRights.BlackOOO
		kingSideMask = 0x6000000000000000  // f8,g8
		queenSideMask = 0x0E00000000000000 // b8,c8,d8
		kingSideTo = Square(62)
		queenSideTo = Square(58)
	}

	from := popLSB(&king) // should only be one king
	toMask := g.kingMoves[from] & ^ownPieces
	for toMask != 0 {
		to := popLSB(&toMask)
		moves = append(moves, Move{From: Square(from), To: Square(to)})
	}

	// Castling, only checks rights and that the squares between the king and rook are empty.
	// Legality (not castling out of, through, or into check) is checked in filterLegalMoves
	occupied := pos.GetOccupied()
	if kingSideRight && occupied&kingSideMask == 0 {
		moves = append(moves, Move{From: Square(from), To: kingSideTo})
	}
	if queenSideRight && occupied&queenSideMask == 0 {
		moves = append(moves, Move{From: Square(from), To: queenSideTo})
	}

	return moves
}

// Checks if a square is attacked in a position for a given color
func (g *Generator) IsSquareAttacked(sq Square, pos *Position, color Color) bool {
	occupied := pos.GetOccupied()
	sqMask := squareMask(sq)

	var pawns, knights, bishops, rooks, queens, king Bitboard
	var forwardOffset int
	if color == White {
		pawns = pos.BlackPawns
		knights = pos.BlackKnights
		bishops = pos.BlackBishops
		rooks = pos.BlackRooks
		queens = pos.BlackQueens
		king = pos.BlackKing
		forwardOffset = 8
	} else {
		pawns = pos.WhitePawns
		knights = pos.WhiteKnights
		bishops = pos.WhiteBishops
		rooks = pos.WhiteRooks
		queens = pos.WhiteQueens
		king = pos.WhiteKing
		forwardOffset = -8
	}

	// to check if sq is attacked by a pawn, do the inverses
	// project a pawn's capture move FROM the target square
	// if that projection lands on an opponent's pawn, the square is under attack
	pawnAttacksLeft := shift(sqMask, forwardOffset-1) & ^fileAMask
	pawnAttackRight := shift(sqMask, forwardOffset+1) & ^fileHMask
	if (pawnAttacksLeft|pawnAttackRight)&pawns != 0 {
		return true
	}

	// same inverse logic, project knight moves from sq and see if any hit opponent knights
	if g.knightMoves[sq]&knights != 0 {
		return true
	}

	if isAttackedBySlidingPiece(sq, g.bishopMoveDirections, bishops|queens, occupied) {
		return true
	}

	if isAttackedBySlidingPiece(sq, g.rookMoveDirections, rooks|queens, occupied) {
		return true
	}

	if g.kingMoves[sq]&king != 0 {
		return true
	}

	return false
}

func isAttackedBySlidingPiece(sq Square, directions []int, attackers, occupied Bitboard) bool {
	// go in each direction from sq, if we hit an attacker sq is attacked
	for _, dir := range directions {
		to := int(sq)
		for {
			prevFile := to % 8
			to += dir
			toMask := squareMask(Square(to))
			newFile := to % 8

			if to < 0 || to > 63 || abs(newFile-prevFile) > 1 {
				break
			}

			if toMask&attackers != 0 {
				return true
			}
			if toMask&occupied != 0 {
				break
			}
		}
	}

	return false
}
