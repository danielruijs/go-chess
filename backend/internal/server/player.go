package server

import (
	"go-chess/internal/chess"
	"sync"
)

type Player struct {
	Name     string
	SendChan chan WSMessage
	queues   map[TimeFormat]struct{}
	queuesMu sync.RWMutex
	Color    chess.Color
}

func NewPlayer(name string) *Player {
	return &Player{
		Name:     name,
		SendChan: make(chan WSMessage, 100),
		queues:   make(map[TimeFormat]struct{}),
	}
}

func (p *Player) JoinQueue(timeFormat TimeFormat) {
	p.queuesMu.Lock()
	defer p.queuesMu.Unlock()
	p.queues[timeFormat] = struct{}{}
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

// TODO: add methods for sending messages etc.
