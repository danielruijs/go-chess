package server

import (
	"fmt"
	"go-chess/internal/chess"
)

type Matchmaker struct {
	Queue chan *Player
}

func NewMatchmaker() *Matchmaker {
	return &Matchmaker{
		Queue: make(chan *Player),
	}
}

func (mm *Matchmaker) JoinQueue(player *Player) {
	mm.Queue <- player
}

func (mm *Matchmaker) Run() {
	for {
		player1 := <-mm.Queue
		fmt.Println("Found player 1", player1.Name)
		player2 := <-mm.Queue
		fmt.Println("Found player 2", player2.Name)

		match := &Match{
			White:     player1,
			Black:     player2,
			Engine:    chess.NewEngine(),
			EventChan: make(chan Event),
		}

		player1.Color = chess.White
		player2.Color = chess.Black

		player1.Match = match
		player2.Match = match

		go match.Run()

		match.EventChan <- Event{
			Player: nil,
			Type:   EventTypeGameStarted,
			Data:   nil,
		}
	}
}
