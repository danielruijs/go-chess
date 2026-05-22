package server

import (
	"go-chess/internal/chess"
	"log"
)

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

type EventHandler interface {
	Handle(m *Match, event Event) (ended bool)
}

var eventHandlers = map[eventType]EventHandler{
	EventTypeMove:        MoveEventHandler{},
	EventTypeGameStarted: GameStartedEventHandler{},
	EventTypeResign:      ResignEventHandler{},
	EventTypeOfferDraw:   OfferDrawEventHandler{},
	EventTypeRespondDraw: RespondDrawEventHandler{},
}

type MoveEventHandler struct{}

func (h MoveEventHandler) Handle(m *Match, event Event) (ended bool) {
	data, ok := event.Data.(MoveData)
	if !ok {
		log.Println("invalid move data format")
		m.metrics.recordWebsocketMessageError(MessageTypeMove, "invalid_payload")
		return false
	}

	loserByTimeout, err := m.Clock.BeforeMove(m.Engine.GetActiveColor())
	if err != nil {
		log.Println("failed to check match clock before move:", err)
		m.metrics.recordWebsocketMessageError(MessageTypeMove, "clock_error")
		return false
	}
	if loserByTimeout != nil {
		m.end(getTimeoutResult(*loserByTimeout))
		return true
	}

	move, err := moveDataToMove(data)
	if err != nil {
		log.Println("invalid move data:", err)
		m.metrics.recordWebsocketMessageError(MessageTypeMove, "invalid_move")
		return false
	}
	result, err := m.Engine.ApplyMove(move, event.Player.GetColor())
	if err != nil {
		if data.Promotion != nil {
			log.Printf("failed to apply move %s -> %s with promotion to %v: %v\n", data.From, data.To, *data.Promotion, err)
		} else {
			log.Printf("failed to apply move %s -> %s: %v\n", data.From, data.To, err)
		}
		m.metrics.recordWebsocketMessageError(MessageTypeMove, "illegal_move")
		return false
	}

	err = m.Clock.AfterMove(chess.GetOppositeColor(m.Engine.GetActiveColor()))
	if err != nil {
		log.Println("failed to update match clock after move:", err)
		return false
	}

	if result != nil {
		m.end(result)
		return true
	}

	return false
}

type GameStartedEventHandler struct{}

func (h GameStartedEventHandler) Handle(m *Match, event Event) (ended bool) {
	m.Clock.Start()
	for _, player := range []*Player{m.Player1, m.Player2} {
		startMatchData := m.getStartMatchData(player)
		err := player.Send(MessageTypeStartMatch, startMatchData)
		if err != nil {
			log.Printf("failed to send %s: %v", MessageTypeStartMatch, err)
			continue
		}
	}
	return false
}

type ResignEventHandler struct{}

func (h ResignEventHandler) Handle(m *Match, event Event) (ended bool) {
	var result *chess.Result
	if event.Player.GetColor() == chess.White {
		result = &chess.Result{Outcome: chess.BlackWin, Reason: chess.Resignation}
	} else {
		result = &chess.Result{Outcome: chess.WhiteWin, Reason: chess.Resignation}
	}
	m.end(result)
	return true
}

type OfferDrawEventHandler struct{}

func (h OfferDrawEventHandler) Handle(m *Match, event Event) (ended bool) {
	if m.DrawOfferedBy != nil {
		log.Println("draw offer already pending")
		m.metrics.recordWebsocketMessageError(MessageTypeOfferDraw, "duplicate_offer")
		return false
	}
	m.DrawOfferedBy = event.Player

	opponent := m.getOpponent(event.Player)
	err := opponent.Send(MessageTypeDrawOffered, nil)
	if err != nil {
		log.Printf("failed to send %s: %v", MessageTypeDrawOffered, err)
	}
	return false
}

type RespondDrawEventHandler struct{}

func (h RespondDrawEventHandler) Handle(m *Match, event Event) (ended bool) {
	data, ok := event.Data.(RespondDrawData)
	if !ok {
		log.Println("invalid respond draw data format")
		m.metrics.recordWebsocketMessageError(MessageTypeRespondDraw, "invalid_payload")
		return false
	}
	if m.DrawOfferedBy == nil {
		log.Println("no draw offer to respond to")
		m.metrics.recordWebsocketMessageError(MessageTypeRespondDraw, "no_pending_offer")
		return false
	}
	if m.DrawOfferedBy == event.Player {
		log.Println("player cannot respond to their own draw offer")
		m.metrics.recordWebsocketMessageError(MessageTypeRespondDraw, "self_response")
		return false
	}
	if !data.Accept {
		// notify opponent that the draw offer was declined
		opponent := m.getOpponent(event.Player)
		err := opponent.Send(MessageTypeDrawDeclined, nil)
		if err != nil {
			log.Printf("failed to send %s: %v", MessageTypeDrawDeclined, err)
		}
		m.DrawOfferedBy = nil
		return false
	}
	// accepted draw
	result := &chess.Result{Outcome: chess.Draw, Reason: chess.AgreedDraw}
	m.end(result)
	return true
}
