package main

import (
	"go-chess/internal/server"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func main() {
	matchmaker := server.NewMatchmaker()
	go matchmaker.Run()

	webSocketHandler := server.WebSocketHandler{
		Upgrader:   websocket.Upgrader{},
		Matchmaker: matchmaker,
	}

	http.HandleFunc("/ws", webSocketHandler.ServeWS)

	log.Print("Started server")
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
