package chess

type messageType string
type messageData any

const (
	MessageTypePosition messageType = "position"
)

// var messageDataMap = map[messageType]messageData{
// 	MessageTypePosition: Position{},
// }

type WSMessage struct {
	Type messageType `json:"type"`
	Data messageData `json:"data"`
}
