package chess

import (
	"fmt"
	"strconv"
	"strings"
)

type Fen string

var FenCharToPiece = map[rune]Piece{
	'P': {Type: Pawn, Color: White},
	'N': {Type: Knight, Color: White},
	'B': {Type: Bishop, Color: White},
	'R': {Type: Rook, Color: White},
	'Q': {Type: Queen, Color: White},
	'K': {Type: King, Color: White},
	'p': {Type: Pawn, Color: Black},
	'n': {Type: Knight, Color: Black},
	'b': {Type: Bishop, Color: Black},
	'r': {Type: Rook, Color: Black},
	'q': {Type: Queen, Color: Black},
	'k': {Type: King, Color: Black},
}

var FenCharToCastling = map[rune]func(*CastlingRights){
	'K': func(c *CastlingRights) { c.WhiteOO = true },
	'Q': func(c *CastlingRights) { c.WhiteOOO = true },
	'k': func(c *CastlingRights) { c.BlackOO = true },
	'q': func(c *CastlingRights) { c.BlackOOO = true },
}

func (f Fen) ToPosition() (Position, error) {
	fields := strings.Split(string(f), " ")
	if len(fields) != 6 {
		return Position{}, fmt.Errorf("fen string must have 6 fields")
	}

	pos, err := parseBoard(fields[0])
	if err != nil {
		return Position{}, err
	}

	activeColor, err := parseColor(fields[1])
	if err != nil {
		return Position{}, err
	}
	pos.ActiveColor = activeColor

	castlingRights, err := parseCastlingRights(fields[2])
	if err != nil {
		return Position{}, err
	}
	pos.CastlingRights = castlingRights

	enPassant, err := parseEnPassant(fields[3])
	if err != nil {
		return Position{}, err
	}
	pos.EnPassant = enPassant

	halfmove, err := parseHalfmove(fields[4])
	if err != nil {
		return Position{}, err
	}
	pos.Halfmove = uint(halfmove)

	fullmove, err := parseFullmove(fields[5])
	if err != nil {
		return Position{}, err
	}
	pos.Fullmove = uint(fullmove)

	return pos, nil
}

func parseBoard(s string) (Position, error) {
	var pos Position

	ranks := strings.Split(s, "/")
	if len(ranks) != 8 {
		return Position{}, fmt.Errorf("fen piece placement must have 8 ranks")
	}
	// Ranks 8 to 1
	for fenRank, rankStr := range ranks {
		file := 0
		// Files a to h
		for _, char := range rankStr {
			// Empty squares
			if char >= '1' && char <= '8' {
				file += int(char - '0')
				continue
			}

			piece, ok := FenCharToPiece[char]
			if !ok {
				return Position{}, fmt.Errorf("invalid fen character: %c", char)
			}

			rank := 7 - fenRank
			mask := coordMask(file, rank)

			switch piece.Color {
			case White:
				switch piece.Type {
				case Pawn:
					pos.WhitePawns |= mask
				case Knight:
					pos.WhiteKnights |= mask
				case Bishop:
					pos.WhiteBishops |= mask
				case Rook:
					pos.WhiteRooks |= mask
				case Queen:
					pos.WhiteQueens |= mask
				case King:
					pos.WhiteKing |= mask
				}
			case Black:
				switch piece.Type {
				case Pawn:
					pos.BlackPawns |= mask
				case Knight:
					pos.BlackKnights |= mask
				case Bishop:
					pos.BlackBishops |= mask
				case Rook:
					pos.BlackRooks |= mask
				case Queen:
					pos.BlackQueens |= mask
				case King:
					pos.BlackKing |= mask
				}
			}

			file++
		}

		if file != 8 {
			return Position{}, fmt.Errorf("fen piece placement must have 8 files")
		}
	}
	return pos, nil
}

func parseColor(s string) (Color, error) {
	switch s {
	case "w":
		return White, nil
	case "b":
		return Black, nil
	default:
		return Color(""), fmt.Errorf("invalid active color in fen: %s", s)
	}
}

func parseCastlingRights(s string) (CastlingRights, error) {
	var castlingRights CastlingRights
	for _, char := range s {
		if char == '-' {
			break
		}
		fn, ok := FenCharToCastling[char]
		if !ok {
			return CastlingRights{}, fmt.Errorf("invalid castling right in fen: %c", char)
		}
		fn(&castlingRights)
	}
	return castlingRights, nil
}

func parseEnPassant(s string) (*Square, error) {
	if s == "-" {
		return nil, nil
	}
	enPassant, err := StrToSquare(s)
	if err != nil {
		return nil, fmt.Errorf("invalid en passant target square in fen: %s", s)
	}
	return &enPassant, nil
}

func parseHalfmove(s string) (int, error) {
	halfmove, err := strconv.Atoi(s)
	if err != nil || halfmove < 0 {
		return 0, fmt.Errorf("invalid halfmove clock in fen: %s", s)
	}
	return halfmove, nil
}

func parseFullmove(s string) (int, error) {
	fullmove, err := strconv.Atoi(s)
	if err != nil || fullmove < 1 {
		return 0, fmt.Errorf("invalid fullmove number in fen: %s", s)
	}
	return fullmove, nil
}
