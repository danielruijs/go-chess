package server

import (
	"encoding/json"
	"go-chess/internal/chess"
)

type MessageType string

const (
	MessageTypePosition  MessageType = "position"
	MessageTypeMove      MessageType = "move"
	MessageTypeJoinMatch MessageType = "join_match"
	MessageTypeError     MessageType = "error"
)

type WSMessage struct {
	Type MessageType     `json:"type"`
	Data json.RawMessage `json:"data"`
}

type PositionData struct {
	Position chess.Position `json:"position"`
}

type MoveData struct {
	Move chess.Move `json:"move"`
}

type JoinMatchData struct {
	PlayerName string `json:"playerName"`
}

type ErrorData struct {
	Message string `json:"message"`
}
