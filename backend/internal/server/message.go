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

type ClockData struct {
	WhiteTimeMs int64 `json:"whiteTimeMs"`
	BlackTimeMs int64 `json:"blackTimeMs"`
	IncrementMs int64 `json:"incrementMs"`
}

type BoardData struct {
	Board      chess.Board            `json:"board"`
	LegalMoves map[string][]LegalMove `json:"legalMoves"`
	PGN        chess.PGN              `json:"pgn"`
	Clock      ClockData              `json:"clock"`
}

type QueueData struct {
	QueueLength int  `json:"queueLength"`
	InQueue     bool `json:"inQueue"`
}

type MatchmakingUpdateData struct {
	Queues map[TimeFormat]QueueData `json:"queues"`
}

type StartMatchData struct {
	Color           chess.Color `json:"color"`
	WhitePlayerName string      `json:"whitePlayerName"`
	BlackPlayerName string      `json:"blackPlayerName"`
	Clock           ClockData   `json:"clock"`
}

type EndMatchData struct {
	Result chess.Result `json:"result"`
}

// inbound
type JoinMatchData struct {
	PlayerName string     `json:"playerName"`
	TimeFormat TimeFormat `json:"timeFormat"`
}

type MoveData struct {
	From      string           `json:"from"`
	To        string           `json:"to"`
	Promotion *chess.PieceType `json:"promotion,omitempty"`
}

type RespondDrawData struct {
	Accept bool `json:"accept"`
}
