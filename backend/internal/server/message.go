package server

import (
	"encoding/json"
	"go-chess/internal/chess"
)

type MessageType string

const (
	// outbound
	MessageTypeBoard             MessageType = "board"
	MessageTypeMatchmakingUpdate MessageType = "matchmaking_update"
	MessageTypeStartMatch        MessageType = "start_match"
	MessageTypeEndMatch          MessageType = "end_match"
	MessageTypeDrawOffered       MessageType = "draw_offered"
	MessageTypeDrawDeclined      MessageType = "draw_declined"
	// inbound
	MessageTypeJoinMatch   MessageType = "join_match"
	MessageTypeMove        MessageType = "move"
	MessageTypeResign      MessageType = "resign"
	MessageTypeOfferDraw   MessageType = "offer_draw"
	MessageTypeRespondDraw MessageType = "respond_draw"
)

type WSMessage struct {
	Type MessageType     `json:"type"`
	Data json.RawMessage `json:"data"`
}

// outbound
type LegalMove struct {
	To        string           `json:"to"`
	Promotion *chess.PieceType `json:"promotion,omitempty"`
}

type BoardData struct {
	Board      chess.Board            `json:"board"`
	LegalMoves map[string][]LegalMove `json:"legalMoves"`
	PGN        chess.PGN              `json:"pgn"`
}

type MatchmakingUpdateData struct {
	QueueLength int  `json:"queueLength"`
	InQueue     bool `json:"inQueue"`
}

type StartMatchData struct {
	Color           chess.Color `json:"color"`
	WhitePlayerName string      `json:"whitePlayerName"`
	BlackPlayerName string      `json:"blackPlayerName"`
}

type EndMatchData struct {
	Result chess.Result `json:"result"`
}

// inbound
type JoinMatchData struct {
	PlayerName string `json:"playerName"`
}

type MoveData struct {
	From      string           `json:"from"`
	To        string           `json:"to"`
	Promotion *chess.PieceType `json:"promotion,omitempty"`
}

type RespondDrawData struct {
	Accept bool `json:"accept"`
}
