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
	isSameRank := m.From.Rank() == m.To.Rank()
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

// returns the move in SAN format, assuming the move has not been made yet
func (m Move) ToSAN(position *Position, generator *Generator) string {
	san := ""
	if m.IsCastling(position) {
		// castling
		if m.To > m.From {
			san += "O-O"
		} else {
			san += "O-O-O"
		}
	} else {
		piece := position.GetPiece(m.From)
		san += PieceTypeToSAN[piece.Type]

		// disambiguation
		if piece.Type != Pawn {
			san += getDisambiguation(m, position, generator)
		}

		// capture
		if position.GetPiece(m.To) != nil || m.IsEnPassant(position) {
			if piece.Type == Pawn {
				// for pawn captures, include the file of the from square
				san += SquareToStr(m.From)[0:1]
			}
			san += "x"
		}

		// to square
		san += SquareToStr(m.To)

		// promotion
		if m.Promotion != nil {
			san += "=" + PieceTypeToSAN[*m.Promotion]
		}
	}

	// check or checkmate
	posCopy := position.GetCopy()
	posCopy.MakeMove(m)
	if posCopy.IsInCheck(generator, posCopy.ActiveColor) {
		if len(generator.GenerateMoves(&posCopy, posCopy.ActiveColor)) == 0 {
			san += "#"
		} else {
			san += "+"
		}
	}

	return san
}

func getDisambiguation(m Move, position *Position, generator *Generator) string {
	piece := position.GetPiece(m.From)
	pieceMoves := generator.GeneratePieceMoves(position, position.ActiveColor, piece.Type)

	var ambiguousMoves []Move
	for _, move := range pieceMoves {
		if move.From != m.From && move.To == m.To {
			ambiguousMoves = append(ambiguousMoves, move)
		}
	}
	if len(ambiguousMoves) == 0 {
		return ""
	}

	sameFile := false
	sameRank := false
	for _, move := range ambiguousMoves {
		if move.From.File() == m.From.File() {
			sameFile = true
		}
		if move.From.Rank() == m.From.Rank() {
			sameRank = true
		}
	}
	if !sameFile {
		return SquareToStr(m.From)[0:1] // file disambiguates
	}
	if !sameRank {
		return SquareToStr(m.From)[1:2] // rank disambiguates
	}
	return SquareToStr(m.From) // full square
}
