package server

import (
	"context"
	"encoding/json"
	"go-chess/internal/auth"
	"go-chess/internal/cache"
	"log"
	"net/http"
	"slices"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	playerCacheTTL        = 15 * time.Minute
	playerCleanupInterval = 1 * time.Minute
)

type WebSocketHandler struct {
	upgrader       websocket.Upgrader
	matchmaker     *Matchmaker
	metrics        *metrics
	allowedOrigins []string
	sessionStore   *auth.SessionStore

	clients   map[*Client]struct{}
	clientsMu sync.RWMutex

	players *cache.Cache[PlayerKey, *Player]
}

func NewWebSocketHandler(ctx context.Context, matchmaker *Matchmaker, allowedOrigins []string, sessionStore *auth.SessionStore) (*WebSocketHandler, error) {
	wsh := &WebSocketHandler{
		upgrader:       websocket.Upgrader{},
		matchmaker:     matchmaker,
		metrics:        matchmaker.metrics,
		allowedOrigins: allowedOrigins,
		clients:        make(map[*Client]struct{}),
		sessionStore:   sessionStore,
	}
	wsh.upgrader.CheckOrigin = wsh.checkOrigin

	cache, err := cache.New[PlayerKey](cache.Options[*Player]{
		Cleanup: &cache.CleanupConfig[*Player]{
			Interval:    playerCleanupInterval,
			ShouldEvict: wsh.shouldEvictPlayer,
			OnEvicted: func(count int) {
				wsh.metrics.recordWebsocketPlayerEvicted(count)
			},
		},
	})
	if err != nil {
		return nil, err
	}
	wsh.players = cache
	wsh.players.StartCleanup(ctx)

	go func() {
		for range matchmaker.UpdateChan {
			wsh.BroadcastMatchmakingUpdates()
		}
	}()

	return wsh, nil
}

func (wsh *WebSocketHandler) shouldEvictPlayer(player *Player, lastUsed time.Time) bool {
	if player.HasClients() ||
		player.IsInQueues() ||
		wsh.matchmaker.GetMatch(player) != nil ||
		time.Since(lastUsed) < playerCacheTTL {
		return false
	}
	return true
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

func (wsh *WebSocketHandler) refreshPlayer(playerKey PlayerKey) {
	if playerKey == "" {
		return
	}

	wsh.players.Touch(playerKey)
}

func (wsh *WebSocketHandler) getOrCreatePlayer(session auth.Session) *Player {
	key := NewPlayerKey(session)

	return wsh.players.GetOrCreate(key, func() *Player {
		player := NewPlayer(key, session.Username, session.DisplayName)
		wsh.metrics.recordWebsocketPlayerCached()
		return player
	})
}

func (wsh *WebSocketHandler) sessionFromRequest(r *http.Request) (auth.Session, bool) {
	sessionID, ok := auth.SessionIDFromRequest(r)

	// Look up the active session if a valid session ID was found.
	if ok {
		if session, exists := wsh.sessionStore.GetSession(sessionID); exists {
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
	if !ok {
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
