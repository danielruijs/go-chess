package server

import (
	"encoding/json"
	"fmt"
)

func (c Client) handleJoinMatch(messageData json.RawMessage, matchmaker *Matchmaker) error {
	var data JoinMatchData
	if err := json.Unmarshal(messageData, &data); err != nil {
		return fmt.Errorf("invalid join match data format: %w", err)
	}

	if c.Player.Match != nil {
		return fmt.Errorf("player is already in a match")
	}

	c.Player.Name = data.PlayerName
	err := matchmaker.Join(c.Player)
	if err != nil {
		return fmt.Errorf("failed to join matchmaking queue: %w", err)
	}

	return nil
}

func (c Client) handleMove(messageData json.RawMessage) error {
	var data MoveData
	if err := json.Unmarshal(messageData, &data); err != nil {
		return fmt.Errorf("invalid move data format: %w", err)
	}

	if c.Player.Match == nil {
		return fmt.Errorf("player is not in a match")
	}

	c.Player.Match.EventChan <- Event{
		Player: c.Player,
		Type:   EventTypeMove,
		Data:   data,
	}

	return nil
}

func (c Client) handleResign() error {
	if c.Player.Match == nil {
		return fmt.Errorf("player is not in a match")
	}

	c.Player.Match.EventChan <- Event{
		Player: c.Player,
		Type:   EventTypeResign,
	}

	return nil
}

func (c Client) handleOfferDraw() error {
	if c.Player.Match == nil {
		return fmt.Errorf("player is not in a match")
	}

	c.Player.Match.EventChan <- Event{
		Player: c.Player,
		Type:   EventTypeOfferDraw,
	}

	return nil
}

func (c Client) handleRespondDraw(messageData json.RawMessage) error {
	var data RespondDrawData
	if err := json.Unmarshal(messageData, &data); err != nil {
		return fmt.Errorf("invalid respond draw data format: %w", err)
	}

	if c.Player.Match == nil {
		return fmt.Errorf("player is not in a match")
	}

	c.Player.Match.EventChan <- Event{
		Player: c.Player,
		Type:   EventTypeRespondDraw,
		Data:   data,
	}

	return nil
}
