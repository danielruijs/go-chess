package server

import (
	"context"
	"fmt"
	"go-chess/internal/chess"
	"log"
	"time"
)

const (
	clockCheckInterval     = 100 * time.Millisecond
	clockBroadcastInterval = 10 * time.Second
)

type Match struct {
	PublicID string
	Player1  *Player
	Player2  *Player

	Engine        *chess.Engine
	Clock         *MatchClock
	EventChan     chan Event
	MatchEnded    chan<- *Match
	DrawOfferedBy *Player

	metrics     *metrics
	eventStorer MatchEventStorer
}

func NewMatch(ctx context.Context, matchStorer MatchStorer, player1, player2 *Player, timeFormat TimeFormat, matchEnded chan<- *Match, metrics *metrics) (*Match, error) {
	publicID, err := GeneratePublicMatchID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate match public ID: %w", err)
	}

	eventStorer, err := matchStorer.CreateMatch(ctx, publicID, player1, player2, timeFormat)
	if err != nil {
		log.Printf("ERROR [NewMatch]: failed to create match record: %v", err)
	}

	return &Match{
		PublicID:    publicID,
		Player1:     player1,
		Player2:     player2,
		Engine:      chess.NewEngine(),
		Clock:       NewMatchClock(timeFormat),
		EventChan:   make(chan Event),
		MatchEnded:  matchEnded,
		metrics:     metrics,
		eventStorer: eventStorer,
	}, nil
}

func (m *Match) persistMatchEvent(ctx context.Context, event Event) {
	if m.eventStorer == nil {
		return
	}
	clockSnap := m.Clock.Snapshot(m.Engine.GetActiveColor())
	m.eventStorer.StoreMatchEvent(ctx, event, clockSnap)
}

func (m *Match) Run(ctx context.Context) {
	clockCheckTicker := time.NewTicker(clockCheckInterval)
	clockBroadcastTicker := time.NewTicker(clockBroadcastInterval)
	defer clockCheckTicker.Stop()
	defer clockBroadcastTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case event := <-m.EventChan:
			handler, ok := eventHandlers[event.Type]
			if !ok {
				log.Println("unknown event type:", event.Type)
				continue
			}

			m.metrics.recordMatchEvent(event.Type)
			result, valid := handler.Handle(m, event)
			if valid {
				m.persistMatchEvent(ctx, event)
			}

			if result != nil {
				m.end(ctx, result)
				return
			}

			m.sendPositionUpdate()
		case <-clockCheckTicker.C:
			if !m.Clock.IsRunning() {
				continue
			}
			loser := m.Clock.Advance(m.Engine.GetActiveColor())
			if loser != nil {
				m.end(ctx, getTimeoutResult(*loser))
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

func (m *Match) end(ctx context.Context, result *chess.Result) {
	m.Clock.Stop(m.Engine.GetActiveColor())
	m.Engine.ApplyResult(result)
	m.metrics.recordMatchFinished(result)

	if m.eventStorer != nil {
		m.eventStorer.StoreGameEndedEvent(ctx, result)
	}

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
