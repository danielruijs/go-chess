package chess

import "encoding/json"

type MessageType string

const (
	MessageTypePosition  MessageType = "position"
	MessageTypeMove      MessageType = "move"
	MessageTypeJoinMatch MessageType = "join_match"
)

type WSMessage struct {
	Type MessageType     `json:"type"`
	Data json.RawMessage `json:"data"`
}
