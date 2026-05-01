package server

import (
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

type WebSocketHandler struct {
	upgrader       websocket.Upgrader
	matchmaker     *Matchmaker
	metrics        *metrics
	allowLocalhost bool
	clients        map[*Client]struct{}
	clientsMu      sync.RWMutex
}

func NewWebSocketHandler(matchmaker *Matchmaker, allowLocalhost bool) *WebSocketHandler {
	wsh := &WebSocketHandler{
		upgrader:       websocket.Upgrader{},
		matchmaker:     matchmaker,
		metrics:        matchmaker.metrics,
		allowLocalhost: allowLocalhost,
		clients:        make(map[*Client]struct{}),
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
	wsh.metrics.recordWebsocketConnectionOpened()
	log.Printf("Client registered. Total clients: %d\n", len(wsh.clients))
}

func (wsh *WebSocketHandler) UnregisterClient(client *Client) {
	wsh.clientsMu.Lock()
	defer wsh.clientsMu.Unlock()
	delete(wsh.clients, client)
	wsh.metrics.recordWebsocketConnectionClosed()
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
	wsh.upgrader.CheckOrigin = func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		// Allow production origin
		if origin == "https://gochess.dev" {
			return true
		}
		// Allow local development origins
		if wsh.allowLocalhost && strings.HasPrefix(origin, "http://localhost") {
			return true
		}
		return false
	}

	conn, err := wsh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error when upgrading connection to websocket: %s", err)
		return
	}

	client := NewClient(conn, wsh.metrics)

	wsh.RegisterClient(client)

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
