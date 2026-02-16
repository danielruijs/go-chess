package chess

type Generator interface {
	generateMoves(position *Position, color Color) []Move
}

var pieceToGenerator = map[PieceType]Generator{
	Pawn: PawnMoveGenerator{},
}

// generates all legal moves in a position for given color
func GenerateMoves(position *Position, color Color) []Move {
	// TODO
	return []Move{}
}

// generates all legal moves in a position for given color and piece type
func GeneratePieceMoves(position *Position, color Color, pieceType PieceType) []Move {
	generator := pieceToGenerator[pieceType]
	return generator.generateMoves(position, color)
}

type PawnMoveGenerator struct{}

func (g PawnMoveGenerator) generateMoves(position *Position, color Color) []Move {
	// regular forward move by 1 if not blocked
	// move forward by 2 if on starting rank and not blocked
	// capture diagonally by 1 if occupied by opponent piece
	// en passant
	// promotion
	return nil
}
