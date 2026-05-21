package server

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/require"
)

func newTestWebSocketServer(t *testing.T) (*WebSocketHandler, *httptest.Server) {
	t.Helper()

	metrics := NewMetrics()
	matchmaker := NewMatchmaker(metrics)
	go matchmaker.Run()

	handler := NewWebSocketHandler(matchmaker, true)
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

func TestWebSocketHandlerAllowsLocalhostOrigins(t *testing.T) {
	_, server := newTestWebSocketServer(t)
	t.Cleanup(server.Close)

	for _, origin := range []string{
		"http://localhost:3000",
		"http://127.0.0.1:3000",
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

	first := handler.getOrCreatePlayer("123e4567-e89b-12d3-a456-426614174000")
	first.Name = "Alice"

	second := handler.getOrCreatePlayer("123e4567-e89b-12d3-a456-426614174000")

	require.Same(t, first, second)
	require.Equal(t, "Alice", second.Name)
}

func TestWebSocketHandlerDoesNotCacheInvalidSessionID(t *testing.T) {
	handler, server := newTestWebSocketServer(t)
	t.Cleanup(server.Close)

	invalidSessionID := strings.Repeat("a", 128)

	first := handler.getOrCreatePlayer(invalidSessionID)
	second := handler.getOrCreatePlayer(invalidSessionID)

	require.NotSame(t, first, second)
	require.Empty(t, handler.players)
}

func TestPlayerBroadcastsToAllRegisteredClients(t *testing.T) {
	player := NewPlayer()
	clientOne := &Client{sendChan: make(chan WSMessage, 1)}
	clientTwo := &Client{sendChan: make(chan WSMessage, 1)}

	player.RegisterClient(clientOne)
	player.RegisterClient(clientTwo)

	err := player.SendCritical(MessageTypeStartMatch, nil)
	require.NoError(t, err)

	messageOne := <-clientOne.sendChan
	messageTwo := <-clientTwo.sendChan

	require.Equal(t, MessageTypeStartMatch, messageOne.Type)
	require.Equal(t, MessageTypeStartMatch, messageTwo.Type)
}
