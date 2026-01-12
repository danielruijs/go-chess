package server

import (
	"go-chess/internal/chess"
	"log"
	"math/rand/v2"
	"net/http"

	"github.com/gorilla/websocket"
)

type WebSocketHandler struct {
	Upgrader websocket.Upgrader
}


func (wsh WebSocketHandler) ServeWS(w http.ResponseWriter, r *http.Request) {
	wsh.Upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	conn, err := wsh.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error when upgrading connection to websocket: %s", err)
		return
	}
	defer func() { _ = conn.Close() }()

	log.Println("New WebSocket connection established")

	for {
		// Read message from client
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		log.Printf("Received message: %s. Message type: %d", message, messageType)

		randomPosition := chess.Position{}
		for i := range chess.BoardSize{
			for j := range chess.BoardSize{
				colors := []chess.Color{chess.White, chess.Black}
				color := colors[rand.IntN(len(colors))]
				randomPosition.Board[i][j] = chess.Piece{
					Type:  "pawn",
					Color: color,
				}
			}
		}

		err = conn.WriteJSON(chess.WSMessage{
			Type: chess.MessageTypePosition,
			Data: randomPosition,
		})
		if err != nil {
			log.Println("Write error:", err)
			break
		}
	}
}
