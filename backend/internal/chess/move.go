package chess

type Move struct {
	From      Square
	To        Square
	Promotion *PieceType
}

func (m Move) IsEnPassant(p *Position) bool {
	piece := p.GetPiece(m.From)
	isPawn := piece.Type == Pawn
	isCapture := (m.From-m.To)%8 != 0
	return isPawn && isCapture && p.EnPassant != nil && *p.EnPassant == m.To
}
