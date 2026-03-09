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
	lastUpdated time.Time
	running     bool
}

func NewMatchClock(timeFormat TimeFormat) *MatchClock {
	return &MatchClock{
		remaining: map[chess.Color]time.Duration{
			chess.White: timeFormat.initial,
			chess.Black: timeFormat.initial,
		},
		increment: timeFormat.increment,
		running:   false,
	}
}

func (c *MatchClock) Start() {
	c.lastUpdated = time.Now()
	c.running = true
}

func (c *MatchClock) Stop(activeColor chess.Color) {
	if c.running {
		c.update(activeColor)
		c.running = false
	}
}

func (c *MatchClock) IsRunning() bool {
	return c.running
}

// Returns the color of the player who ran out of time, or nil if no one ran out of time
func (c *MatchClock) BeforeMove(moverColor chess.Color) (*chess.Color, error) {
	if !c.running {
		return nil, fmt.Errorf("clock is not running")
	}
	return c.Advance(moverColor), nil
}

// Updates the clock after a move is made
func (c *MatchClock) AfterMove(moverColor chess.Color) error {
	if !c.running {
		return fmt.Errorf("clock is not running")
	}

	c.remaining[moverColor] += c.increment
	c.lastUpdated = time.Now()

	return nil
}

func (c *MatchClock) Snapshot(activeColor chess.Color) ClockData {
	c.update(activeColor)

	return ClockData{
		WhiteTimeMs: c.remaining[chess.White].Milliseconds(),
		BlackTimeMs: c.remaining[chess.Black].Milliseconds(),
		IncrementMs: c.increment.Milliseconds(),
	}
}

// Updates the clock and return the color of the player who ran out of time, or nil if no one ran out of time
func (c *MatchClock) Advance(activeColor chess.Color) *chess.Color {
	c.update(activeColor)
	return c.getTimeout()
}

func (c *MatchClock) update(activeColor chess.Color) {
	now := time.Now()
	elapsed := now.Sub(c.lastUpdated)
	if elapsed <= 0 {
		return
	}

	c.remaining[activeColor] -= elapsed
	if c.remaining[activeColor] < 0 {
		c.remaining[activeColor] = 0
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

func (tf TimeFormat) String() string {
	return fmt.Sprintf("%d+%d", int(tf.initial.Minutes()), int(tf.increment.Seconds()))
}

func TimeFormatToMs(tf TimeFormat) TimeFormatMs {
	return TimeFormatMs{
		InitialMs:   tf.initial.Milliseconds(),
		IncrementMs: tf.increment.Milliseconds(),
	}
}

func MsToTimeFormat(tf TimeFormatMs) TimeFormat {
	return TimeFormat{
		initial:   time.Duration(tf.InitialMs) * time.Millisecond,
		increment: time.Duration(tf.IncrementMs) * time.Millisecond,
	}
}
