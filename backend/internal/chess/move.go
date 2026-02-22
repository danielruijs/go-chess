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

func (m Move) IsCastling(p *Position) bool {
	piece := p.GetPiece(m.From)
	isKing := piece.Type == King
	isSameRank := m.From/8 == m.To/8
	isLongMove := abs(int(m.To)-int(m.From)) > 1
	return isKing && isLongMove && isSameRank
}
