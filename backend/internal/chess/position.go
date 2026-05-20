package chess

type Bitboard uint64 // 64-bit integer representing the presence of pieces on the board, bit 0 (LSB) -> a1, bit 1 -> b1, ..., bit 63 (MSB) -> h8

type CastlingRights struct {
	WhiteOO  bool
	WhiteOOO bool
	BlackOO  bool
	BlackOOO bool
}

type Position struct {
	WhitePawns   Bitboard
	WhiteKnights Bitboard
	WhiteBishops Bitboard
	WhiteRooks   Bitboard
	WhiteQueens  Bitboard
	WhiteKing    Bitboard

	BlackPawns   Bitboard
	BlackKnights Bitboard
	BlackBishops Bitboard
	BlackRooks   Bitboard
	BlackQueens  Bitboard
	BlackKing    Bitboard

	ActiveColor    Color          // Color to move
	CastlingRights CastlingRights // Castling rights
	EnPassant      *Square        // En passant target square, square over which pawn just moved when moving two squares
	Halfmove       uint           // Halfmove clock, number of halfmoves since last capture or pawn move, for fifty-move rule
	Fullmove       uint           // Fullmove number
}

type PositionKey struct {
	WhitePawns, WhiteKnights, WhiteBishops, WhiteRooks, WhiteQueens, WhiteKing Bitboard
	BlackPawns, BlackKnights, BlackBishops, BlackRooks, BlackQueens, BlackKing Bitboard
	ActiveColor                                                                Color
	Castling                                                                   CastlingRights
	EnPassant                                                                  *Square
}

func NewInitialPosition() *Position {
	pos, err := StartingPositionFEN.ToPosition()
	if err != nil {
		panic("failed to parse starting position FEN: " + err.Error())
	}
	return &pos
}

func (p *Position) GetOccupied() Bitboard {
	return p.getOccupiedWhite() | p.getOccupiedBlack()
}

func (p *Position) getOccupiedWhite() Bitboard {
	return p.WhitePawns | p.WhiteKnights | p.WhiteBishops | p.WhiteRooks | p.WhiteQueens | p.WhiteKing
}

func (p *Position) getOccupiedBlack() Bitboard {
	return p.BlackPawns | p.BlackKnights | p.BlackBishops | p.BlackRooks | p.BlackQueens | p.BlackKing
}

func (p *Position) GetOccupiedByColor(color Color) Bitboard {
	if color == White {
		return p.getOccupiedWhite()
	} else {
		return p.getOccupiedBlack()
	}
}

func (p *Position) GetPieceBitboard(color Color, pieceType PieceType) *Bitboard {
	switch color {
	case White:
		switch pieceType {
		case Pawn:
			return &p.WhitePawns
		case Knight:
			return &p.WhiteKnights
		case Bishop:
			return &p.WhiteBishops
		case Rook:
			return &p.WhiteRooks
		case Queen:
			return &p.WhiteQueens
		case King:
			return &p.WhiteKing
		}
	case Black:
		switch pieceType {
		case Pawn:
			return &p.BlackPawns
		case Knight:
			return &p.BlackKnights
		case Bishop:
			return &p.BlackBishops
		case Rook:
			return &p.BlackRooks
		case Queen:
			return &p.BlackQueens
		case King:
			return &p.BlackKing
		}
	}
	panic("invalid color or piece type in GetPieceBitboard")
}

func (p *Position) GetPiece(sq Square) *Piece {
	mask := squareMask(sq)

	for _, color := range []Color{White, Black} {
		for _, pieceType := range []PieceType{Pawn, Knight, Bishop, Rook, Queen, King} {
			pieceBitboard := p.GetPieceBitboard(color, pieceType)
			if *pieceBitboard&mask != 0 {
				return &Piece{Type: pieceType, Color: color}
			}
		}
	}

	return nil
}

func (p *Position) GetBoard() Board {
	var board Board
	for file := range 8 {
		for rank := range 8 {
			square := coordsToSquare(file, rank)
			piece := p.GetPiece(square)
			board[file][rank] = piece
		}
	}
	return board
}

func (p *Position) GetCopy() Position {
	copy := *p

	if p.EnPassant != nil {
		ep := *p.EnPassant
		copy.EnPassant = &ep
	}

	return copy
}

func (p *Position) setPiece(sq Square, piece *Piece) {
	if piece == nil {
		return
	}
	mask := squareMask(sq)
	pieceBitboard := p.GetPieceBitboard(piece.Color, piece.Type)
	*pieceBitboard |= mask
}

func (p *Position) removePiece(sq Square) {
	mask := ^squareMask(sq)
	piece := p.GetPiece(sq)
	pieceBitboard := p.GetPieceBitboard(piece.Color, piece.Type)
	*pieceBitboard &= mask
}

// Makes a move, assuming it is legal
func (p *Position) MakeMove(move Move) {
	piece := p.GetPiece(move.From)
	occupied := p.GetOccupied()
	toMask := squareMask(move.To)
	isCapture := false
	var capturedPiece *Piece
	var capturedSquare Square

	if move.IsEnPassant(p) {
		// en passant
		isCapture = true
		p.removePiece(move.From)
		p.setPiece(move.To, piece)
		if piece.Color == White {
			capturedSquare = Square(move.To - 8)
		} else {
			capturedSquare = Square(move.To + 8)
		}
		capturedPiece = p.GetPiece(capturedSquare)
		p.removePiece(capturedSquare)
	} else if move.IsCastling(p) {
		// castling
		var rookFrom, rookTo Square
		if move.To > move.From {
			// king-side
			rookFrom = move.From + 3
			rookTo = move.From + 1
		} else {
			// queen-side
			rookFrom = move.From - 4
			rookTo = move.From - 1
		}
		p.removePiece(move.From)
		p.setPiece(move.To, piece)
		p.removePiece(rookFrom)
		p.setPiece(rookTo, &Piece{Type: Rook, Color: piece.Color})
	} else if move.Promotion != nil {
		// promotion, can be capture
		if occupied&toMask != 0 {
			isCapture = true
			capturedPiece = p.GetPiece(move.To)
			capturedSquare = move.To
			p.removePiece(move.To)
		}
		p.removePiece(move.From)
		p.setPiece(move.To, &Piece{Type: *move.Promotion, Color: piece.Color})
	} else if occupied&toMask != 0 {
		// captures
		isCapture = true
		capturedPiece = p.GetPiece(move.To)
		capturedSquare = move.To
		p.removePiece(move.To)
		p.removePiece(move.From)
		p.setPiece(move.To, piece)
	} else {
		// regular moves
		p.removePiece(move.From)
		p.setPiece(move.To, piece)
	}

	p.updateCastlingRights(move, piece, capturedPiece, capturedSquare)

	// update en passant square for double pawn moves
	p.EnPassant = nil
	if piece.Type == Pawn && abs(int(move.To)-int(move.From)) == 16 {
		if piece.Color == White {
			p.EnPassant = new(Square(move.To - 8))
		} else {
			p.EnPassant = new(Square(move.To + 8))
		}
	}

	// update halfmove clock
	if piece.Type == Pawn || isCapture {
		p.Halfmove = 0
	} else {
		p.Halfmove++
	}

	// update fullmove number
	if piece.Color == Black {
		p.Fullmove++
	}

	// switch active color
	p.ActiveColor = GetOppositeColor(p.ActiveColor)
}

func (p *Position) updateCastlingRights(move Move, piece *Piece, capturedPiece *Piece, capturedSquare Square) {
	switch piece.Type {
	case King:
		if piece.Color == White {
			p.CastlingRights.WhiteOO = false
			p.CastlingRights.WhiteOOO = false
		} else {
			p.CastlingRights.BlackOO = false
			p.CastlingRights.BlackOOO = false
		}
	case Rook:
		if piece.Color == White {
			if move.From == Square(7) { // h1
				p.CastlingRights.WhiteOO = false
			} else if move.From == Square(0) { // a1
				p.CastlingRights.WhiteOOO = false
			}
		} else {
			if move.From == Square(63) { // h8
				p.CastlingRights.BlackOO = false
			} else if move.From == Square(56) { // a8
				p.CastlingRights.BlackOOO = false
			}
		}
	}

	// if a rook was captured, clear the corresponding castling right
	if capturedPiece != nil && capturedPiece.Type == Rook {
		// determine which square the rook was on when captured
		switch capturedSquare {
		case Square(7): // h1
			p.CastlingRights.WhiteOO = false
		case Square(0): // a1
			p.CastlingRights.WhiteOOO = false
		case Square(63): // h8
			p.CastlingRights.BlackOO = false
		case Square(56): // a8
			p.CastlingRights.BlackOOO = false
		}
	}
}

func (p *Position) IsInCheck(g *Generator, color Color) bool {
	var king Bitboard
	if color == White {
		king = p.WhiteKing
	} else {
		king = p.BlackKing
	}

	kingSq := Square(popLSB(&king))
	return g.IsSquareAttacked(kingSq, p, color)
}

func (p *Position) Key() PositionKey {
	return PositionKey{
		WhitePawns:   p.WhitePawns,
		WhiteKnights: p.WhiteKnights,
		WhiteBishops: p.WhiteBishops,
		WhiteRooks:   p.WhiteRooks,
		WhiteQueens:  p.WhiteQueens,
		WhiteKing:    p.WhiteKing,
		BlackPawns:   p.BlackPawns,
		BlackKnights: p.BlackKnights,
		BlackBishops: p.BlackBishops,
		BlackRooks:   p.BlackRooks,
		BlackQueens:  p.BlackQueens,
		BlackKing:    p.BlackKing,
		ActiveColor:  p.ActiveColor,
		Castling:     p.CastlingRights,
		EnPassant:    p.EnPassant,
	}
}
