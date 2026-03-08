package server

import (
	"fmt"
	"go-chess/internal/chess"
	"time"
)

type TimeFormat struct {
	initial   time.Duration
	increment time.Duration
}

type MatchClock struct {
	remaining   map[chess.Color]time.Duration
	increment   time.Duration
	activeColor chess.Color
	lastUpdated time.Time
	running     bool
}

func NewMatchClock(timeFormat TimeFormat) *MatchClock {
	return &MatchClock{
		remaining: map[chess.Color]time.Duration{
			chess.White: timeFormat.initial,
			chess.Black: timeFormat.initial,
		},
		increment:   timeFormat.increment,
		activeColor: chess.White,
		running:     false,
	}
}

func (c *MatchClock) Start() {
	c.lastUpdated = time.Now()
	c.running = true
}

func (c *MatchClock) Stop() {
	if c.running {
		c.update()
		c.running = false
	}
}

func (c *MatchClock) IsRunning() bool {
	return c.running
}

// Returns the color of the player who ran out of time, or nil if no one ran out of time
func (c *MatchClock) BeforeMove() (*chess.Color, error) {
	if !c.running {
		return nil, fmt.Errorf("clock is not running")
	}
	return c.Advance(), nil
}

// Updates the clock after a move is made
func (c *MatchClock) AfterMove() error {
	if !c.running {
		return fmt.Errorf("clock is not running")
	}

	c.remaining[c.activeColor] += c.increment
	c.lastUpdated = time.Now()
	c.activeColor = chess.GetOppositeColor(c.activeColor)

	return nil
}

func (c *MatchClock) Snapshot() ClockData {
	c.update()

	return ClockData{
		WhiteTimeMs: c.remaining[chess.White].Milliseconds(),
		BlackTimeMs: c.remaining[chess.Black].Milliseconds(),
		IncrementMs: c.increment.Milliseconds(),
	}
}

// Updates the clock and return the color of the player who ran out of time, or nil if no one ran out of time
func (c *MatchClock) Advance() *chess.Color {
	c.update()
	return c.getTimeout()
}

func (c *MatchClock) update() {
	now := time.Now()
	elapsed := now.Sub(c.lastUpdated)
	if elapsed <= 0 {
		return
	}

	c.remaining[c.activeColor] -= elapsed
	if c.remaining[c.activeColor] < 0 {
		c.remaining[c.activeColor] = 0
	}
	c.lastUpdated = now
}

func (c *MatchClock) getTimeout() *chess.Color {
	if !c.running {
		return nil
	}
	for color, remaining := range c.remaining {
		if remaining <= 0 {
			return &color
		}
	}
	return nil
}
