package chess

import "fmt"

type PGN string

// appends a move to the PGN string, assuming the move has not been made yet
func (p PGN) AppendMove(move Move, position *Position, generator *Generator) PGN {
	san := move.ToSAN(position, generator)
	if position.ActiveColor == White {
		return p + PGN(fmt.Sprintf("%d.%s ", position.Fullmove, san))
	} else if position.ActiveColor == Black && len(p) == 0 {
		// starting from black move
		return p + PGN(fmt.Sprintf("%d... %s ", position.Fullmove, san))
	} else {
		return p + PGN(fmt.Sprintf("%s ", san))
	}
}

func (p PGN) AppendResult(result *Result) PGN {
	if len(p) > 0 && p[len(p)-1] != ' ' {
		p += " "
	}

	switch result.Outcome {
	case WhiteWin:
		return p + "1-0"
	case BlackWin:
		return p + "0-1"
	case Draw:
		return p + "1/2-1/2"
	default:
		return p + "*"
	}
}
