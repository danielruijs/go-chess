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

		rnd := rand.IntN(2)
		randomPosition := chess.Position{}
		switch rnd {
		case 0:
			randomPosition, err = chess.StartingPositionFEN.ToPosition()
			if err != nil {
				log.Println("Error generating starting position:", err)
				continue
			}
		case 1:
			randomPosition = chess.Position{}
			for i := range chess.BoardSize {
				for j := range chess.BoardSize {
					colors := []chess.Color{chess.White, chess.Black}
					color := colors[rand.IntN(len(colors))]
					randomPosition.Board[i][j] = chess.Piece{
						Type:  chess.Pawn,
						Color: color,
					}
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
