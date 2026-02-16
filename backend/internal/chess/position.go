package chess

import (
	"fmt"
)

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

func (p *Position) ValidateMove(move Move, color Color) error {
	pieceToMove := p.GetPiece(move.From)
	if pieceToMove == (Piece{}) {
		return fmt.Errorf("no piece at source square")
	}
	if color != p.ActiveColor {
		return fmt.Errorf("not %s's turn to move", color)
	}
	if pieceToMove.Color != color {
		return fmt.Errorf("piece at source square does not belong to player")
	}

	moves := GeneratePieceMoves(p, color, pieceToMove.Type)
	for _, m := range moves {
		if m == move {
			return nil
		}
	}
	return fmt.Errorf("invalid move")
}
