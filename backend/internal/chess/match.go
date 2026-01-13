package chess

import "fmt"

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
			move, ok := event.Data.(Move)
			if !ok {
				fmt.Println("invalid move data format")
				continue
			}
			err := m.Position.ApplyMove(move)
			if err != nil {
				fmt.Println("invalid move:", err)
				continue
			}
			m.Moves = append(m.Moves, move)

			m.White.SendChan <- WSMessage{
				Type: MessageTypePosition,
				Data: *m.Position,
			}
			m.Black.SendChan <- WSMessage{
				Type: MessageTypePosition,
				Data: *m.Position,
			}
		case EventTypeGameStarted:
			fmt.Println("Started match")
			m.White.SendChan <- WSMessage{
				Type: MessageTypePosition,
				Data: *m.Position,
			}
			m.Black.SendChan <- WSMessage{
				Type: MessageTypePosition,
				Data: *m.Position,
			}
		}
		fmt.Println(event.Player, event.Type)
	}
}
