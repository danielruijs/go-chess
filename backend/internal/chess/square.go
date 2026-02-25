package chess

import "fmt"

type Square uint8 // 0 -> a1, 1 -> b1, ..., 63 -> h8

func StrToSquare(str string) (Square, error) {
	if !isValidStr(str) {
		return 0, fmt.Errorf("invalid square string: %s", str)
	}
	file := int(str[0] - 'a')
	rank := int(str[1] - '1')
	return Square(file + rank*8), nil
}

func isValidStr(str string) bool {
	if len(str) != 2 {
		return false
	}
	file := str[0]
	rank := str[1]
	return file >= 'a' && file <= 'h' && rank >= '1' && rank <= '8'
}

func SquareToStr(sq Square) string {
	file := sq % 8
	rank := sq / 8
	return fmt.Sprintf("%c%d", 'a'+file, rank+1)
}

func coordsToSquare(file, rank int) Square {
	return Square(file + rank*8)
}

// Returns bitboard mask for given square
func squareMask(sq Square) Bitboard {
	return Bitboard(1) << sq
}

// Returns bitboard mask for given coordinates
func coordMask(file, rank int) Bitboard {
	return squareMask(coordsToSquare(file, rank))
}

func (sq Square) Color() Color {
	file := sq % 8
	rank := sq / 8
	if (rank+file)%2 == 0 {
		return Black
	}
	return White
}
