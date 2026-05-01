package server

import (
	"net/http"
	"time"

	"go-chess/internal/chess"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type metrics struct {
	registry *prometheus.Registry

	onlineConnections prometheus.Gauge
	connectionsTotal  prometheus.Counter
	connectionsClosed prometheus.Counter
	messagesReceived  *prometheus.CounterVec
	messageErrors     *prometheus.CounterVec

	queueSize        *prometheus.GaugeVec
	queueJoinsTotal  *prometheus.CounterVec
	queueLeavesTotal prometheus.Counter

	activeMatches   prometheus.Gauge
	matchesStarted  *prometheus.CounterVec
	matchesFinished prometheus.Counter
	matchDuration   prometheus.Histogram
	eventsTotal     *prometheus.CounterVec
	resultsTotal    *prometheus.CounterVec
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
		connectionsTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "go_chess_websocket_connections_total",
			Help: "Total number of websocket connections opened.",
		}),
		connectionsClosed: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "go_chess_websocket_connections_closed_total",
			Help: "Total number of websocket connections closed.",
		}),
		messagesReceived: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "go_chess_websocket_messages_received_total",
			Help: "Total websocket messages received by message type.",
		}, []string{"type"}),
		messageErrors: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "go_chess_websocket_message_errors_total",
			Help: "Total websocket message handling errors by type and category.",
		}, []string{"type", "category"}),
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
		matchDuration: prometheus.NewHistogram(prometheus.HistogramOpts{
			Name:    "go_chess_match_duration_seconds",
			Help:    "Match duration in seconds.",
			Buckets: prometheus.DefBuckets,
		}),
		eventsTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "go_chess_match_events_total",
			Help: "Total match events by type.",
		}, []string{"type"}),
		resultsTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "go_chess_match_results_total",
			Help: "Total match results by outcome and reason.",
		}, []string{"outcome", "reason"}),
	}

	registry.MustRegister(
		m.onlineConnections,
		m.connectionsTotal,
		m.connectionsClosed,
		m.messagesReceived,
		m.messageErrors,
		m.queueSize,
		m.queueJoinsTotal,
		m.queueLeavesTotal,
		m.activeMatches,
		m.matchesStarted,
		m.matchesFinished,
		m.matchDuration,
		m.eventsTotal,
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

func (m *metrics) recordWebsocketMessageReceived(messageType MessageType) {
	m.messagesReceived.WithLabelValues(string(messageType)).Inc()
}

func (m *metrics) recordWebsocketMessageError(messageType MessageType, category string) {
	m.messageErrors.WithLabelValues(string(messageType), category).Inc()
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

func (m *metrics) recordMatchFinished(duration time.Duration, result *chess.Result) {
	m.activeMatches.Dec()
	m.matchesFinished.Inc()
	m.matchDuration.Observe(duration.Seconds())
	if result != nil {
		m.resultsTotal.WithLabelValues(string(result.Outcome), string(result.Reason)).Inc()
	}
}

func (m *metrics) recordMatchEvent(eventType eventType) {
	m.eventsTotal.WithLabelValues(string(eventType)).Inc()
}
