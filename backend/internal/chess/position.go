package chess

import "log"

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

func NewInitialPosition() *Position {
	pos, err := StartingPositionFEN.ToPosition()
	if err != nil {
		log.Fatalf("failed to create initial position: %v", err)
	}
	return &pos
}

func (p *Position) GetOccupied() Bitboard {
	return p.GetOccupiedWhite() | p.GetOccupiedBlack()
}

func (p *Position) GetOccupiedWhite() Bitboard {
	return p.WhitePawns | p.WhiteKnights | p.WhiteBishops | p.WhiteRooks | p.WhiteQueens | p.WhiteKing
}

func (p *Position) GetOccupiedBlack() Bitboard {
	return p.BlackPawns | p.BlackKnights | p.BlackBishops | p.BlackRooks | p.BlackQueens | p.BlackKing
}

func (p *Position) GetPiece(sq Square) Piece {
	mask := squareMask(sq)

	switch {
	case p.WhitePawns&mask != 0:
		return Piece{Type: Pawn, Color: White}
	case p.WhiteKnights&mask != 0:
		return Piece{Type: Knight, Color: White}
	case p.WhiteBishops&mask != 0:
		return Piece{Type: Bishop, Color: White}
	case p.WhiteRooks&mask != 0:
		return Piece{Type: Rook, Color: White}
	case p.WhiteQueens&mask != 0:
		return Piece{Type: Queen, Color: White}
	case p.WhiteKing&mask != 0:
		return Piece{Type: King, Color: White}

	case p.BlackPawns&mask != 0:
		return Piece{Type: Pawn, Color: Black}
	case p.BlackKnights&mask != 0:
		return Piece{Type: Knight, Color: Black}
	case p.BlackBishops&mask != 0:
		return Piece{Type: Bishop, Color: Black}
	case p.BlackRooks&mask != 0:
		return Piece{Type: Rook, Color: Black}
	case p.BlackQueens&mask != 0:
		return Piece{Type: Queen, Color: Black}
	case p.BlackKing&mask != 0:
		return Piece{Type: King, Color: Black}
	}

	return Piece{}
}

func (p *Position) GetBoard() Board {
	var board Board
	for file := range BoardSize {
		for rank := range BoardSize {
			square := coordsToSquare(file, rank)
			piece := p.GetPiece(square)
			board[file][rank] = piece
		}
	}
	return board
}
