package server

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type WebSocketHandler struct {
	Upgrader websocket.Upgrader
}


func (wsh WebSocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := wsh.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error when upgrading connection to websocket: %s", err)
		return
	}
	defer c.Close()
}
