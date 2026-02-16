package chess

type Move struct {
	From      Square
	To        Square
	Promotion *PieceType
}
