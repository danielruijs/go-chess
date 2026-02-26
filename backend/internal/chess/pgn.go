package chess

import "fmt"

type PGN string

func (p PGN) AppendMove(move Move, position *Position) PGN {
	san := move.ToSAN(position)
	if position.ActiveColor == White {
		return p + PGN(fmt.Sprintf("%d.%s ", position.Fullmove, san))
	} else {
		return p + PGN(fmt.Sprintf("%s ", san))
	}
	//TODO
}
