package server

import "go-chess/internal/chess"

type Player struct {
	Name     string
	SendChan chan WSMessage
	queues   map[TimeFormat]struct{}
	Color    chess.Color
}

func (p *Player) JoinQueue(timeFormat TimeFormat) {
	p.queues[timeFormat] = struct{}{}
}

func (p *Player) LeaveQueues() {
	p.queues = make(map[TimeFormat]struct{})
}

func (p *Player) IsInQueue(timeFormat TimeFormat) bool {
	_, inQueue := p.queues[timeFormat]
	return inQueue
}

// TODO: add methods for sending messages etc.
