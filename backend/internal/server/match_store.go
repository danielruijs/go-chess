package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"go-chess/internal/chess"
	"go-chess/internal/db/sqlc"

	"github.com/jackc/pgx/v5/pgtype"
)

type MatchStorer interface {
	CreateMatch(ctx context.Context, publicID string, player1, player2 *Player, timeFormat TimeFormat) (MatchEventStorer, error)
	Close()
}

type MatchEventStorer interface {
	StoreMatchEvent(ctx context.Context, event Event, clockSnap ClockData)
	StoreGameEndedEvent(ctx context.Context, result *chess.Result)
}

type MatchStore struct {
	queries *sqlc.Queries
	wg      sync.WaitGroup
}

type MatchEventStore struct {
	store      *MatchStore
	internalID int64
	seqNum     int32
}

func NewMatchStore(queries *sqlc.Queries) *MatchStore {
	return &MatchStore{
		queries: queries,
	}
}

func (s *MatchStore) CreateMatch(ctx context.Context, publicID string, player1, player2 *Player, timeFormat TimeFormat) (MatchEventStorer, error) {
	if err := timeFormat.Validate(); err != nil {
		return nil, fmt.Errorf("invalid time format: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
	defer cancel()

	whitePlayer := player1
	blackPlayer := player2
	if player1.GetColor() != chess.White {
		whitePlayer = player2
		blackPlayer = player1
	}

	var whiteID, blackID pgtype.Int8
	if whitePlayer.Username != "" {
		u, err := s.queries.GetUserByUsername(ctx, whitePlayer.Username)
		if err == nil {
			whiteID = pgtype.Int8{Int64: u.ID, Valid: true}
		}
	}
	if blackPlayer.Username != "" {
		u, err := s.queries.GetUserByUsername(ctx, blackPlayer.Username)
		if err == nil {
			blackID = pgtype.Int8{Int64: u.ID, Valid: true}
		}
	}

	createdMatch, err := s.queries.CreateMatch(ctx, sqlc.CreateMatchParams{
		PublicID:         publicID,
		WhiteUserID:      whiteID,
		BlackUserID:      blackID,
		WhiteDisplayName: whitePlayer.DisplayName,
		BlackDisplayName: blackPlayer.DisplayName,
		InitialTimeMs:    timeFormat.initial.Milliseconds(),
		IncrementMs:      timeFormat.increment.Milliseconds(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to insert match: %w", err)
	}

	return &MatchEventStore{
		store:      s,
		internalID: createdMatch.ID,
	}, nil
}

func (s *MatchStore) Close() {
	s.wg.Wait()
}

func (ms *MatchEventStore) storeEvent(ctx context.Context, eventType sqlc.MatchEventType, payload []byte) {
	if len(payload) == 0 {
		payload = []byte("{}")
	}

	ms.seqNum++
	ms.store.wg.Add(1)
	go func(matchID int64, seqNum int32, p []byte) {
		defer ms.store.wg.Done()
		ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
		defer cancel()

		_, err := ms.store.queries.InsertMatchEvent(ctx, sqlc.InsertMatchEventParams{
			MatchID:   matchID,
			SeqNum:    seqNum,
			EventType: eventType,
			Payload:   p,
		})
		if err != nil {
			log.Printf("ERROR [MatchStore]: failed to store event %s (seq=%d) for match %d: %v", eventType, seqNum, matchID, err)
		}
	}(ms.internalID, ms.seqNum, payload)
}

func (ms *MatchEventStore) StoreMatchEvent(ctx context.Context, event Event, clockSnap ClockData) {
	switch event.Type {
	case EventTypeGameStarted:
		ms.storeEvent(ctx, sqlc.MatchEventTypeGameStarted, nil)

	case EventTypeMove:
		data, ok := event.Data.(MoveData)
		if !ok {
			return
		}
		payload, err := json.Marshal(MovePayload{
			From:        data.From,
			To:          data.To,
			Promotion:   data.Promotion,
			WhiteTimeMs: clockSnap.WhiteTimeMs,
			BlackTimeMs: clockSnap.BlackTimeMs,
		})
		if err != nil {
			log.Printf("ERROR [MatchStore]: failed to marshal move payload: %v", err)
			return
		}
		ms.storeEvent(ctx, sqlc.MatchEventTypeMove, payload)
	}
}

func (ms *MatchEventStore) StoreGameEndedEvent(ctx context.Context, result *chess.Result) {
	if result == nil {
		return
	}
	payload, err := json.Marshal(GameEndedPayload{
		Outcome: result.Outcome,
		Reason:  result.Reason,
	})
	if err != nil {
		log.Printf("ERROR [MatchStore]: failed to marshal game_ended payload: %v", err)
		return
	}
	ms.storeEvent(ctx, sqlc.MatchEventTypeGameEnded, payload)
}
