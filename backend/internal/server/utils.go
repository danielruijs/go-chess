package server

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"

	"go-chess/internal/chess"
)

const (
	base62Alphabet   = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	publicMatchIDLen = 12
)

// GeneratePublicMatchID creates a random 12-character base62 string for match routing.
func GeneratePublicMatchID() (string, error) {
	result := make([]byte, publicMatchIDLen)
	alphabetLen := big.NewInt(int64(len(base62Alphabet)))

	for i := range publicMatchIDLen {
		n, err := rand.Int(rand.Reader, alphabetLen)
		if err != nil {
			return "", err
		}
		result[i] = base62Alphabet[n.Int64()]
	}

	return string(result), nil
}

func IsValidPublicMatchID(id string) bool {
	if len(id) != publicMatchIDLen {
		return false
	}
	for _, c := range id {
		if !strings.ContainsRune(base62Alphabet, c) {
			return false
		}
	}
	return true
}

func moveDataToMove(data MoveData) (chess.Move, error) {
	from, err := chess.StrToSquare(data.From)
	if err != nil {
		return chess.Move{}, fmt.Errorf("invalid from square: %w", err)
	}
	to, err := chess.StrToSquare(data.To)
	if err != nil {
		return chess.Move{}, fmt.Errorf("invalid to square: %w", err)
	}

	return chess.Move{
		From:      from,
		To:        to,
		Promotion: data.Promotion,
	}, nil
}

func moveListToLegalMoves(moves []chess.Move) map[string][]LegalMove {
	legalMoves := make(map[string][]LegalMove)
	for _, move := range moves {
		fromStr := chess.SquareToStr(move.From)
		legalMove := LegalMove{
			To:        chess.SquareToStr(move.To),
			Promotion: move.Promotion,
		}
		legalMoves[fromStr] = append(legalMoves[fromStr], legalMove)
	}
	return legalMoves
}
