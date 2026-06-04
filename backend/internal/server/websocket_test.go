package server

import (
	"encoding/json"
	"go-chess/internal/auth"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/require"
)

func newTestWebSocketServer(t *testing.T) (*WebSocketHandler, *httptest.Server) {
	t.Helper()

	metrics := NewMetrics()
	matchmaker := NewMatchmaker(metrics)
	go matchmaker.Run()

	sessionStore := auth.NewSessionStore()
	handler := NewWebSocketHandler(matchmaker, []string{"http://localhost:5173", "http://127.0.0.1:5173"}, sessionStore)
	server := httptest.NewServer(http.HandlerFunc(handler.ServeWS))

	return handler, server
}

func dialTestWebSocket(t *testing.T, serverURL string, origin string) (*websocket.Conn, *http.Response) {
	t.Helper()

	parsedURL, err := url.Parse(serverURL)
	require.NoError(t, err)
	parsedURL.Scheme = "ws"

	header := http.Header{}
	header.Set("Origin", origin)

	conn, response, err := websocket.DefaultDialer.Dial(parsedURL.String(), header)
	require.NoError(t, err)

	return conn, response
}

func dialTestWebSocketWithCookie(t *testing.T, serverURL string, origin string, sessionID auth.SessionID) (*websocket.Conn, *http.Response) {
	t.Helper()

	parsedURL, err := url.Parse(serverURL)
	require.NoError(t, err)
	parsedURL.Scheme = "ws"

	header := http.Header{}
	header.Set("Origin", origin)
	header.Set("Cookie", string(auth.SessionCookieName)+"="+string(sessionID))

	conn, response, err := websocket.DefaultDialer.Dial(parsedURL.String(), header)
	require.NoError(t, err)

	return conn, response
}

func sendTestMessage(t *testing.T, conn *websocket.Conn, messageType MessageType, data any) {
	t.Helper()

	rawData, err := json.Marshal(data)
	require.NoError(t, err)

	require.NoError(t, conn.WriteJSON(WSMessage{Type: messageType, Data: rawData}))
}

func readTestMessage(t *testing.T, conn *websocket.Conn) WSMessage {
	t.Helper()

	var message WSMessage
	require.NoError(t, conn.ReadJSON(&message))

	return message
}

func waitForTestMessageType(t *testing.T, conn *websocket.Conn, expectedType MessageType) WSMessage {
	t.Helper()

	require.NoError(t, conn.SetReadDeadline(time.Now().Add(5*time.Second)))
	defer func() {
		_ = conn.SetReadDeadline(time.Time{})
	}()

	for {
		message := readTestMessage(t, conn)
		if message.Type == expectedType {
			return message
		}
	}
}

func TestWebSocketHandlerAllowsLocalhostOrigins(t *testing.T) {
	_, server := newTestWebSocketServer(t)
	t.Cleanup(server.Close)

	for _, origin := range []string{
		"http://localhost:5173",
		"http://127.0.0.1:5173",
	} {
		t.Run(origin, func(t *testing.T) {
			conn, response := dialTestWebSocket(t, server.URL, origin)
			t.Cleanup(func() {
				_ = conn.Close()
				if response != nil && response.Body != nil {
					_ = response.Body.Close()
				}
			})
		})
	}
}

func TestWebSocketHandlerDeniesOriginBypassAndCountsIt(t *testing.T) {
	handler, server := newTestWebSocketServer(t)
	t.Cleanup(server.Close)

	parsedURL, err := url.Parse(server.URL)
	require.NoError(t, err)
	parsedURL.Scheme = "ws"

	header := http.Header{}
	header.Set("Origin", "http://localhost.evil.com")

	conn, response, err := websocket.DefaultDialer.Dial(parsedURL.String(), header)
	require.Error(t, err)
	require.Nil(t, conn)
	if response != nil && response.Body != nil {
		t.Cleanup(func() {
			_ = response.Body.Close()
		})
	}

	require.InDelta(t, 1.0, testutil.ToFloat64(handler.metrics.connectionsDenied), 0.0001)
}

func TestWebSocketHandlerReusesPlayerForSessionID(t *testing.T) {
	handler, server := newTestWebSocketServer(t)
	t.Cleanup(server.Close)

	sessionID := auth.SessionID("123e4567-e89b-12d3-a456-426614174000")
	session := handler.sessionStore.CreateSessionWithID(sessionID, "", "Anon")
	first := handler.getOrCreatePlayer(session)
	first.DisplayName = "Alice"

	second := handler.getOrCreatePlayer(session)

	require.Same(t, first, second)
	require.Equal(t, "Alice", second.DisplayName)
}

func TestWebSocketHandlerRestoresActiveMatchOnReconnect(t *testing.T) {
	_, server := newTestWebSocketServer(t)
	t.Cleanup(server.Close)

	const origin = "http://localhost:5173"
	timeFormat := TimeFormat{initial: time.Minute, increment: 0}
	firstSessionID := auth.SessionID("123e4567-e89b-12d3-a456-426614174000")
	secondSessionID := auth.SessionID("123e4567-e89b-12d3-a456-426614174001")

	firstConn, firstResponse := dialTestWebSocketWithCookie(t, server.URL, origin, firstSessionID)
	t.Cleanup(func() {
		_ = firstConn.Close()
		if firstResponse != nil && firstResponse.Body != nil {
			_ = firstResponse.Body.Close()
		}
	})

	secondConn, secondResponse := dialTestWebSocketWithCookie(t, server.URL, origin, secondSessionID)
	t.Cleanup(func() {
		_ = secondConn.Close()
		if secondResponse != nil && secondResponse.Body != nil {
			_ = secondResponse.Body.Close()
		}
	})

	// Discard player_info message and matchmaking_update message
	require.Equal(t, MessageTypePlayerInfo, readTestMessage(t, firstConn).Type)
	require.Equal(t, MessageTypeMatchmakingUpdate, readTestMessage(t, firstConn).Type)
	require.Equal(t, MessageTypePlayerInfo, readTestMessage(t, secondConn).Type)
	require.Equal(t, MessageTypeMatchmakingUpdate, readTestMessage(t, secondConn).Type)

	// Both clients join the same matchmaking queue with the same time format, which should pair them together
	sendTestMessage(t, firstConn, MessageTypeJoinMatch, JoinMatchData{
		TimeFormat: TimeFormatMs{InitialMs: timeFormat.initial.Milliseconds(), IncrementMs: timeFormat.increment.Milliseconds()},
	})
	sendTestMessage(t, secondConn, MessageTypeJoinMatch, JoinMatchData{
		TimeFormat: TimeFormatMs{InitialMs: timeFormat.initial.Milliseconds(), IncrementMs: timeFormat.increment.Milliseconds()},
	})

	// Wait for both clients to receive the start match and board messages
	require.Equal(t, MessageTypeStartMatch, waitForTestMessageType(t, firstConn, MessageTypeStartMatch).Type)
	require.Equal(t, MessageTypeStartMatch, waitForTestMessageType(t, secondConn, MessageTypeStartMatch).Type)
	require.Equal(t, MessageTypeBoard, waitForTestMessageType(t, firstConn, MessageTypeBoard).Type)
	require.Equal(t, MessageTypeBoard, waitForTestMessageType(t, secondConn, MessageTypeBoard).Type)

	// Simulate a disconnect and reconnect for the first player
	require.NoError(t, firstConn.Close())
	reconnectConn, reconnectResponse := dialTestWebSocketWithCookie(t, server.URL, origin, firstSessionID)
	t.Cleanup(func() {
		_ = reconnectConn.Close()
		if reconnectResponse != nil && reconnectResponse.Body != nil {
			_ = reconnectResponse.Body.Close()
		}
	})

	// Discard player_info message
	require.Equal(t, MessageTypePlayerInfo, readTestMessage(t, reconnectConn).Type)

	// The first player should receive the start match and board messages again, indicating that their active match was restored
	require.Equal(t, MessageTypeStartMatch, waitForTestMessageType(t, reconnectConn, MessageTypeStartMatch).Type)
	require.Equal(t, MessageTypeBoard, waitForTestMessageType(t, reconnectConn, MessageTypeBoard).Type)
}

func TestWebSocketHandlerTracksCachedPlayerMetric(t *testing.T) {
	handler, server := newTestWebSocketServer(t)
	t.Cleanup(server.Close)

	sessionID := auth.SessionID("123e4567-e89b-12d3-a456-426614174000")
	session := handler.sessionStore.CreateSessionWithID(sessionID, "", "Anon")
	player := handler.getOrCreatePlayer(session)
	require.Same(t, player, handler.players[player.Key].player)
	require.InDelta(t, 1.0, testutil.ToFloat64(handler.metrics.cachedPlayers), 0.0001)

	handler.players[player.Key].lastUsed = time.Now().Add(-2 * handler.playerCacheTTL)
	handler.cleanupInactivePlayers()

	require.InDelta(t, 0.0, testutil.ToFloat64(handler.metrics.cachedPlayers), 0.0001)
}

func TestWebSocketHandlerCleansUpInactiveCachedPlayers(t *testing.T) {
	handler, server := newTestWebSocketServer(t)
	t.Cleanup(server.Close)

	sessionID := auth.SessionID("123e4567-e89b-12d3-a456-426614174000")
	session := handler.sessionStore.CreateSessionWithID(sessionID, "", "Anon")
	player := handler.getOrCreatePlayer(session)
	require.Same(t, player, handler.players[player.Key].player)

	handler.players[player.Key].lastUsed = time.Now().Add(-2 * handler.playerCacheTTL)
	handler.cleanupInactivePlayers()

	_, ok := handler.players[player.Key]
	require.False(t, ok)
}

func TestWebSocketHandlerKeepsCachedPlayersWithActiveClients(t *testing.T) {
	handler, server := newTestWebSocketServer(t)
	t.Cleanup(server.Close)

	sessionID := auth.SessionID("123e4567-e89b-12d3-a456-426614174000")
	session := handler.sessionStore.CreateSessionWithID(sessionID, "", "Anon")
	player := handler.getOrCreatePlayer(session)
	player.RegisterClient(&Client{sendChan: make(chan WSMessage, 1)})
	handler.players[player.Key].lastUsed = time.Now().Add(-2 * handler.playerCacheTTL)

	handler.cleanupInactivePlayers()

	_, ok := handler.players[player.Key]
	require.True(t, ok)
}

func TestWebSocketHandlerKeepsCachedPlayersInQueue(t *testing.T) {
	handler, server := newTestWebSocketServer(t)
	t.Cleanup(server.Close)

	timeFormat := TimeFormat{initial: time.Minute, increment: 0}
	sessionID := auth.SessionID("123e4567-e89b-12d3-a456-426614174000")
	session := handler.sessionStore.CreateSessionWithID(sessionID, "", "Anon")
	player := handler.getOrCreatePlayer(session)
	require.NoError(t, handler.matchmaker.Join(player, timeFormat))
	handler.players[player.Key].lastUsed = time.Now().Add(-2 * handler.playerCacheTTL)

	handler.cleanupInactivePlayers()

	_, ok := handler.players[player.Key]
	require.True(t, ok)
}

func TestWebSocketHandlerKeepsCachedPlayersInMatch(t *testing.T) {
	handler, server := newTestWebSocketServer(t)
	t.Cleanup(server.Close)

	timeFormat := TimeFormat{initial: time.Minute, increment: 0}
	session1 := handler.sessionStore.CreateSessionWithID(auth.SessionID("123e4567-e89b-12d3-a456-426614174000"), "", "Anon")
	player1 := handler.getOrCreatePlayer(session1)
	player2 := NewPlayer("anon:test-2", "", "Test Player 2")
	match := NewMatch(player1, player2, timeFormat, handler.matchmaker.matchEnded, handler.matchmaker.metrics)
	handler.matchmaker.RegisterMatch(match)
	handler.players[player1.Key].lastUsed = time.Now().Add(-2 * handler.playerCacheTTL)

	handler.cleanupInactivePlayers()

	_, ok := handler.players[player1.Key]
	require.True(t, ok)
}

func TestWebSocketHandlerDifferentSessionsGetDifferentPlayers(t *testing.T) {
	handler, server := newTestWebSocketServer(t)
	t.Cleanup(server.Close)

	session1 := handler.sessionStore.CreateSession("", "Anon 1")
	session2 := handler.sessionStore.CreateSession("", "Anon 2")

	first := handler.getOrCreatePlayer(session1)
	second := handler.getOrCreatePlayer(session2)

	require.NotSame(t, first, second)
}

func TestPlayerBroadcastsToAllRegisteredClients(t *testing.T) {
	player := NewPlayer("anon:test-1", "", "Test Player")
	clientOne := &Client{sendChan: make(chan WSMessage, 1)}
	clientTwo := &Client{sendChan: make(chan WSMessage, 1)}

	player.RegisterClient(clientOne)
	player.RegisterClient(clientTwo)

	player.Send(MessageTypeStartMatch, nil)

	messageOne := <-clientOne.sendChan
	messageTwo := <-clientTwo.sendChan

	require.Equal(t, MessageTypeStartMatch, messageOne.Type)
	require.Equal(t, MessageTypeStartMatch, messageTwo.Type)
}
