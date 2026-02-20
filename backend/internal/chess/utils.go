package chess

import "math/bits"

func shift(b Bitboard, n int) Bitboard {
	if n > 0 {
		return b << n
	} else {
		return b >> (-n)
	}
}

// Clears the least significant bit and returns its index
func popLSB(b *Bitboard) int {
	index := bits.TrailingZeros64(uint64(*b))
	*b &= *b - 1
	return index
}
