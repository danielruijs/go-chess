package server

import (
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type playerCacheEntry struct {
	player   *Player
	lastUsed time.Time
}

const (
	defaultPlayerCacheTTL      = 15 * time.Minute
	defaultPlayerCleanupPeriod = 1 * time.Minute
)

type WebSocketHandler struct {
	upgrader       websocket.Upgrader
	matchmaker     *Matchmaker
	metrics        *metrics
	allowLocalhost bool
	playerCacheTTL time.Duration
	cleanupPeriod  time.Duration

	clients   map[*Client]struct{}
	clientsMu sync.RWMutex

	players   map[string]*playerCacheEntry
	playersMu sync.RWMutex
}

const SESSION_ID_COOKIE_NAME = "go-chess.sessionId"

func NewWebSocketHandler(matchmaker *Matchmaker, allowLocalhost bool) *WebSocketHandler {
	wsh := &WebSocketHandler{
		upgrader:       websocket.Upgrader{},
		matchmaker:     matchmaker,
		metrics:        matchmaker.metrics,
		allowLocalhost: allowLocalhost,
		clients:        make(map[*Client]struct{}),
		playerCacheTTL: defaultPlayerCacheTTL,
		cleanupPeriod:  defaultPlayerCleanupPeriod,
		players:        make(map[string]*playerCacheEntry),
	}
	wsh.upgrader.CheckOrigin = wsh.checkOrigin

	go func() {
		for range matchmaker.UpdateChan {
			wsh.BroadcastMatchmakingUpdates()
		}
	}()

	go wsh.runPlayerCleanup()

	return wsh
}

func (wsh *WebSocketHandler) runPlayerCleanup() {
	ticker := time.NewTicker(wsh.cleanupPeriod)
	defer ticker.Stop()

	for range ticker.C {
		wsh.cleanupInactivePlayers()
	}
}

func (wsh *WebSocketHandler) cleanupInactivePlayers() {
	wsh.playersMu.Lock()
	defer wsh.playersMu.Unlock()

	for sessionID, entry := range wsh.players {
		if entry.player.HasClients() ||
			entry.player.IsInQueues() ||
			wsh.matchmaker.GetMatch(entry.player) != nil ||
			time.Since(entry.lastUsed) < wsh.playerCacheTTL {
			continue
		}

		delete(wsh.players, sessionID)
		wsh.metrics.recordWebsocketPlayerEvicted()
	}
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

func (wsh *WebSocketHandler) refreshPlayer(sessionID string) {
	if sessionID == "" {
		return
	}

	wsh.playersMu.Lock()
	defer wsh.playersMu.Unlock()

	if entry, ok := wsh.players[sessionID]; ok {
		entry.lastUsed = time.Now()
	}
}

func (wsh *WebSocketHandler) getOrCreatePlayer(sessionID string) *Player {
	if sessionID == "" || !isValidSessionID(sessionID) {
		return NewPlayer()
	}

	wsh.playersMu.Lock()
	defer wsh.playersMu.Unlock()

	if entry, ok := wsh.players[sessionID]; ok {
		entry.lastUsed = time.Now()
		return entry.player
	}

	player := NewPlayer()
	wsh.players[sessionID] = &playerCacheEntry{player: player, lastUsed: time.Now()}
	wsh.metrics.recordWebsocketPlayerCached()
	return player
}

func (wsh *WebSocketHandler) sessionIDFromRequest(r *http.Request) string {
	cookie, err := r.Cookie(SESSION_ID_COOKIE_NAME)
	if err != nil {
		return ""
	}
	if !isValidSessionID(cookie.Value) {
		return ""
	}
	return cookie.Value
}

func isValidSessionID(sessionID string) bool {
	_, err := uuid.Parse(sessionID)
	return err == nil
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
	client.Player.Send(MessageTypeMatchmakingUpdate, updateData)
}

func (wsh *WebSocketHandler) checkOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	// Allow production origin
	if origin == "https://gochess.dev" {
		return true
	}

	// Allow local development origins
	if wsh.allowLocalhost {
		u, err := url.Parse(origin)
		if err == nil {
			hostname := u.Hostname()
			if hostname == "localhost" || hostname == "127.0.0.1" {
				return true
			}
		}
	}

	wsh.metrics.recordWebsocketConnectionDenied()
	return false
}

func (wsh *WebSocketHandler) ServeWS(w http.ResponseWriter, r *http.Request) {
	conn, err := wsh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error when upgrading connection to websocket: %s", err)
		return
	}

	sessionID := wsh.sessionIDFromRequest(r)
	player := wsh.getOrCreatePlayer(sessionID)
	client := NewClient(conn, player, wsh.metrics)
	player.RegisterClient(client)

	wsh.RegisterClient(client)

	defer func() {
		player.UnregisterClient(client)
		close(client.Done)
		wsh.matchmaker.LeaveAll(client.Player)
		wsh.UnregisterClient(client)
		wsh.refreshPlayer(sessionID)
		_ = client.Conn.Close()
	}()

	go client.SendMessages()
	if match := wsh.matchmaker.GetMatch(client.Player); match != nil {
		match.sendCurrentState(client.Player)
	}
	queueStats := wsh.matchmaker.GetQueueStats()
	wsh.sendMatchmakingUpdate(client, queueStats)
	client.ReceiveMessages(wsh.matchmaker)
}
