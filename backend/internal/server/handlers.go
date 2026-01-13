package server

import (
	"fmt"
	"go-chess/internal/chess"
	"go-chess/internal/matchmaker"
)

func (c Client) handleMove(message chess.WSMessage) error {
	if c.Player.Match == nil {
		return fmt.Errorf("player is not in a match")
	}

	move, ok := message.Data.(chess.Move)
	if !ok {
		return fmt.Errorf("invalid move data format")
	}

	c.Player.Match.EventChan <- chess.Event{
		Player: c.Player,
		Type:   chess.EventTypeMove,
		Data:   move,
	}

	return nil
}

func (c Client) handleJoinMatch(matchmaker *matchmaker.Matchmaker) error {
	if c.Player.Match != nil {
		return fmt.Errorf("player is already in a match")
	}

	matchmaker.JoinQueue(c.Player)

	return nil
}
