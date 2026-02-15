package server

import (
	"encoding/json"
	"fmt"
	"go-chess/internal/chess"
	"go-chess/internal/matchmaker"
)

func (c Client) handleMove(message chess.WSMessage) error {
	var data chess.MoveData
	if err := json.Unmarshal(message.Data, &data); err != nil {
		return fmt.Errorf("invalid move data format: %w", err)
	}

	if c.Player.Match == nil {
		return fmt.Errorf("player is not in a match")
	}

	c.Player.Match.EventChan <- chess.Event{
		Player: c.Player,
		Type:   chess.EventTypeMove,
		Data:   data,
	}

	return nil
}

func (c Client) handleJoinMatch(message chess.WSMessage, matchmaker *matchmaker.Matchmaker) error {
	var data chess.JoinMatchData
	if err := json.Unmarshal(message.Data, &data); err != nil {
		return fmt.Errorf("invalid join match data format: %w", err)
	}

	if c.Player.Match != nil {
		return fmt.Errorf("player is already in a match")
	}

	c.Player.Name = data.PlayerName
	matchmaker.JoinQueue(c.Player)

	return nil
}
