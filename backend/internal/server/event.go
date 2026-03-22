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

func (m *Match) handleMoveEvent(event Event) (ended bool) {
	data, ok := event.Data.(MoveData)
	if !ok {
		log.Println("invalid move data format")
		return false
	}

	loserByTimeout, err := m.Clock.BeforeMove(m.Engine.GetActiveColor())
	if err != nil {
		log.Println("failed to check match clock before move:", err)
		return false
	}
	if loserByTimeout != nil {
		m.end(getTimeoutResult(*loserByTimeout))
		return true
	}

	move, err := moveDataToMove(data)
	if err != nil {
		log.Println("invalid move data:", err)
		return false
	}
	result, err := m.Engine.ApplyMove(move, event.Player.GetColor())
	if err != nil {
		if data.Promotion != nil {
			log.Printf("failed to apply move %s -> %s with promotion to %v: %v\n", data.From, data.To, *data.Promotion, err)
		} else {
			log.Printf("failed to apply move %s -> %s: %v\n", data.From, data.To, err)
		}
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

func (m *Match) handleGameStartedEvent() {
	m.Clock.Start()
	for _, player := range []*Player{m.Player1, m.Player2} {
		startMatchData := StartMatchData{
			Color:           player.GetColor(),
			WhitePlayerName: m.getPlayerByColor(chess.White).Name,
			BlackPlayerName: m.getPlayerByColor(chess.Black).Name,
			Clock:           m.Clock.Snapshot(m.Engine.GetActiveColor()),
		}
		err := player.SendCritical(MessageTypeStartMatch, startMatchData)
		if err != nil {
			log.Printf("failed to send start match message to %s: %v", player.Name, err)
			continue
		}
	}
}

func (m *Match) handleResignEvent(event Event) (ended bool) {
	var result *chess.Result
	if event.Player.GetColor() == chess.White {
		result = &chess.Result{Outcome: chess.BlackWin, Reason: chess.Resignation}
	} else {
		result = &chess.Result{Outcome: chess.WhiteWin, Reason: chess.Resignation}
	}
	m.end(result)
	return true
}

func (m *Match) handleOfferDrawEvent(event Event) {
	if m.DrawOfferedBy != nil {
		log.Println("draw offer already pending")
		return
	}
	m.DrawOfferedBy = event.Player

	opponent := m.getOpponent(event.Player)
	err := opponent.SendCritical(MessageTypeDrawOffered, nil)
	if err != nil {
		log.Printf("failed to send draw offered message to %s: %v", opponent.Name, err)
	}
}

func (m *Match) handleRespondDrawEvent(event Event) (ended bool) {
	data, ok := event.Data.(RespondDrawData)
	if !ok {
		log.Println("invalid respond draw data format")
		return false
	}
	if m.DrawOfferedBy == nil {
		log.Println("no draw offer to respond to")
		return false
	}
	if m.DrawOfferedBy == event.Player {
		log.Println("player cannot respond to their own draw offer")
		return false
	}
	if !data.Accept {
		// notify opponent that the draw offer was declined
		opponent := m.getOpponent(event.Player)
		err := opponent.SendCritical(MessageTypeDrawDeclined, nil)
		if err != nil {
			log.Printf("failed to send draw declined message to %s: %v", opponent.Name, err)
		}
		m.DrawOfferedBy = nil
		return false
	}
	// accepted draw
	result := &chess.Result{Outcome: chess.Draw, Reason: chess.AgreedDraw}
	m.end(result)
	return true
}

func (m *Match) getOpponent(p *Player) *Player {
	if p == m.Player1 {
		return m.Player2
	}
	return m.Player1
}
