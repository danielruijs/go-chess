package server

import (
	"go-chess/internal/chess"
	"log"
	"time"
)

const (
	clockCheckInterval     = 100 * time.Millisecond
	clockBroadcastInterval = 10 * time.Second
)

type Match struct {
	Player1 *Player
	Player2 *Player

	Engine     *chess.Engine
	Clock      *MatchClock
	EventChan  chan Event
	MatchEnded chan<- *Match

	DrawOfferedBy *Player
	metrics       *metrics
}

func NewMatch(player1, player2 *Player, timeFormat TimeFormat, matchEnded chan<- *Match, metrics *metrics) *Match {
	return &Match{
		Player1:    player1,
		Player2:    player2,
		Engine:     chess.NewEngine(),
		Clock:      NewMatchClock(timeFormat),
		EventChan:  make(chan Event),
		MatchEnded: matchEnded,
		metrics:    metrics,
	}
}

func (m *Match) Run() {
	clockCheckTicker := time.NewTicker(clockCheckInterval)
	clockBroadcastTicker := time.NewTicker(clockBroadcastInterval)
	defer clockCheckTicker.Stop()
	defer clockBroadcastTicker.Stop()

	for {
		select {
		case event := <-m.EventChan:
			handler, ok := eventHandlers[event.Type]
			if !ok {
				log.Println("unknown event type:", event.Type)
				continue
			}

			m.metrics.recordMatchEvent(event.Type)
			ended := handler.Handle(m, event)
			if ended {
				return
			}

			m.sendPositionUpdate()
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
			m.sendPositionUpdate()
		}
	}
}

func (m *Match) sendPositionUpdate() {
	for _, player := range []*Player{m.Player1, m.Player2} {
		boardData := m.getBoardData(player)
		player.Send(MessageTypeBoard, boardData)
	}
}

func (m *Match) sendCurrentState(player *Player) {
	player.Send(MessageTypeStartMatch, m.getStartMatchData(player))
	player.Send(MessageTypeBoard, m.getBoardData(player))
}

func (m *Match) getBoardData(player *Player) BoardData {
	legalMovesList := m.Engine.GetLegalMoves(player.GetColor())
	return BoardData{
		Board:       m.Engine.GetBoard(),
		LegalMoves:  moveListToLegalMoves(legalMovesList),
		PGN:         m.Engine.GetPGN(),
		ActiveColor: m.Engine.GetActiveColor(),
		Clock:       m.Clock.Snapshot(m.Engine.GetActiveColor()),
	}
}

func (m *Match) getStartMatchData(player *Player) StartMatchData {
	return StartMatchData{
		Color:           player.GetColor(),
		WhitePlayerName: m.getPlayerByColor(chess.White).DisplayName,
		BlackPlayerName: m.getPlayerByColor(chess.Black).DisplayName,
		Clock:           m.Clock.Snapshot(m.Engine.GetActiveColor()),
	}
}

func (m *Match) sendFinalPositionUpdate() {
	for _, player := range []*Player{m.Player1, m.Player2} {
		boardData := BoardData{
			Board:       m.Engine.GetBoard(),
			LegalMoves:  map[string][]LegalMove{}, // no legal moves
			PGN:         m.Engine.GetPGN(),
			ActiveColor: m.Engine.GetActiveColor(),
			Clock:       m.Clock.Snapshot(m.Engine.GetActiveColor()),
		}
		player.Send(MessageTypeBoard, boardData)
	}
}

func (m *Match) end(result *chess.Result) {
	m.Clock.Stop(m.Engine.GetActiveColor())
	m.Engine.ApplyResult(result)
	m.metrics.recordMatchFinished(result)
	m.sendFinalPositionUpdate()
	m.sendMatchEnd(*result)
	close(m.EventChan)
	m.MatchEnded <- m
	log.Printf("ended match between %s and %s with result: %s\n", m.Player1.DisplayName, m.Player2.DisplayName, result.Outcome)
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
		player.Send(MessageTypeEndMatch, resultData)
	}
}

func (m *Match) getPlayerByColor(color chess.Color) *Player {
	if m.Player1.GetColor() == color {
		return m.Player1
	}
	return m.Player2
}

func (m *Match) getOpponent(p *Player) *Player {
	if p == m.Player1 {
		return m.Player2
	}
	return m.Player1
}
