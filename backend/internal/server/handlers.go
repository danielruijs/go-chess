package server

import (
	"log"
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
	defer conn.Close()

	log.Println("New WebSocket connection established")

	for {
		// Read message from client
		mt, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		log.Printf("Received message: %s", message)

		// Echo back to client for now
		if err := conn.WriteMessage(mt, append([]byte("Server: "), message...)); err != nil {
			log.Println("Write error:", err)
			break
		}
	}
}
