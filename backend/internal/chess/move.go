package chess

type Move struct {
	Piece     PieceType
	From      Square
	To        Square
	Promotion *PieceType
}
