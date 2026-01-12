package chess

import "fmt"

type Square string

func StrToSquare(s string) (Square, error) {
	sq := Square(s)
	if !sq.IsValid() {
		return "", fmt.Errorf("invalid square: %s", s)
	}
	return sq, nil
}

func (s Square) IsValid() bool {
	if len(s) != 2 {
		return false
	}
	row := s[0]
	col := s[1]
	return row >= 'a' && row <= 'h' && col >= '1' && col <= '8'
}

func (b Board) GetSquare(sq Square) (Piece, error) {
	if !sq.IsValid() {
		return Piece{}, fmt.Errorf("invalid square: %s", sq)
	}
	row := int(sq[1] - '1')
	col := int(sq[0] - 'a')
	return b[row][col], nil
}
