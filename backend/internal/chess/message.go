package chess

type messageType string
type messageData any

const (
	MessageTypePosition  messageType = "position"
	MessageTypeMove      messageType = "move"
	MessageTypeJoinMatch messageType = "join_match"
)

type WSMessage struct {
	Type messageType `json:"type"`
	Data messageData `json:"data"`
}
