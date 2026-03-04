package server

type eventType string
type eventData any

const (
	EventTypeMove        eventType = "move"
	EventTypeGameStarted eventType = "game_started"
	EventTypeResign      eventType = "resign"
	EventTypeOfferDraw   eventType = "offer_draw"
	EventTypeRespondDraw eventType = "respond_draw"
)

type Event struct {
	Player *Player
	Type   eventType
	Data   eventData
}
