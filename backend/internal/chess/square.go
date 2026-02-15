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
	file := s[0]
	rank := s[1]
	return file >= 'a' && file <= 'h' && rank >= '1' && rank <= '8'
}
