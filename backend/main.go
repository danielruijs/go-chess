package main

import (
	"go-chess/internal/server"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func main() {
	webSocketHandler := server.WebSocketHandler{
		Upgrader: websocket.Upgrader{},
	}

	http.HandleFunc("/ws", webSocketHandler.ServeWS)

	log.Print("Starting server...")
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
