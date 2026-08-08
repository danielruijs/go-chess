package server

import (
	"net/http"

	"go-chess/internal/chess"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type metrics struct {
	registry *prometheus.Registry

	onlineConnections    prometheus.Gauge
	cachedPlayers        prometheus.Gauge
	connectionsTotal     prometheus.Counter
	connectionsClosed    prometheus.Counter
	connectionsDenied    prometheus.Counter
	messagesSent         *prometheus.CounterVec
	messagesReceived     *prometheus.CounterVec
	messageSendErrors    *prometheus.CounterVec
	messageReceiveErrors *prometheus.CounterVec

	queueSize        *prometheus.GaugeVec
	queueJoinsTotal  *prometheus.CounterVec
	queueLeavesTotal prometheus.Counter

	activeMatches    prometheus.Gauge
	matchesStarted   *prometheus.CounterVec
	matchesFinished  prometheus.Counter
	eventsTotal      *prometheus.CounterVec
	eventErrorsTotal *prometheus.CounterVec
	resultsTotal     *prometheus.CounterVec
}

func NewMetrics() *metrics {
	registry := prometheus.NewRegistry()
	registry.MustRegister(collectors.NewGoCollector())
	registry.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	m := &metrics{
		registry: registry,
		onlineConnections: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "go_chess_websocket_connections_active",
			Help: "Current number of open websocket connections.",
		}),
		cachedPlayers: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "go_chess_websocket_players_cached",
			Help: "Current number of cached websocket players.",
		}),
		connectionsTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "go_chess_websocket_connections_total",
			Help: "Total number of websocket connections opened.",
		}),
		connectionsClosed: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "go_chess_websocket_connections_closed_total",
			Help: "Total number of websocket connections closed.",
		}),
		connectionsDenied: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "go_chess_websocket_connections_denied_total",
			Help: "Total number of websocket connections denied due to origin checks.",
		}),
		messagesSent: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "go_chess_websocket_messages_sent_total",
			Help: "Total websocket messages sent to player clients by message type.",
		}, []string{"type"}),
		messagesReceived: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "go_chess_websocket_messages_received_total",
			Help: "Total websocket messages received by message type.",
		}, []string{"type"}),
		messageSendErrors: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "go_chess_websocket_message_send_errors_total",
			Help: "Total websocket message send errors by type.",
		}, []string{"type"}),
		messageReceiveErrors: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "go_chess_websocket_message_receive_errors_total",
			Help: "Total websocket message handling errors by type.",
		}, []string{"type"}),
		queueSize: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "go_chess_matchmaking_queue_size",
			Help: "Current queue size per time format.",
		}, []string{"time_format"}),
		queueJoinsTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "go_chess_matchmaking_joins_total",
			Help: "Total matchmaking queue joins by time format.",
		}, []string{"time_format"}),
		queueLeavesTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "go_chess_matchmaking_leaves_total",
			Help: "Total matchmaking queue leaves.",
		}),
		activeMatches: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "go_chess_matches_active",
			Help: "Current number of active matches.",
		}),
		matchesStarted: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "go_chess_matches_started_total",
			Help: "Total matches started by time format.",
		}, []string{"time_format"}),
		matchesFinished: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "go_chess_matches_finished_total",
			Help: "Total matches finished.",
		}),
		eventsTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "go_chess_match_events_total",
			Help: "Total match events by type.",
		}, []string{"type"}),
		eventErrorsTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "go_chess_match_event_errors_total",
			Help: "Total match event handling errors by event type.",
		}, []string{"type"}),
		resultsTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "go_chess_match_results_total",
			Help: "Total match results by outcome and reason.",
		}, []string{"outcome", "reason"}),
	}

	registry.MustRegister(
		m.onlineConnections,
		m.cachedPlayers,
		m.connectionsTotal,
		m.connectionsClosed,
		m.connectionsDenied,
		m.messagesSent,
		m.messagesReceived,
		m.messageSendErrors,
		m.messageReceiveErrors,
		m.queueSize,
		m.queueJoinsTotal,
		m.queueLeavesTotal,
		m.activeMatches,
		m.matchesStarted,
		m.matchesFinished,
		m.eventsTotal,
		m.eventErrorsTotal,
		m.resultsTotal,
	)

	return m
}

func (m *metrics) MetricsHandler() http.Handler {
	return promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{})
}

func (m *metrics) recordWebsocketConnectionOpened() {
	m.onlineConnections.Inc()
	m.connectionsTotal.Inc()
}

func (m *metrics) recordWebsocketConnectionClosed() {
	m.onlineConnections.Dec()
	m.connectionsClosed.Inc()
}

func (m *metrics) recordWebsocketPlayerCached() {
	m.cachedPlayers.Inc()
}

func (m *metrics) recordWebsocketPlayerEvicted(count int) {
	m.cachedPlayers.Sub(float64(count))
}

func (m *metrics) recordWebsocketConnectionDenied() {
	m.connectionsDenied.Inc()
}

func (m *metrics) recordWebsocketMessageSent(messageType MessageType) {
	m.messagesSent.WithLabelValues(string(messageType)).Inc()
}

func (m *metrics) recordWebsocketMessageReceived(messageType MessageType) {
	m.messagesReceived.WithLabelValues(string(messageType)).Inc()
}

func (m *metrics) recordWebsocketMessageSendError(messageType MessageType) {
	m.messageSendErrors.WithLabelValues(string(messageType)).Inc()
}

func (m *metrics) recordWebsocketMessageReceiveError(messageType MessageType) {
	m.messageReceiveErrors.WithLabelValues(string(messageType)).Inc()
}

func (m *metrics) recordQueueJoin(timeFormat TimeFormat, queueSize int) {
	label := timeFormat.String()
	m.queueJoinsTotal.WithLabelValues(label).Inc()
	m.queueSize.WithLabelValues(label).Set(float64(queueSize))
}

func (m *metrics) recordQueueLeave(timeFormat TimeFormat, queueSize int) {
	label := timeFormat.String()
	m.queueLeavesTotal.Inc()
	m.queueSize.WithLabelValues(label).Set(float64(queueSize))
}

func (m *metrics) recordMatchStarted(timeFormat TimeFormat) {
	m.activeMatches.Inc()
	m.matchesStarted.WithLabelValues(timeFormat.String()).Inc()
}

func (m *metrics) recordMatchFinished(result *chess.Result) {
	m.activeMatches.Dec()
	m.matchesFinished.Inc()
	if result != nil {
		m.resultsTotal.WithLabelValues(string(result.Outcome), string(result.Reason)).Inc()
	}
}

func (m *metrics) recordMatchEvent(eventType eventType) {
	m.eventsTotal.WithLabelValues(string(eventType)).Inc()
}

func (m *metrics) recordMatchEventError(eventType eventType) {
	m.eventErrorsTotal.WithLabelValues(string(eventType)).Inc()
}
