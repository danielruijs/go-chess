package server

import (
	"context"
	"encoding/json"
	"fmt"
	"go-chess/internal/chess"
	"log"
	"sync"
	"time"
)

type Player struct {
	Name  string
	color chess.Color

	queues   map[TimeFormat]struct{}
	queuesMu sync.RWMutex

	clients   map[*Client]struct{}
	clientsMu sync.RWMutex
}

func NewPlayer() *Player {
	return &Player{
		Name:    "",
		queues:  make(map[TimeFormat]struct{}),
		clients: make(map[*Client]struct{}),
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

func (p *Player) SendInformational(msgType MessageType, data any) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("failed to marshal %s data for %s: %v", msgType, p.Name, err)
		return
	}
	msg := WSMessage{
		Type: msgType,
		Data: jsonData,
	}

	for _, client := range p.getClientsSnapshot() {
		select {
		case client.sendChan <- msg:
		default:
			log.Printf("Skipping message for %s", p.Name)
		}
	}
}

func (p *Player) SendCritical(msgType MessageType, data any) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal %s data for %s: %v", msgType, p.Name, err)
	}
	msg := WSMessage{
		Type: msgType,
		Data: jsonData,
	}

	var failed int
	for _, client := range p.getClientsSnapshot() {
		select {
		case client.sendChan <- msg:
		case <-ctx.Done():
			failed++
		}
	}
	if failed > 0 {
		return fmt.Errorf("failed to send %s message to %d clients for %s", msgType, failed, p.Name)
	}
	return nil
}
