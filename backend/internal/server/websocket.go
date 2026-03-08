package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type WebSocketHandler struct {
	upgrader   websocket.Upgrader
	matchmaker *Matchmaker
	clients    map[*Client]struct{}
}

func NewWebSocketHandler(matchmaker *Matchmaker) *WebSocketHandler {
	wsh := &WebSocketHandler{
		upgrader:   websocket.Upgrader{},
		matchmaker: matchmaker,
		clients:    make(map[*Client]struct{}),
	}

	go func() {
		for range matchmaker.UpdateChan {
			wsh.BroadcastMatchmakingUpdates()
		}
	}()

	return wsh
}

func (wsh *WebSocketHandler) RegisterClient(client *Client) {
	wsh.clients[client] = struct{}{}
	log.Printf("Client registered. Total clients: %d\n", len(wsh.clients))
}

func (wsh *WebSocketHandler) UnregisterClient(client *Client) {
	delete(wsh.clients, client)
	log.Printf("Client unregistered. Total clients: %d\n", len(wsh.clients))
}

func (wsh *WebSocketHandler) BroadcastMatchmakingUpdates() {
	for client := range wsh.clients {
		wsh.sendMatchmakingUpdate(client)
	}
}

func (wsh *WebSocketHandler) sendMatchmakingUpdate(client *Client) {
	updateData := wsh.matchmaker.GetMatchmakingUpdate(client.Player)
	data, err := json.Marshal(updateData)
	if err != nil {
		log.Printf("failed to marshal matchmaking update data: %v", err)
		return
	}
	client.Player.SendChan <- WSMessage{
		Type: MessageTypeMatchmakingUpdate,
		Data: data,
	}
}

func (wsh *WebSocketHandler) ServeWS(w http.ResponseWriter, r *http.Request) {
	wsh.upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	conn, err := wsh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error when upgrading connection to websocket: %s", err)
		return
	}

	client := &Client{
		Conn: conn,
		Player: &Player{
			Name:     "",
			SendChan: make(chan WSMessage),
		},
	}

	wsh.RegisterClient(client)
	log.Println("New WebSocket connection established")

	defer func() {
		wsh.matchmaker.Leave(client.Player)
		wsh.UnregisterClient(client)
		close(client.Player.SendChan)
		_ = client.Conn.Close()
	}()

	go client.SendMessages()
	wsh.sendMatchmakingUpdate(client)
	client.ReceiveMessages(wsh.matchmaker)
}
