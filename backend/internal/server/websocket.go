package server

import (
	"encoding/json"
	"go-chess/internal/auth"
	"log"
	"net/http"
	"slices"
	"sync"
	"time"

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
	allowedOrigins []string
	playerCacheTTL time.Duration
	cleanupPeriod  time.Duration
	sessionStore   *auth.SessionStore

	clients   map[*Client]struct{}
	clientsMu sync.RWMutex

	players   map[auth.PlayerKey]*playerCacheEntry
	playersMu sync.RWMutex
}

func NewWebSocketHandler(matchmaker *Matchmaker, allowedOrigins []string, sessionStore *auth.SessionStore) *WebSocketHandler {
	wsh := &WebSocketHandler{
		upgrader:       websocket.Upgrader{},
		matchmaker:     matchmaker,
		metrics:        matchmaker.metrics,
		allowedOrigins: allowedOrigins,
		clients:        make(map[*Client]struct{}),
		playerCacheTTL: defaultPlayerCacheTTL,
		cleanupPeriod:  defaultPlayerCleanupPeriod,
		players:        make(map[auth.PlayerKey]*playerCacheEntry),
		sessionStore:   sessionStore,
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

	for playerKey, entry := range wsh.players {
		if entry.player.HasClients() ||
			entry.player.IsInQueues() ||
			wsh.matchmaker.GetMatch(entry.player) != nil ||
			time.Since(entry.lastUsed) < wsh.playerCacheTTL {
			continue
		}

		delete(wsh.players, playerKey)
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

func (wsh *WebSocketHandler) refreshPlayer(playerKey auth.PlayerKey) {
	if playerKey == "" {
		return
	}

	wsh.playersMu.Lock()
	defer wsh.playersMu.Unlock()

	if entry, ok := wsh.players[playerKey]; ok {
		entry.lastUsed = time.Now()
	}
}

func (wsh *WebSocketHandler) getOrCreatePlayer(session auth.Session) *Player {
	key := session.PlayerKey()

	wsh.playersMu.Lock()
	defer wsh.playersMu.Unlock()

	if entry, ok := wsh.players[key]; ok {
		entry.lastUsed = time.Now()
		return entry.player
	}

	player := NewPlayer(key, session.Username, session.DisplayName)
	wsh.players[key] = &playerCacheEntry{player: player, lastUsed: time.Now()}
	wsh.metrics.recordWebsocketPlayerCached()
	return player
}

func (wsh *WebSocketHandler) sessionFromRequest(r *http.Request) (auth.Session, bool) {
	// Retrieve the session cookie from the handshake request if it exists and is valid.
	var sessionID auth.SessionID
	if cookie, err := r.Cookie(auth.SessionCookieName); err == nil {
		val := auth.SessionID(cookie.Value)
		if auth.IsValidSessionID(val) {
			sessionID = val
		}
	}

	// Look up the active session if a valid session ID was found.
	if sessionID != "" {
		if session, ok := wsh.sessionStore.GetSession(sessionID); ok {
			// Consider the session authenticated if it has a non-empty username
			// (i.e. not an anonymous session).
			isAuthenticated := session.Username != ""
			return session, isAuthenticated
		}
	}

	// No valid session found — fall back to an anonymous session.
	// Reuse the cookie's sessionID if it was valid (preserves reconnect continuity),
	// otherwise generate a new one.
	fallbackID := sessionID
	if fallbackID == "" {
		fallbackID = auth.GenerateSessionID()
	}

	return wsh.sessionStore.CreateAnonSessionWithID(fallbackID), false
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
	if slices.Contains(wsh.allowedOrigins, origin) {
		return true
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

	session, isAuthenticated := wsh.sessionFromRequest(r)
	player := wsh.getOrCreatePlayer(session)
	client := NewClient(conn, player, wsh.metrics)
	player.RegisterClient(client)

	wsh.RegisterClient(client)

	defer func() {
		player.UnregisterClient(client)
		client.Close()
		wsh.matchmaker.LeaveAll(client.Player)
		wsh.UnregisterClient(client)
		wsh.refreshPlayer(player.Key)
		_ = client.Conn.Close()
	}()

	go client.SendMessages()

	// Send initial player info immediately after connection is opened and registered.
	playerInfo := PlayerInfoData{
		Username:        player.Username,
		DisplayName:     player.DisplayName,
		IsAuthenticated: isAuthenticated,
	}
	infoBytes, err := json.Marshal(playerInfo)
	if err != nil {
		log.Printf("WARN: failed to marshal %s for %s: %v", MessageTypePlayerInfo, player.DisplayName, err)
		return
	}
	client.sendChan <- WSMessage{
		Type: MessageTypePlayerInfo,
		Data: json.RawMessage(infoBytes),
	}

	if match := wsh.matchmaker.GetMatch(client.Player); match != nil {
		match.EventChan <- Event{Type: EventTypePlayerReconnected, Player: client.Player}
	}
	queueStats := wsh.matchmaker.GetQueueStats()
	wsh.sendMatchmakingUpdate(client, queueStats)
	client.ReceiveMessages(wsh.matchmaker)
}
