package chess

func shift(b Bitboard, n int) Bitboard {
	if n > 0 {
		return b << n
	} else {
		return b >> (-n)
	}
}
