package chess

import (
	"encoding/json"
	"fmt"
)

type Player struct {
	Name     string
	Match    *Match
	SendChan chan WSMessage
}

type Match struct {
	White Player
	Black Player

	Moves     []Move
	Position  *Position
	EventChan chan Event
}

func NewInitialPosition() *Position {
	pos, err := StartingPositionFEN.ToPosition()
	if err != nil {
		panic(fmt.Sprintf("failed to create initial position: %v", err))
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
			err := m.Position.ApplyMove(moveData.Move)
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

func (m *Match) sendError(message string) error {
	errorData := ErrorData{Message: message}
	errorDataJson, err := json.Marshal(errorData)
	if err != nil {
		return fmt.Errorf("failed to marshal error data: %v", err)
	}
	m.White.SendChan <- WSMessage{
		Type: MessageTypeError,
		Data: errorDataJson,
	}
	m.Black.SendChan <- WSMessage{
		Type: MessageTypeError,
		Data: errorDataJson,
	}
	return nil
}
