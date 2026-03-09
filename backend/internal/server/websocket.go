package server

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type WebSocketHandler struct {
	upgrader   websocket.Upgrader
	matchmaker *Matchmaker
	clients    map[*Client]struct{}
	clientsMu  sync.RWMutex
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
	wsh.clientsMu.Lock()
	defer wsh.clientsMu.Unlock()
	wsh.clients[client] = struct{}{}
	log.Printf("Client registered. Total clients: %d\n", len(wsh.clients))
}

func (wsh *WebSocketHandler) UnregisterClient(client *Client) {
	wsh.clientsMu.Lock()
	defer wsh.clientsMu.Unlock()
	delete(wsh.clients, client)
	log.Printf("Client unregistered. Total clients: %d\n", len(wsh.clients))
}

func (wsh *WebSocketHandler) BroadcastMatchmakingUpdates() {
	queueStats := wsh.matchmaker.GetQueueStats()
	for _, client := range wsh.getClientsSnapshot() {
		wsh.sendMatchmakingUpdate(client, queueStats)
	}
}

func (wsh *WebSocketHandler) getClientsSnapshot() []*Client {
	wsh.clientsMu.RLock()
	defer wsh.clientsMu.RUnlock()

	clients := make([]*Client, 0, len(wsh.clients))
	for client := range wsh.clients {
		clients = append(clients, client)
	}
	return clients
}

func (wsh *WebSocketHandler) sendMatchmakingUpdate(client *Client, queueStats map[TimeFormat]int) {
	queues := make([]QueueData, 0, len(queueStats))
	for timeFormat, queueLength := range queueStats {
		queues = append(queues, QueueData{
			TimeFormat:  TimeFormatToMs(timeFormat),
			QueueLength: queueLength,
			InQueue:     client.Player.IsInQueue(timeFormat),
		})
	}
	updateData := MatchmakingUpdateData{Queues: queues}
	client.Player.SendInformational(MessageTypeMatchmakingUpdate, updateData)
}

func (wsh *WebSocketHandler) ServeWS(w http.ResponseWriter, r *http.Request) {
	wsh.upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	conn, err := wsh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error when upgrading connection to websocket: %s", err)
		return
	}

	client := NewClient(conn)

	wsh.RegisterClient(client)
	log.Println("New WebSocket connection established")

	defer func() {
		close(client.Done)
		wsh.matchmaker.Leave(client.Player)
		wsh.UnregisterClient(client)
		_ = client.Conn.Close()
	}()

	go client.SendMessages()
	queueStats := wsh.matchmaker.GetQueueStats()
	wsh.sendMatchmakingUpdate(client, queueStats)
	client.ReceiveMessages(wsh.matchmaker)
}
