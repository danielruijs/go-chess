package server

import (
	"encoding/json"
	"go-chess/internal/auth"
	"go-chess/internal/chess"
	"log"
	"sync"
)

// PlayerKey represents a unique identifier for a player's cached session.
// It is derived from a user's username (for authenticated users) or their session ID (for anonymous users).
type PlayerKey string

func NewPlayerKey(session auth.Session) PlayerKey {
	if session.Username != "" {
		return PlayerKey("user:" + session.Username)
	}
	return PlayerKey("anon:" + string(session.ID))
}

type Player struct {
	Key         PlayerKey
	Username    string
	DisplayName string
	color       chess.Color

	queues   map[TimeFormat]struct{}
	queuesMu sync.RWMutex

	clients   map[*Client]struct{}
	clientsMu sync.RWMutex
}

func NewPlayer(key PlayerKey, username, displayName string) *Player {
	return &Player{
		Key:         key,
		Username:    username,
		DisplayName: displayName,
		queues:      make(map[TimeFormat]struct{}),
		clients:     make(map[*Client]struct{}),
	}
}

func (p *Player) SetColor(color chess.Color) {
	p.color = color
}

func (p *Player) GetColor() chess.Color {
	return p.color
}

func (p *Player) JoinQueue(timeFormat TimeFormat) {
	p.queuesMu.Lock()
	defer p.queuesMu.Unlock()
	p.queues[timeFormat] = struct{}{}
}

func (p *Player) LeaveQueue(timeFormat TimeFormat) {
	p.queuesMu.Lock()
	defer p.queuesMu.Unlock()
	delete(p.queues, timeFormat)
}

func (p *Player) LeaveQueues() {
	p.queuesMu.Lock()
	defer p.queuesMu.Unlock()
	p.queues = make(map[TimeFormat]struct{})
}

func (p *Player) IsInQueue(timeFormat TimeFormat) bool {
	p.queuesMu.RLock()
	defer p.queuesMu.RUnlock()
	_, inQueue := p.queues[timeFormat]
	return inQueue
}

func (p *Player) IsInQueues() bool {
	p.queuesMu.RLock()
	defer p.queuesMu.RUnlock()
	return len(p.queues) > 0
}

func (p *Player) HasClients() bool {
	p.clientsMu.RLock()
	defer p.clientsMu.RUnlock()
	return len(p.clients) > 0
}

func (p *Player) GetQueues() []TimeFormat {
	p.queuesMu.RLock()
	defer p.queuesMu.RUnlock()
	queues := make([]TimeFormat, 0, len(p.queues))
	for timeFormat := range p.queues {
		queues = append(queues, timeFormat)
	}
	return queues
}

func (p *Player) RegisterClient(client *Client) {
	p.clientsMu.Lock()
	defer p.clientsMu.Unlock()
	p.clients[client] = struct{}{}
}

func (p *Player) UnregisterClient(client *Client) {
	p.clientsMu.Lock()
	defer p.clientsMu.Unlock()
	delete(p.clients, client)
}

func (p *Player) getClientsSnapshot() []*Client {
	p.clientsMu.RLock()
	defer p.clientsMu.RUnlock()

	clients := make([]*Client, 0, len(p.clients))
	for client := range p.clients {
		clients = append(clients, client)
	}
	return clients
}

func (p *Player) Send(msgType MessageType, data any) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("WARN: failed to marshal %s for %s: %v", msgType, p.DisplayName, err)
		return
	}
	msg := WSMessage{
		Type: msgType,
		Data: jsonData,
	}

	for _, client := range p.getClientsSnapshot() {
		select {
		case <-client.Done:
			// Skip clients that are already closed
			continue
		case client.sendChan <- msg:
		default:
			log.Printf("WARN: send buffer full, evicting client for player %s (msg=%s)", p.DisplayName, msgType)
			client.metrics.recordWebsocketMessageSendError(msgType, "buffer_full_evicted")
			client.Close()
			p.UnregisterClient(client)
		}
	}
}
