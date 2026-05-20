package server

import (
	"context"
	"encoding/json"
	"errors"
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

	sendChan chan WSMessage
}

func NewPlayer() *Player {
	return &Player{
		Name:     "",
		queues:   make(map[TimeFormat]struct{}),
		sendChan: make(chan WSMessage, 100),
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

func (p *Player) GetQueues() []TimeFormat {
	p.queuesMu.RLock()
	defer p.queuesMu.RUnlock()
	queues := make([]TimeFormat, 0, len(p.queues))
	for timeFormat := range p.queues {
		queues = append(queues, timeFormat)
	}
	return queues
}

func (p *Player) GetSendChannel() chan WSMessage {
	return p.sendChan
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

	select {
	case p.sendChan <- msg:
	default:
		log.Printf("Skipping message for %s", p.Name)
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

	select {
	case p.sendChan <- msg:
		return nil
	case <-ctx.Done():
		return errors.New("timeout sending message")
	}
}
