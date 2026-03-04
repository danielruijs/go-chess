package chess

import "math/bits"

type Outcome string

const (
	WhiteWin Outcome = "white_win"
	BlackWin Outcome = "black_win"
	Draw     Outcome = "draw"
)

type Reason string

const (
	Checkmate            Reason = "checkmate"
	Stalemate            Reason = "stalemate"
	ThreefoldRepetition  Reason = "threefold_repetition"
	FiftyMoveRule        Reason = "fifty_move_rule"
	InsufficientMaterial Reason = "insufficient_material"
	Resignation          Reason = "resignation"
	AgreedDraw           Reason = "AGREED_DRAW"
	// Timeout              Reason = "TIMEOUT"
)

type Result struct {
	Outcome Outcome `json:"outcome"`
	Reason  Reason  `json:"reason"`
}

// Returns a Result if the game has ended and nil otherwise.
func getResult(generator *Generator, position *Position, positions map[PositionKey]int) *Result {
	legalMoves := generator.GenerateMoves(position, position.ActiveColor)
	isInCheck := position.IsInCheck(generator, position.ActiveColor)

	if position.isCheckmate(isInCheck, legalMoves) {
		if position.ActiveColor == White {
			return &Result{Outcome: BlackWin, Reason: Checkmate}
		} else {
			return &Result{Outcome: WhiteWin, Reason: Checkmate}
		}
	}
	if position.isStalemate(isInCheck, legalMoves) {
		return &Result{Outcome: Draw, Reason: Stalemate}
	}
	if position.isThreefoldRepetition(positions) {
		return &Result{Outcome: Draw, Reason: ThreefoldRepetition}
	}
	if position.isFiftyMoveRule() {
		return &Result{Outcome: Draw, Reason: FiftyMoveRule}
	}
	if position.isInsufficientMaterial() {
		return &Result{Outcome: Draw, Reason: InsufficientMaterial}
	}
	return nil
}

func (p *Position) isCheckmate(isInCheck bool, legalMoves []Move) bool {
	return isInCheck && len(legalMoves) == 0
}

func (p *Position) isStalemate(isInCheck bool, legalMoves []Move) bool {
	return !isInCheck && len(legalMoves) == 0
}

func (p *Position) isThreefoldRepetition(positions map[PositionKey]int) bool {
	return positions[p.Key()] >= 3
}

func (p *Position) isFiftyMoveRule() bool {
	return p.Halfmove >= 100
}

func (p *Position) isInsufficientMaterial() bool {
	if (p.WhitePawns | p.BlackPawns | p.WhiteRooks | p.BlackRooks | p.WhiteQueens | p.BlackQueens) != 0 {
		return false
	}

	whiteMinorPieces := p.WhiteKnights | p.WhiteBishops
	blackMinorPieces := p.BlackKnights | p.BlackBishops

	whiteCount := bits.OnesCount64(uint64(whiteMinorPieces))
	blackCount := bits.OnesCount64(uint64(blackMinorPieces))

	// K vs K
	if whiteCount == 0 && blackCount == 0 {
		return true
	}

	// K + N vs K or K + B vs K
	if whiteCount == 1 && blackCount == 0 {
		return true
	}
	if whiteCount == 0 && blackCount == 1 {
		return true
	}

	// K + B vs K + B with same color bishops
	whiteB := p.WhiteBishops
	blackB := p.BlackBishops
	if bits.OnesCount64(uint64(whiteB)) == 1 && bits.OnesCount64(uint64(blackB)) == 1 {
		whiteSq := popLSB(&whiteB)
		blackSq := popLSB(&blackB)
		if Square(whiteSq).Color() == Square(blackSq).Color() {
			return true
		}
	}

	return false
}
