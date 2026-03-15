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

	if matchmaker.GetMatch(c.Player) != nil {
		return fmt.Errorf("player is already in a match")
	}

	c.Player.Name = data.PlayerName
	timeFormat := MsToTimeFormat(data.TimeFormat)
	err := matchmaker.Join(c.Player, timeFormat)
	if err != nil {
		return fmt.Errorf("failed to join matchmaking queue: %w", err)
	}

	return nil
}

func (c Client) handleMove(messageData json.RawMessage, matchmaker *Matchmaker) error {
	var data MoveData
	if err := json.Unmarshal(messageData, &data); err != nil {
		return fmt.Errorf("invalid move data format: %w", err)
	}

	match := matchmaker.GetMatch(c.Player)
	if match == nil {
		return fmt.Errorf("player is not in a match")
	}

	match.EventChan <- Event{
		Player: c.Player,
		Type:   EventTypeMove,
		Data:   data,
	}

	return nil
}

func (c Client) handleResign(matchmaker *Matchmaker) error {
	match := matchmaker.GetMatch(c.Player)
	if match == nil {
		return fmt.Errorf("player is not in a match")
	}

	match.EventChan <- Event{
		Player: c.Player,
		Type:   EventTypeResign,
	}

	return nil
}

func (c Client) handleOfferDraw(matchmaker *Matchmaker) error {
	match := matchmaker.GetMatch(c.Player)
	if match == nil {
		return fmt.Errorf("player is not in a match")
	}

	match.EventChan <- Event{
		Player: c.Player,
		Type:   EventTypeOfferDraw,
	}

	return nil
}

func (c Client) handleRespondDraw(messageData json.RawMessage, matchmaker *Matchmaker) error {
	var data RespondDrawData
	if err := json.Unmarshal(messageData, &data); err != nil {
		return fmt.Errorf("invalid respond draw data format: %w", err)
	}

	match := matchmaker.GetMatch(c.Player)
	if match == nil {
		return fmt.Errorf("player is not in a match")
	}

	match.EventChan <- Event{
		Player: c.Player,
		Type:   EventTypeRespondDraw,
		Data:   data,
	}

	return nil
}
