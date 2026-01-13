package chess

type eventType string
type eventData any

const (
	EventTypeMove        eventType = "move"
	EventTypeGameStarted eventType = "game_started"
)

type Event struct {
	Player *Player   `json:"player"`
	Type   eventType `json:"type"`
	Data   eventData `json:"data"`
}
