package server

import (
	"errors"
	"go-chess/internal/chess"
	"log"
	"math/rand"
	"sync"
)

type Matchmaker struct {
	queue      map[TimeFormat]map[*Player]struct{}
	actions    chan func()
	UpdateChan chan struct{}
	matchEnded chan *Match

	activeMatches   map[*Player]*Match
	activeMatchesMu sync.RWMutex
}

func NewMatchmaker() *Matchmaker {
	return &Matchmaker{
		queue:         make(map[TimeFormat]map[*Player]struct{}),
		actions:       make(chan func()),
		UpdateChan:    make(chan struct{}),
		matchEnded:    make(chan *Match),
		activeMatches: make(map[*Player]*Match),
	}
}

func (mm *Matchmaker) GetMatch(player *Player) *Match {
	mm.activeMatchesMu.RLock()
	defer mm.activeMatchesMu.RUnlock()
	return mm.activeMatches[player]
}

func (mm *Matchmaker) RegisterMatch(match *Match) {
	mm.activeMatchesMu.Lock()
	defer mm.activeMatchesMu.Unlock()
	mm.activeMatches[match.Player1] = match
	mm.activeMatches[match.Player2] = match
}

func (mm *Matchmaker) UnregisterMatch(match *Match) {
	mm.activeMatchesMu.Lock()
	defer mm.activeMatchesMu.Unlock()
	delete(mm.activeMatches, match.Player1)
	delete(mm.activeMatches, match.Player2)
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
		log.Printf("Player %s joined %v. Queue size: %d\n", player.Name, timeFormat, len(mm.queue[timeFormat]))

		errChan <- nil
	}
	return <-errChan
}

func (mm *Matchmaker) Leave(player *Player) {
	mm.actions <- func() {
		player.LeaveQueues()
		for timeFormat := range mm.queue {
			if _, ok := mm.queue[timeFormat][player]; ok {
				delete(mm.queue[timeFormat], player)
				log.Printf("Player %s left %v. Queue size: %d\n", player.Name, timeFormat, len(mm.queue[timeFormat]))
			}
		}
	}
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

func (mm *Matchmaker) Run() {
	for {
		select {
		case action := <-mm.actions:
			action()

			select {
			case mm.UpdateChan <- struct{}{}:
			default:
			}

			for timeFormat := range mm.queue {
				if len(mm.queue[timeFormat]) < 2 {
					continue
				}
				p1, p2 := mm.matchPlayers(timeFormat)

				mm.startMatch(p1, p2, timeFormat)
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
	for timeFormat := range mm.queue {
		delete(mm.queue[timeFormat], p1)
		delete(mm.queue[timeFormat], p2)
	}
	p1.LeaveQueues()
	p2.LeaveQueues()
	return p1, p2
}

func (mm *Matchmaker) startMatch(player1, player2 *Player, timeFormat TimeFormat) {
	match := &Match{
		Player1:    player1,
		Player2:    player2,
		Engine:     chess.NewEngine(),
		Clock:      NewMatchClock(timeFormat),
		EventChan:  make(chan Event),
		MatchEnded: mm.matchEnded,
	}

	// Randomly assign colors
	if rand.Intn(2) == 0 {
		player1.SetColor(chess.White)
		player2.SetColor(chess.Black)
	} else {
		player1.SetColor(chess.Black)
		player2.SetColor(chess.White)
	}

	mm.RegisterMatch(match)

	log.Printf("Starting match: %s (color: %s) vs %s (color: %s)\n", player1.Name, player1.GetColor(), player2.Name, player2.GetColor())

	go match.Run()

	match.EventChan <- Event{Type: EventTypeGameStarted}
}
