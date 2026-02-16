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

	Engine    chess.Engine
	EventChan chan Event
}

func (m *Match) Run() {
	for event := range m.EventChan {
		switch event.Type {
		case EventTypeMove:
			data, ok := event.Data.(MoveData)
			if !ok {
				fmt.Println("invalid move data format")
				m.sendError("invalid move data format")
				continue
			}
			move, err := moveDataToMove(data)
			if err != nil {
				fmt.Println("invalid move data:", err)
				m.sendError(fmt.Sprintf("invalid move data: %v", err))
				continue
			}
			err = m.Engine.ApplyMove(move, event.Player.Color)
			if err != nil {
				fmt.Printf("failed to apply move %s -> %s: %v\n", data.From, data.To, err)
				m.sendError(fmt.Sprintf("invalid move: %v", err))
				continue
			}
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
	positionData := BoardData{Board: m.Engine.GetBoard()}
	data, err := json.Marshal(positionData)
	if err != nil {
		return fmt.Errorf("failed to marshal position: %v", err)
	}
	m.White.SendChan <- WSMessage{
		Type: MessageTypeBoard,
		Data: data,
	}
	m.Black.SendChan <- WSMessage{
		Type: MessageTypeBoard,
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
