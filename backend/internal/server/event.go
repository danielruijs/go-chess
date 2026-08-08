package server

import (
	"errors"
	"fmt"
	"go-chess/internal/chess"
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
	Handle(m *Match, event Event) (result *chess.Result, err error)
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

func (h MoveEventHandler) Handle(m *Match, event Event) (result *chess.Result, err error) {
	data, ok := event.Data.(MoveData)
	if !ok {
		return nil, errors.New("invalid move data format")
	}

	loserByTimeout, err := m.Clock.BeforeMove(m.Engine.GetActiveColor())
	if err != nil {
		return nil, fmt.Errorf("failed to check match clock before move: %w", err)
	}
	if loserByTimeout != nil {
		return getTimeoutResult(*loserByTimeout), nil
	}

	move, err := moveDataToMove(data)
	if err != nil {
		return nil, fmt.Errorf("invalid move data: %w", err)
	}
	result, err = m.Engine.ApplyMove(move, event.Player.GetColor())
	if err != nil {
		if data.Promotion != nil {
			return nil, fmt.Errorf("failed to apply move %s -> %s with promotion to %v: %w", data.From, data.To, *data.Promotion, err)
		}
		return nil, fmt.Errorf("failed to apply move %s -> %s: %w", data.From, data.To, err)
	}

	err = m.Clock.AfterMove(chess.GetOppositeColor(m.Engine.GetActiveColor()))
	if err != nil {
		return nil, fmt.Errorf("failed to update match clock after move: %w", err)
	}

	return result, nil
}

type GameStartedEventHandler struct{}

func (h GameStartedEventHandler) Handle(m *Match, event Event) (result *chess.Result, err error) {
	m.Clock.Start()
	for _, player := range []*Player{m.Player1, m.Player2} {
		startMatchData := m.getStartMatchData(player)
		player.Send(MessageTypeStartMatch, startMatchData)
	}
	return nil, nil
}

// re-sends the current match state to a specific player
type PlayerReconnectedEventHandler struct{}

func (h PlayerReconnectedEventHandler) Handle(m *Match, event Event) (result *chess.Result, err error) {
	if event.Player != nil {
		m.sendCurrentState(event.Player)
	}
	return nil, nil
}

type ResignEventHandler struct{}

func (h ResignEventHandler) Handle(m *Match, event Event) (result *chess.Result, err error) {
	if event.Player != nil && event.Player.GetColor() == chess.White {
		return &chess.Result{Outcome: chess.BlackWin, Reason: chess.Resignation}, nil
	}
	return &chess.Result{Outcome: chess.WhiteWin, Reason: chess.Resignation}, nil
}

type OfferDrawEventHandler struct{}

func (h OfferDrawEventHandler) Handle(m *Match, event Event) (result *chess.Result, err error) {
	if m.DrawOfferedBy != nil {
		return nil, errors.New("draw offer already pending")
	}
	m.DrawOfferedBy = event.Player

	opponent := m.getOpponent(event.Player)
	opponent.Send(MessageTypeDrawOffered, nil)
	return nil, nil
}

type RespondDrawEventHandler struct{}

func (h RespondDrawEventHandler) Handle(m *Match, event Event) (result *chess.Result, err error) {
	data, ok := event.Data.(RespondDrawData)
	if !ok {
		return nil, errors.New("invalid respond draw data format")
	}
	if m.DrawOfferedBy == nil {
		return nil, errors.New("no draw offer to respond to")
	}
	if m.DrawOfferedBy == event.Player {
		return nil, errors.New("player cannot respond to their own draw offer")
	}
	if !data.Accept {
		// notify opponent that the draw offer was declined
		opponent := m.getOpponent(event.Player)
		opponent.Send(MessageTypeDrawDeclined, nil)
		m.DrawOfferedBy = nil
		return nil, nil
	}
	// accepted draw
	return &chess.Result{Outcome: chess.Draw, Reason: chess.AgreedDraw}, nil
}
