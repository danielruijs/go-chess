package matchmaker

import (
	"fmt"
	"go-chess/internal/chess"
)

type Matchmaker struct {
	Queue chan *chess.Player
}

func NewMatchmaker() *Matchmaker {
	return &Matchmaker{
		Queue: make(chan *chess.Player),
	}
}

func (mm *Matchmaker) JoinQueue(player *chess.Player) {
	mm.Queue <- player
}

func (mm *Matchmaker) Run() {
	for {
		player1 := <-mm.Queue
		fmt.Println("Found player 1", player1.Name)
		player2 := <-mm.Queue
		fmt.Println("Found player 2", player2.Name)

		match := &chess.Match{
			White:     player1,
			Black:     player2,
			Moves:     []chess.Move{},
			Position:  chess.NewInitialPosition(),
			EventChan: make(chan chess.Event),
		}

		player1.Match = match
		player2.Match = match

		go match.Run()

		match.EventChan <- chess.Event{
			Player: nil,
			Type:   chess.EventTypeGameStarted,
			Data:   nil,
		}
	}
}
