package server

import (
	"go-chess/internal/chess"
	"go-chess/internal/matchmaker"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type WebSocketHandler struct {
	Upgrader   websocket.Upgrader
	Matchmaker *matchmaker.Matchmaker
}

type Client struct {
	Conn   *websocket.Conn
	Player *chess.Player
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
		Player: &chess.Player{
			Name:     "test",
			SendChan: make(chan chess.WSMessage),
		},
	}

	go client.ReceiveMessages(wsh.Matchmaker)
	go client.SendMessages()

	log.Println("New WebSocket connection established")
}

func (c Client) ReceiveMessages(matchmaker *matchmaker.Matchmaker) {
	defer func() { _ = c.Conn.Close() }()
	for {
		var message chess.WSMessage
		err := c.Conn.ReadJSON(&message)
		if err != nil {
			log.Println("Failed to unmarshal message:", err)
			return
		}

		switch message.Type {
		case chess.MessageTypeJoinMatch:
			err := c.handleJoinMatch(matchmaker)
			if err != nil {
				log.Println("Error handling join match:", err)
			}
		case chess.MessageTypeMove:
			err := c.handleMove(message)
			if err != nil {
				log.Println("Error handling move:", err)
			}
			continue
		default:
			log.Println("Unsupported message type:", message.Type)
			return
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
