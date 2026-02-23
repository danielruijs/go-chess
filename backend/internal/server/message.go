package server

import (
	"encoding/json"
	"go-chess/internal/chess"
)

type MessageType string

const (
	MessageTypeBoard      MessageType = "board"
	MessageTypeMove       MessageType = "move"
	MessageTypeJoinMatch  MessageType = "join_match"
	MessageTypeStartMatch MessageType = "start_match"
)

type WSMessage struct {
	Type MessageType     `json:"type"`
	Data json.RawMessage `json:"data"`
}

type LegalMove struct {
	To        string           `json:"to"`
	Promotion *chess.PieceType `json:"promotion,omitempty"`
}

type BoardData struct {
	Board      chess.Board            `json:"board"`
	LegalMoves map[string][]LegalMove `json:"legalMoves"`
}

type MoveData struct {
	From      string           `json:"from"`
	To        string           `json:"to"`
	Promotion *chess.PieceType `json:"promotion,omitempty"`
}

type JoinMatchData struct {
	PlayerName string `json:"playerName"`
}

type StartMatchData struct {
	WhitePlayerName string `json:"whitePlayerName"`
	BlackPlayerName string `json:"blackPlayerName"`
}
