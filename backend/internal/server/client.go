package server

import (
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	websocketWriteTimeout = 10 * time.Second
	websocketPongTimeout  = 60 * time.Second
	websocketPingInterval = 30 * time.Second
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
	c.closeOnce.Do(func() {
		close(c.Done)
		if c.Conn != nil {
			_ = c.Conn.Close()
		}
	})
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
	// Set read deadline and extend it on every incoming pong keepalive frame.
	c.Conn.SetReadDeadline(time.Now().Add(websocketPongTimeout))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(websocketPongTimeout))
		return nil
	})

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

		c.Conn.SetReadDeadline(time.Now().Add(websocketPongTimeout))
		c.metrics.recordWebsocketMessageReceived(message.Type)

		handler, ok := messageHandlers[message.Type]
		if !ok {
			c.metrics.recordWebsocketMessageReceiveError(message.Type)
			log.Println("Unsupported message type:", message.Type)
			continue
		}

		err = handler.Handle(c, message.Data, matchmaker)
		if err != nil {
			c.metrics.recordWebsocketMessageReceiveError(message.Type)
			log.Printf("Error handling message with type %v: %v\n", message.Type, err)
		}
	}
}

func (c *Client) SendMessages() {
	// Send periodic pings to keep intermediate proxies from timing out idle connections.
	pingTicker := time.NewTicker(websocketPingInterval)
	defer pingTicker.Stop()

	for {
		select {
		case message := <-c.sendChan:
			c.Conn.SetWriteDeadline(time.Now().Add(websocketWriteTimeout))
			err := c.Conn.WriteJSON(message)
			if err != nil {
				c.metrics.recordWebsocketMessageSendError(message.Type)
				log.Println("Error sending message:", err)
				c.Close()
				return
			}
			c.metrics.recordWebsocketMessageSent(message.Type)
		case <-pingTicker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(websocketWriteTimeout))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Println("Error sending ping:", err)
				c.Close()
				return
			}
		case <-c.Done:
			return
		}
	}
}
