package server

import (
	"encoding/json"
	"fmt"
	"go-chess/internal/chess"
	"log"
)

type Player struct {
	Name     string
	Match    *Match
	SendChan chan WSMessage
	Color    chess.Color
}

type Match struct {
	White *Player
	Black *Player

	Moves     []chess.Move
	Position  *chess.Position
	EventChan chan Event
}

func NewInitialPosition() *chess.Position {
	pos, err := chess.StartingPositionFEN.ToPosition()
	if err != nil {
		log.Fatalf("failed to create initial position: %v", err)
	}
	return &pos
}

func (m *Match) Run() {
	for event := range m.EventChan {
		switch event.Type {
		case EventTypeMove:
			moveData, ok := event.Data.(MoveData)
			if !ok {
				fmt.Println("invalid move data format")
				m.sendError("invalid move data format")
				continue
			}
			err := m.Position.ApplyMove(moveData.Move, event.Player.Color)
			if err != nil {
				fmt.Println("invalid move:", err)
				m.sendError(fmt.Sprintf("invalid move: %v", err))
				continue
			}
			m.Moves = append(m.Moves, moveData.Move)
		case EventTypeGameStarted:
			fmt.Println("Started match")
		}
		err := m.sendPositionUpdate()
		if err != nil {
			fmt.Println("failed to send position update:", err)
		}
	}
}

func (m *Match) sendPositionUpdate() error {
	positionData := PositionData{Position: *m.Position}
	data, err := json.Marshal(positionData)
	if err != nil {
		return fmt.Errorf("failed to marshal position: %v", err)
	}
	m.White.SendChan <- WSMessage{
		Type: MessageTypePosition,
		Data: data,
	}
	m.Black.SendChan <- WSMessage{
		Type: MessageTypePosition,
		Data: data,
	}
	return nil
}

func (m *Match) sendError(message string) {
	errorData := ErrorData{Message: message}
	errorDataJson, err := json.Marshal(errorData)
	if err != nil {
		log.Fatal("failed to marshal error data:", err)
	}
	m.White.SendChan <- WSMessage{
		Type: MessageTypeError,
		Data: errorDataJson,
	}
	m.Black.SendChan <- WSMessage{
		Type: MessageTypeError,
		Data: errorDataJson,
	}
}
