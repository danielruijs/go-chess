package server

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type WebSocketHandler struct {
	Upgrader   websocket.Upgrader
	Matchmaker *Matchmaker
}

type Client struct {
	Conn   *websocket.Conn
	Player *Player
}

func (wsh WebSocketHandler) ServeWS(w http.ResponseWriter, r *http.Request) {
	wsh.Upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	conn, err := wsh.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error when upgrading connection to websocket: %s", err)
		return
	}

	client := Client{
		Conn: conn,
		Player: &Player{
			Name:     "",
			SendChan: make(chan WSMessage),
		},
	}

	go client.ReceiveMessages(wsh.Matchmaker)
	go client.SendMessages()

	log.Println("New WebSocket connection established")
}

func (c Client) ReceiveMessages(matchmaker *Matchmaker) {
	defer func() {
		matchmaker.Leave(c.Player)
		_ = c.Conn.Close()
	}()
	for {
		var message WSMessage
		err := c.Conn.ReadJSON(&message)
		if err != nil {
			log.Println("Failed to unmarshal message:", err)
			return
		}

		switch message.Type {
		case MessageTypeJoinMatch:
			err := c.handleJoinMatch(message, matchmaker)
			if err != nil {
				log.Println("Error handling join match:", err)
			}
		case MessageTypeMove:
			err := c.handleMove(message)
			if err != nil {
				log.Println("Error handling move:", err)
			}
		default:
			log.Println("Unsupported message type:", message.Type)
		}
	}
}

func (c Client) SendMessages() {
	defer func() { _ = c.Conn.Close() }()
	for message := range c.Player.SendChan {
		err := c.Conn.WriteJSON(message)
		if err != nil {
			log.Println("Write error:", err)
			break
		}
	}
}
