package server

import (
	"encoding/json"
	"fmt"
)

type MessageHandler interface {
	Handle(c *Client, messageData json.RawMessage, matchmaker *Matchmaker) error
}

var messageHandlers = map[MessageType]MessageHandler{
	MessageTypeJoinMatch:   JoinMatchMessageHandler{},
	MessageTypeMove:        MoveMessageHandler{},
	MessageTypeResign:      ResignMessageHandler{},
	MessageTypeOfferDraw:   OfferDrawMessageHandler{},
	MessageTypeRespondDraw: RespondDrawMessageHandler{},
}

type JoinMatchMessageHandler struct{}

func (h JoinMatchMessageHandler) Handle(c *Client, messageData json.RawMessage, matchmaker *Matchmaker) error {
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

type MoveMessageHandler struct{}

func (h MoveMessageHandler) Handle(c *Client, messageData json.RawMessage, matchmaker *Matchmaker) error {
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

type ResignMessageHandler struct{}

func (h ResignMessageHandler) Handle(c *Client, messageData json.RawMessage, matchmaker *Matchmaker) error {
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

type OfferDrawMessageHandler struct{}

func (h OfferDrawMessageHandler) Handle(c *Client, messageData json.RawMessage, matchmaker *Matchmaker) error {
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

type RespondDrawMessageHandler struct{}

func (h RespondDrawMessageHandler) Handle(c *Client, messageData json.RawMessage, matchmaker *Matchmaker) error {
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
