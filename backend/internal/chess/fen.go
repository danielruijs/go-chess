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

	board, err := parseBoard(fields[0])
	if err != nil {
		return Position{}, err
	}

	activeColor, err := parseColor(fields[1])
	if err != nil {
		return Position{}, err
	}

	castlingRights, err := parseCastlingRights(fields[2])
	if err != nil {
		return Position{}, err
	}

	enPassant, err := parseEnPassant(fields[3])
	if err != nil {
		return Position{}, err
	}

	halfmove, err := parseHalfmove(fields[4])
	if err != nil {
		return Position{}, err
	}

	fullmove, err := parseFullmove(fields[5])
	if err != nil {
		return Position{}, err
	}

	return Position{
		Board:          board,
		ActiveColor:    activeColor,
		CastlingRights: castlingRights,
		EnPassant:      enPassant,
		Halfmove:       uint(halfmove),
		Fullmove:       uint(fullmove),
	}, nil
}

func parseBoard(s string) (Board, error) {
	var board Board
	rows := strings.Split(s, "/")
	if len(rows) != BoardSize {
		return Board{}, fmt.Errorf("fen piece placement must have %d rows", BoardSize)
	}
	// Rows 8 to 1
	for i, row := range rows {
		j := 0
		for _, char := range row {
			if j >= BoardSize {
				return Board{}, fmt.Errorf("fen piece placement has too many columns")
			}
			// Empty squares
			if char >= '1' && char <= '8' {
				emptySquares := int(char - '0')
				for k := range emptySquares {
					board[i][j+k] = Piece{}
				}
				j += emptySquares
				continue
			}
			// Pieces
			piece, ok := FenCharToPiece[char]
			if !ok {
				return Board{}, fmt.Errorf("invalid fen character: %c", char)
			}
			board[i][j] = piece
			j++
		}
		if j != BoardSize {
			return Board{}, fmt.Errorf("fen piece placement must have %d columns", BoardSize)
		}
	}
	return board, nil
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

func parseEnPassant(s string) (Square, error) {
	if s == "-" {
		return Square(""), nil
	}
	enPassant, err := StrToSquare(s)
	if err != nil {
		return Square(""), fmt.Errorf("invalid en passant target square in fen: %s", s)
	}
	return enPassant, nil
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
