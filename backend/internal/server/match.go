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
	Player1 *Player
	Player2 *Player

	Engine    *chess.Engine
	EventChan chan Event
}

func (m *Match) Run() {
	for event := range m.EventChan {
		switch event.Type {
		case EventTypeMove:
			data, ok := event.Data.(MoveData)
			if !ok {
				log.Println("invalid move data format")
				continue
			}
			move, err := moveDataToMove(data)
			if err != nil {
				log.Println("invalid move data:", err)
				continue
			}
			result, err := m.Engine.ApplyMove(move, event.Player.Color)
			if err != nil {
				log.Printf("failed to apply move %s -> %s, promotion to %v: %v\n", data.From, data.To, *data.Promotion, err)
				continue
			}
			if result != nil {
				err := m.sendFinalPositionUpdate()
				if err != nil {
					log.Println("failed to send final position update:", err)
				}
				m.sendMatchEnd(*result)
				close(m.EventChan)
				return
			}
		case EventTypeGameStarted:
			for _, player := range []*Player{m.Player1, m.Player2} {
				startMatchData := StartMatchData{
					WhitePlayerName: m.getPlayerByColor(chess.White).Name,
					BlackPlayerName: m.getPlayerByColor(chess.Black).Name,
				}
				data, err := json.Marshal(startMatchData)
				if err != nil {
					log.Printf("failed to marshal start match data: %v", err)
					continue
				}
				player.SendChan <- WSMessage{
					Type: MessageTypeStartMatch,
					Data: data,
				}
			}
		}
		err := m.sendPositionUpdate()
		if err != nil {
			log.Println("failed to send position update:", err)
		}
	}
}

func (m *Match) sendPositionUpdate() error {
	for _, player := range []*Player{m.Player1, m.Player2} {
		legalMovesList := m.Engine.GetLegalMoves(player.Color)
		boardData := BoardData{
			Board:      m.Engine.GetBoard(),
			LegalMoves: moveListToLegalMoves(legalMovesList),
			PGN:        m.Engine.GetPGN(),
		}
		data, err := json.Marshal(boardData)
		if err != nil {
			return fmt.Errorf("failed to marshal board data: %v", err)
		}
		player.SendChan <- WSMessage{
			Type: MessageTypeBoard,
			Data: data,
		}
	}
	return nil
}

func (m *Match) sendFinalPositionUpdate() error {
	for _, player := range []*Player{m.Player1, m.Player2} {
		boardData := BoardData{
			Board:      m.Engine.GetBoard(),
			LegalMoves: map[string][]LegalMove{}, // no legal moves
			PGN:        m.Engine.GetPGN(),
		}
		data, err := json.Marshal(boardData)
		if err != nil {
			return fmt.Errorf("failed to marshal board data: %v", err)
		}
		player.SendChan <- WSMessage{
			Type: MessageTypeBoard,
			Data: data,
		}
	}
	return nil
}

func (m *Match) sendMatchEnd(result chess.Result) {
	resultData := EndMatchData{
		Result: result,
	}
	data, err := json.Marshal(resultData)
	if err != nil {
		log.Printf("failed to marshal end match data: %v", err)
		return
	}
	for _, player := range []*Player{m.Player1, m.Player2} {
		player.SendChan <- WSMessage{
			Type: MessageTypeEndMatch,
			Data: data,
		}
	}

}

func (m *Match) getPlayerByColor(color chess.Color) *Player {
	if m.Player1.Color == color {
		return m.Player1
	}
	return m.Player2
}
