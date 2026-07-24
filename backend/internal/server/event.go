package server

import (
	"go-chess/internal/chess"
	"log"
)

type eventType string
type eventData any

const (
	EventTypeMove              eventType = "move"
	EventTypeGameStarted       eventType = "game_started"
	EventTypeResign            eventType = "resign"
	EventTypeOfferDraw         eventType = "offer_draw"
	EventTypeRespondDraw       eventType = "respond_draw"
	EventTypePlayerReconnected eventType = "player_reconnected"
)

type Event struct {
	Player *Player
	Type   eventType
	Data   eventData
}

type EventHandler interface {
	Handle(m *Match, event Event) (result *chess.Result, valid bool)
}

var eventHandlers = map[eventType]EventHandler{
	EventTypeMove:              MoveEventHandler{},
	EventTypeGameStarted:       GameStartedEventHandler{},
	EventTypeResign:            ResignEventHandler{},
	EventTypeOfferDraw:         OfferDrawEventHandler{},
	EventTypeRespondDraw:       RespondDrawEventHandler{},
	EventTypePlayerReconnected: PlayerReconnectedEventHandler{},
}

type MoveEventHandler struct{}

func (h MoveEventHandler) Handle(m *Match, event Event) (result *chess.Result, valid bool) {
	data, ok := event.Data.(MoveData)
	if !ok {
		log.Println("invalid move data format")
		m.metrics.recordWebsocketMessageError(MessageTypeMove, "invalid_payload")
		return nil, false
	}

	loserByTimeout, err := m.Clock.BeforeMove(m.Engine.GetActiveColor())
	if err != nil {
		log.Println("failed to check match clock before move:", err)
		m.metrics.recordWebsocketMessageError(MessageTypeMove, "clock_error")
		return nil, false
	}
	if loserByTimeout != nil {
		return getTimeoutResult(*loserByTimeout), false
	}

	move, err := moveDataToMove(data)
	if err != nil {
		log.Println("invalid move data:", err)
		m.metrics.recordWebsocketMessageError(MessageTypeMove, "invalid_move")
		return nil, false
	}
	result, err = m.Engine.ApplyMove(move, event.Player.GetColor())
	if err != nil {
		if data.Promotion != nil {
			log.Printf("failed to apply move %s -> %s with promotion to %v: %v", data.From, data.To, *data.Promotion, err)
		} else {
			log.Printf("failed to apply move %s -> %s: %v", data.From, data.To, err)
		}
		m.metrics.recordWebsocketMessageError(MessageTypeMove, "illegal_move")
		return nil, false
	}

	err = m.Clock.AfterMove(chess.GetOppositeColor(m.Engine.GetActiveColor()))
	if err != nil {
		log.Println("failed to update match clock after move:", err)
		return nil, false
	}

	return result, true
}

type GameStartedEventHandler struct{}

func (h GameStartedEventHandler) Handle(m *Match, event Event) (result *chess.Result, valid bool) {
	m.Clock.Start()
	for _, player := range []*Player{m.Player1, m.Player2} {
		startMatchData := m.getStartMatchData(player)
		player.Send(MessageTypeStartMatch, startMatchData)
	}
	return nil, true
}

// re-sends the current match state to a specific player
type PlayerReconnectedEventHandler struct{}

func (h PlayerReconnectedEventHandler) Handle(m *Match, event Event) (result *chess.Result, valid bool) {
	if event.Player != nil {
		m.sendCurrentState(event.Player)
	}
	return nil, true
}

type ResignEventHandler struct{}

func (h ResignEventHandler) Handle(m *Match, event Event) (result *chess.Result, valid bool) {
	if event.Player.GetColor() == chess.White {
		return &chess.Result{Outcome: chess.BlackWin, Reason: chess.Resignation}, false
	}
	return &chess.Result{Outcome: chess.WhiteWin, Reason: chess.Resignation}, false
}

type OfferDrawEventHandler struct{}

func (h OfferDrawEventHandler) Handle(m *Match, event Event) (result *chess.Result, valid bool) {
	if m.DrawOfferedBy != nil {
		log.Println("draw offer already pending")
		m.metrics.recordWebsocketMessageError(MessageTypeOfferDraw, "duplicate_offer")
		return nil, false
	}
	m.DrawOfferedBy = event.Player

	opponent := m.getOpponent(event.Player)
	opponent.Send(MessageTypeDrawOffered, nil)
	return nil, true
}

type RespondDrawEventHandler struct{}

func (h RespondDrawEventHandler) Handle(m *Match, event Event) (result *chess.Result, valid bool) {
	data, ok := event.Data.(RespondDrawData)
	if !ok {
		log.Println("invalid respond draw data format")
		m.metrics.recordWebsocketMessageError(MessageTypeRespondDraw, "invalid_payload")
		return nil, false
	}
	if m.DrawOfferedBy == nil {
		log.Println("no draw offer to respond to")
		m.metrics.recordWebsocketMessageError(MessageTypeRespondDraw, "no_pending_offer")
		return nil, false
	}
	if m.DrawOfferedBy == event.Player {
		log.Println("player cannot respond to their own draw offer")
		m.metrics.recordWebsocketMessageError(MessageTypeRespondDraw, "self_response")
		return nil, false
	}
	if !data.Accept {
		// notify opponent that the draw offer was declined
		opponent := m.getOpponent(event.Player)
		opponent.Send(MessageTypeDrawDeclined, nil)
		m.DrawOfferedBy = nil
		return nil, true
	}
	// accepted draw
	return &chess.Result{Outcome: chess.Draw, Reason: chess.AgreedDraw}, false
}
