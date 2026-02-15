package chess

type eventType string
type eventData any

const (
	EventTypeMove        eventType = "move"
	EventTypeGameStarted eventType = "game_started"
)

type Event struct {
	Player *Player
	Type   eventType
	Data   eventData
}
