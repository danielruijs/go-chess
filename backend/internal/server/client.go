package server

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn   *websocket.Conn
	Player *Player

	Done      chan struct{}
	closeOnce sync.Once
	sendChan  chan WSMessage

	metrics *metrics
}

func (c *Client) Close() {
	c.closeOnce.Do(func() { close(c.Done) })
}

func NewClient(conn *websocket.Conn, player *Player, metrics *metrics) *Client {
	return &Client{
		Conn:     conn,
		Player:   player,
		Done:     make(chan struct{}),
		sendChan: make(chan WSMessage, 100),
		metrics:  metrics,
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

		c.metrics.recordWebsocketMessageReceived(message.Type)

		handler, ok := messageHandlers[message.Type]
		if !ok {
			c.metrics.recordWebsocketMessageError(message.Type, "unsupported_type")
			log.Println("Unsupported message type:", message.Type)
			continue
		}

		err = handler.Handle(c, message.Data, matchmaker)
		if err != nil {
			c.metrics.recordWebsocketMessageError(message.Type, "handler_error")
			log.Printf("Error handling message with type %v: %v\n", message.Type, err)
		}
	}
}

func (c *Client) SendMessages() {
	for {
		select {
		case message := <-c.sendChan:
			err := c.Conn.WriteJSON(message)
			if err != nil {
				c.metrics.recordWebsocketMessageSendError(message.Type, "write_error")
				log.Println("Error sending message:", err)
			} else {
				c.metrics.recordWebsocketMessageSent(message.Type)
			}
		case <-c.Done:
			return
		}
	}
}
