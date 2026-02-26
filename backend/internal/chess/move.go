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

func samePromotion(a, b *PieceType) bool {
	// either both are nil or both are not nil and same value
	return a == b || (a != nil && b != nil && *a == *b)
}

func ContainsMove(moves []Move, m Move) bool {
	for _, move := range moves {
		if move.From == m.From &&
			move.To == m.To &&
			samePromotion(move.Promotion, m.Promotion) {
			return true
		}
	}
	return false
}

func (m Move) ToSAN(position *Position) string {
	san := ""
	piece := position.GetPiece(m.From)
	san += PieceTypeToSAN[piece.Type]

	// capture
	if position.GetPiece(m.To) != nil || m.IsEnPassant(position) {
		if piece.Type == Pawn {
			// for pawn captures, include the file of the from square
			san += SquareToStr(m.From)[0:1]
		}
		san += "x"
	}

	san += SquareToStr(m.To)
	return san
}
