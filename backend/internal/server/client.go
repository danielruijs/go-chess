package server

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn   *websocket.Conn
	Player *Player
	Done   chan struct{}
}

func NewClient(conn *websocket.Conn) *Client {
	done := make(chan struct{})
	player := NewPlayer("", done)

	return &Client{
		Conn:   conn,
		Player: player,
		Done:   done,
	}
}

func (c *Client) ReceiveMessages(matchmaker *Matchmaker) {
	for {
		var message WSMessage
		err := c.Conn.ReadJSON(&message)
		if err != nil {
			expectedCodes := []int{websocket.CloseNormalClosure, websocket.CloseGoingAway}
			if websocket.IsUnexpectedCloseError(err, expectedCodes...) {
				log.Println("WebSocket error:", err)
			}
			return
		}

		handler, ok := messageHandlers[message.Type]
		if !ok {
			log.Println("Unsupported message type:", message.Type)
			continue
		}

		err = handler.Handle(c, message.Data, matchmaker)
		if err != nil {
			log.Printf("Error handling message with type %v: %v\n", message.Type, err)
		}
	}
}

func (c *Client) SendMessages() {
	ch := c.Player.GetSendChannel()
	for {
		select {
		case message := <-ch:
			err := c.Conn.WriteJSON(message)
			if err != nil {
				log.Println("Error sending message:", err)
			}
		case <-c.Done:
			return
		}
	}
}
