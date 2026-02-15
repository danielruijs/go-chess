package chess

import "encoding/json"

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
	Position Position `json:"position"`
}

type MoveData struct {
	Move Move `json:"move"`
}

type JoinMatchData struct {
	PlayerName string `json:"playerName"`
}

type ErrorData struct {
	Message string `json:"message"`
}
