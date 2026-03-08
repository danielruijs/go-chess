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
				data, ok := event.Data.(MoveData)
				if !ok {
					log.Println("invalid move data format")
					continue
				}

				loserByTimeout, err := m.Clock.BeforeMove()
				if err != nil {
					log.Println("failed to check match clock before move:", err)
					continue
				}
				if loserByTimeout != nil {
					m.end(getTimeoutResult(*loserByTimeout))
					return
				}

				move, err := moveDataToMove(data)
				if err != nil {
					log.Println("invalid move data:", err)
					continue
				}
				result, err := m.Engine.ApplyMove(move, event.Player.GetColor())
				if err != nil {
					if data.Promotion != nil {
						log.Printf("failed to apply move %s -> %s with promotion to %v: %v\n", data.From, data.To, *data.Promotion, err)
					} else {
						log.Printf("failed to apply move %s -> %s: %v\n", data.From, data.To, err)
					}
					continue
				}

				err = m.Clock.AfterMove()
				if err != nil {
					log.Println("failed to update match clock after move:", err)
					continue
				}

				if result != nil {
					m.end(result)
					return
				}
			case EventTypeGameStarted:
				m.Clock.Start()
				for _, player := range []*Player{m.Player1, m.Player2} {
					startMatchData := StartMatchData{
						Color:           player.GetColor(),
						WhitePlayerName: m.getPlayerByColor(chess.White).Name,
						BlackPlayerName: m.getPlayerByColor(chess.Black).Name,
						Clock:           m.Clock.Snapshot(),
					}
					err := player.SendCritical(MessageTypeStartMatch, startMatchData)
					if err != nil {
						log.Printf("failed to send start match message to %s: %v", player.Name, err)
						continue
					}
				}
			case EventTypeResign:
				var result *chess.Result
				if event.Player.GetColor() == chess.White {
					result = &chess.Result{Outcome: chess.BlackWin, Reason: chess.Resignation}
				} else {
					result = &chess.Result{Outcome: chess.WhiteWin, Reason: chess.Resignation}
				}
				m.end(result)
				return
			case EventTypeOfferDraw:
				if m.DrawOfferedBy != nil {
					log.Println("draw offer already pending")
					continue
				}
				m.DrawOfferedBy = event.Player

				receivingPlayer := m.Player1
				if event.Player == m.Player1 {
					receivingPlayer = m.Player2
				}
				err := receivingPlayer.SendCritical(MessageTypeDrawOffered, nil)
				if err != nil {
					log.Printf("failed to send draw offered message to %s: %v", receivingPlayer.Name, err)
				}
			case EventTypeRespondDraw:
				data, ok := event.Data.(RespondDrawData)
				if !ok {
					log.Println("invalid respond draw data format")
					continue
				}
				if m.DrawOfferedBy == nil {
					log.Println("no draw offer to respond to")
					continue
				}
				if m.DrawOfferedBy == event.Player {
					log.Println("player cannot respond to their own draw offer")
					continue
				}
				if !data.Accept {
					// notify opponent that the draw offer was declined
					receivingPlayer := m.Player1
					if event.Player == m.Player1 {
						receivingPlayer = m.Player2
					}
					err := receivingPlayer.SendCritical(MessageTypeDrawDeclined, nil)
					if err != nil {
						log.Printf("failed to send draw declined message to %s: %v", receivingPlayer.Name, err)
					}
					m.DrawOfferedBy = nil
					continue
				}
				// accepted draw
				result := &chess.Result{Outcome: chess.Draw, Reason: chess.AgreedDraw}
				m.end(result)
				return
			}

			err := m.sendPositionUpdate(true)
			if err != nil {
				log.Println("failed to send position update:", err)
			}
		case <-clockCheckTicker.C:
			if !m.Clock.IsRunning() {
				continue
			}
			loser := m.Clock.Advance()
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
			Board:      m.Engine.GetBoard(),
			LegalMoves: moveListToLegalMoves(legalMovesList),
			PGN:        m.Engine.GetPGN(),
			Clock:      m.Clock.Snapshot(),
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
			Board:      m.Engine.GetBoard(),
			LegalMoves: map[string][]LegalMove{}, // no legal moves
			PGN:        m.Engine.GetPGN(),
			Clock:      m.Clock.Snapshot(),
		}
		err := player.SendCritical(MessageTypeBoard, boardData)
		if err != nil {
			return fmt.Errorf("failed to send final board data to %s: %v", player.Name, err)
		}
	}
	return nil
}

func (m *Match) end(result *chess.Result) {
	m.Clock.Stop()
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
