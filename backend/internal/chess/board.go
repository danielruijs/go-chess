package chess

import "fmt"

type Board [BoardSize][BoardSize]Piece // Files(columns) a-h, ranks(rows) 1-8

func (b Board) GetPiece(sq Square) (Piece, error) {
	if !sq.IsValid() {
		return Piece{}, fmt.Errorf("invalid square: %s", sq)
	}
	file := int(sq[0] - 'a')
	rank := int(sq[1] - '1')
	return b[file][rank], nil
}
