package server

import (
	"errors"
	"go-chess/internal/chess"
	"log"
	"math/rand"
)

type Matchmaker struct {
	queue      map[*Player]struct{}
	actions    chan func()
	UpdateChan chan struct{}
}

func NewMatchmaker() *Matchmaker {
	return &Matchmaker{
		queue:      make(map[*Player]struct{}),
		actions:    make(chan func()),
		UpdateChan: make(chan struct{}),
	}
}

func (mm *Matchmaker) Join(player *Player) error {
	errChan := make(chan error)
	mm.actions <- func() {
		if _, ok := mm.queue[player]; ok {
			errChan <- errors.New("player is already in the matchmaking queue")
			return
		}

		mm.queue[player] = struct{}{}
		log.Printf("Player %s joined. Queue size: %d\n", player.Name, len(mm.queue))

		errChan <- nil
	}
	return <-errChan
}

func (mm *Matchmaker) Leave(player *Player) {
	mm.actions <- func() {
		if _, ok := mm.queue[player]; !ok {
			return
		}
		delete(mm.queue, player)

		log.Printf("Player %s left. Queue size: %d\n", player.Name, len(mm.queue))
	}
}

func (mm *Matchmaker) Run() {
	for action := range mm.actions {
		action()

		select {
		case mm.UpdateChan <- struct{}{}:
		default:
		}

		for len(mm.queue) >= 2 {
			var p1, p2 *Player
			for p := range mm.queue {
				if p1 == nil {
					p1 = p
				} else {
					p2 = p
					break
				}
			}
			delete(mm.queue, p1)
			delete(mm.queue, p2)

			startMatch(p1, p2)
		}
	}
}

func startMatch(player1, player2 *Player) {
	match := &Match{
		Player1:   player1,
		Player2:   player2,
		Engine:    chess.NewEngine(),
		EventChan: make(chan Event),
	}

	// Randomly assign colors
	if rand.Intn(2) == 0 {
		player1.Color, player2.Color = chess.White, chess.Black
	} else {
		player1.Color, player2.Color = chess.Black, chess.White
	}

	player1.Match = match
	player2.Match = match

	log.Printf("Starting match: %s (color: %s) vs %s (color: %s)\n", player1.Name, player1.Color, player2.Name, player2.Color)

	go match.Run()

	match.EventChan <- Event{Type: EventTypeGameStarted}
}
