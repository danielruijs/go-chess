package server

import (
	"context"
	"errors"
	"go-chess/internal/cache"
	"go-chess/internal/chess"
	"log"
	"math/rand"
)

type Matchmaker struct {
	queue      map[TimeFormat]map[*Player]struct{}
	actions    chan func()
	UpdateChan chan struct{}
	matchEnded chan *Match

	activeMatches *cache.Cache[*Player, *Match]

	metrics     *metrics
	matchStorer MatchStorer
}

func NewMatchmaker(metrics *metrics, matchStorer MatchStorer) (*Matchmaker, error) {
	activeMatches, err := cache.New[*Player](cache.Options[*Match]{})
	if err != nil {
		return nil, err
	}
	return &Matchmaker{
		queue:         make(map[TimeFormat]map[*Player]struct{}),
		actions:       make(chan func()),
		UpdateChan:    make(chan struct{}),
		matchEnded:    make(chan *Match),
		activeMatches: activeMatches,
		metrics:       metrics,
		matchStorer:   matchStorer,
	}, nil
}

func (mm *Matchmaker) GetMatch(player *Player) *Match {
	match, _ := mm.activeMatches.Get(player)
	return match
}

func (mm *Matchmaker) RegisterMatch(match *Match) {
	mm.activeMatches.Set(match.Player1, match)
	mm.activeMatches.Set(match.Player2, match)
}

func (mm *Matchmaker) UnregisterMatch(match *Match) {
	mm.activeMatches.Delete(match.Player1)
	mm.activeMatches.Delete(match.Player2)
}

func (mm *Matchmaker) Join(player *Player, timeFormat TimeFormat) error {
	errChan := make(chan error, 1)
	mm.actions <- func() {
		if _, ok := mm.queue[timeFormat][player]; ok {
			errChan <- errors.New("player is already in the matchmaking queue")
			return
		}

		if _, ok := mm.queue[timeFormat]; !ok {
			mm.queue[timeFormat] = make(map[*Player]struct{})
		}
		mm.queue[timeFormat][player] = struct{}{}
		player.JoinQueue(timeFormat)
		mm.metrics.recordQueueJoin(timeFormat, len(mm.queue[timeFormat]))
		log.Printf("Player %s joined %v. Queue size: %d\n", player.DisplayName, timeFormat, len(mm.queue[timeFormat]))

		errChan <- nil
	}
	return <-errChan
}

func (mm *Matchmaker) removePlayerFromQueue(player *Player, timeFormat TimeFormat) {
	if _, exists := mm.queue[timeFormat][player]; exists {
		delete(mm.queue[timeFormat], player)
		mm.metrics.recordQueueLeave(timeFormat, len(mm.queue[timeFormat]))
		log.Printf("Player %s left %v. Queue size: %d\n", player.DisplayName, timeFormat, len(mm.queue[timeFormat]))
	}
}

func (mm *Matchmaker) removePlayersFromQueues(players ...*Player) {
	for _, player := range players {
		for _, timeFormat := range player.GetQueues() {
			mm.removePlayerFromQueue(player, timeFormat)
		}
		player.LeaveQueues()
	}
}

func (mm *Matchmaker) Leave(player *Player, timeFormat TimeFormat) {
	done := make(chan struct{})
	mm.actions <- func() {
		defer close(done)
		mm.removePlayerFromQueue(player, timeFormat)
		player.LeaveQueue(timeFormat)
	}
	<-done
}

func (mm *Matchmaker) LeaveAll(player *Player) {
	done := make(chan struct{})
	mm.actions <- func() {
		defer close(done)
		mm.removePlayersFromQueues(player)
	}
	<-done
}

// Returns a snapshot of queue lengths for each time format
func (mm *Matchmaker) GetQueueStats() map[TimeFormat]int {
	update := make(chan map[TimeFormat]int, 1)
	mm.actions <- func() {
		queueStats := make(map[TimeFormat]int)
		for timeFormat, players := range mm.queue {
			queueStats[timeFormat] = len(players)
		}
		update <- queueStats
	}
	return <-update
}

func (mm *Matchmaker) Run(ctx context.Context) {
	defer close(mm.UpdateChan)
	for {
		select {
		case <-ctx.Done():
			return
		case action := <-mm.actions:
			action()

			for timeFormat, players := range mm.queue {
				if len(players) < 2 {
					continue
				}
				p1, p2 := mm.matchPlayers(timeFormat)

				mm.startMatch(ctx, p1, p2, timeFormat)
			}

			select {
			case mm.UpdateChan <- struct{}{}:
			default:
			}
		case match := <-mm.matchEnded:
			mm.UnregisterMatch(match)
		}
	}
}

func (mm *Matchmaker) matchPlayers(timeFormat TimeFormat) (*Player, *Player) {
	var p1, p2 *Player
	for p := range mm.queue[timeFormat] {
		if p1 == nil {
			p1 = p
		} else {
			p2 = p
			break
		}
	}

	if p1 == nil || p2 == nil {
		return nil, nil
	}

	mm.removePlayersFromQueues(p1, p2)
	return p1, p2
}

func (mm *Matchmaker) startMatch(ctx context.Context, player1, player2 *Player, timeFormat TimeFormat) {
	// Randomly assign colors
	if rand.Intn(2) == 0 {
		player1.SetColor(chess.White)
		player2.SetColor(chess.Black)
	} else {
		player1.SetColor(chess.Black)
		player2.SetColor(chess.White)
	}

	match, err := NewMatch(ctx, mm.matchStorer, player1, player2, timeFormat, mm.matchEnded, mm.metrics)
	if err != nil {
		log.Printf("ERROR [startMatch]: failed to create match: %v", err)
		return
	}

	mm.RegisterMatch(match)
	mm.metrics.recordMatchStarted(timeFormat)

	log.Printf("Starting match %s: %s (color: %s) vs %s (color: %s)\n", match.PublicID, player1.DisplayName, player1.GetColor(), player2.DisplayName, player2.GetColor())

	go match.Run(ctx)

	match.EventChan <- Event{Type: EventTypeGameStarted}
}
