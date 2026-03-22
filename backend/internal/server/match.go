package server

import (
	"fmt"
	"go-chess/internal/chess"
	"log"
	"time"
)

const (
	clockCheckInterval     = 100 * time.Millisecond
	clockBroadcastInterval = 1 * time.Second
)

type Match struct {
	Player1 *Player
	Player2 *Player

	Engine     *chess.Engine
	Clock      *MatchClock
	EventChan  chan Event
	MatchEnded chan<- *Match

	DrawOfferedBy *Player
}

func (m *Match) Run() {
	clockCheckTicker := time.NewTicker(clockCheckInterval)
	clockBroadcastTicker := time.NewTicker(clockBroadcastInterval)
	defer clockCheckTicker.Stop()
	defer clockBroadcastTicker.Stop()

	for {
		select {
		case event := <-m.EventChan:
			switch event.Type {
			case EventTypeMove:
				ended := m.handleMoveEvent(event)
				if ended {
					return
				}
			case EventTypeGameStarted:
				m.handleGameStartedEvent()
			case EventTypeResign:
				ended := m.handleResignEvent(event)
				if ended {
					return
				}
			case EventTypeOfferDraw:
				m.handleOfferDrawEvent(event)
			case EventTypeRespondDraw:
				ended := m.handleRespondDrawEvent(event)
				if ended {
					return
				}
			}

			err := m.sendPositionUpdate(true)
			if err != nil {
				log.Println("failed to send position update:", err)
			}
		case <-clockCheckTicker.C:
			if !m.Clock.IsRunning() {
				continue
			}
			loser := m.Clock.Advance(m.Engine.GetActiveColor())
			if loser != nil {
				m.end(getTimeoutResult(*loser))
				return
			}
		case <-clockBroadcastTicker.C:
			if !m.Clock.IsRunning() {
				continue
			}
			err := m.sendPositionUpdate(false)
			if err != nil {
				log.Println("failed to send position update:", err)
			}
		}
	}
}

func (m *Match) sendPositionUpdate(isCritical bool) error {
	for _, player := range []*Player{m.Player1, m.Player2} {
		legalMovesList := m.Engine.GetLegalMoves(player.GetColor())
		boardData := BoardData{
			Board:       m.Engine.GetBoard(),
			LegalMoves:  moveListToLegalMoves(legalMovesList),
			PGN:         m.Engine.GetPGN(),
			ActiveColor: m.Engine.GetActiveColor(),
			Clock:       m.Clock.Snapshot(m.Engine.GetActiveColor()),
		}
		if isCritical {
			err := player.SendCritical(MessageTypeBoard, boardData)
			if err != nil {
				return fmt.Errorf("failed to send board data to %s: %v", player.Name, err)
			}
		} else {
			player.SendInformational(MessageTypeBoard, boardData)
		}
	}
	return nil
}

func (m *Match) sendFinalPositionUpdate() error {
	for _, player := range []*Player{m.Player1, m.Player2} {
		boardData := BoardData{
			Board:       m.Engine.GetBoard(),
			LegalMoves:  map[string][]LegalMove{}, // no legal moves
			PGN:         m.Engine.GetPGN(),
			ActiveColor: m.Engine.GetActiveColor(),
			Clock:       m.Clock.Snapshot(m.Engine.GetActiveColor()),
		}
		err := player.SendCritical(MessageTypeBoard, boardData)
		if err != nil {
			return fmt.Errorf("failed to send final board data to %s: %v", player.Name, err)
		}
	}
	return nil
}

func (m *Match) end(result *chess.Result) {
	m.Clock.Stop(m.Engine.GetActiveColor())
	m.Engine.ApplyResult(result)
	err := m.sendFinalPositionUpdate()
	if err != nil {
		log.Println("failed to send final position update:", err)
	}
	m.sendMatchEnd(*result)
	close(m.EventChan)
	m.MatchEnded <- m
	log.Printf("ended match between %s and %s with result: %s\n", m.Player1.Name, m.Player2.Name, result.Outcome)
}

func getTimeoutResult(loser chess.Color) *chess.Result {
	if loser == chess.White {
		return &chess.Result{Outcome: chess.BlackWin, Reason: chess.Timeout}
	}
	return &chess.Result{Outcome: chess.WhiteWin, Reason: chess.Timeout}
}

func (m *Match) sendMatchEnd(result chess.Result) {
	resultData := EndMatchData{
		Result: result,
	}
	for _, player := range []*Player{m.Player1, m.Player2} {
		err := player.SendCritical(MessageTypeEndMatch, resultData)
		if err != nil {
			log.Printf("failed to send end match message to %s: %v", player.Name, err)
		}
	}

}

func (m *Match) getPlayerByColor(color chess.Color) *Player {
	if m.Player1.GetColor() == color {
		return m.Player1
	}
	return m.Player2
}
