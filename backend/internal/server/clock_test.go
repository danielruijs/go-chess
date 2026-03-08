package server

import (
	"go-chess/internal/chess"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMatchClock(t *testing.T) {
	initialTime := 3 * time.Minute
	increment := 2 * time.Second
	clock := NewMatchClock(TimeFormat{initial: initialTime, increment: increment})

	clock.Start()
	assert.True(t, clock.IsRunning())
	assert.Equal(t, chess.White, clock.activeColor)

	// White thinks for 10 seconds
	whiteThinkTime := 10 * time.Second
	clock.lastUpdated = time.Now().Add(-whiteThinkTime)
	loser, err := clock.BeforeMove()
	assert.Nil(t, err)
	assert.Nil(t, loser)

	whiteAfterThink := clock.remaining[chess.White]
	assert.InDelta(t, (initialTime - whiteThinkTime).Milliseconds(), whiteAfterThink.Milliseconds(), 1)

	// White makes a move
	err = clock.AfterMove()
	assert.Nil(t, err)
	assert.Equal(t, chess.Black, clock.activeColor)

	whiteAfterMove := clock.remaining[chess.White]
	assert.InDelta(t, (initialTime - whiteThinkTime + increment).Milliseconds(), whiteAfterMove.Milliseconds(), 1)

	// Black thinks for 5 seconds
	blackThinkTime := 5 * time.Second
	clock.lastUpdated = time.Now().Add(-blackThinkTime / 2)

	// Take snapshot in the middle of think
	snapshot := clock.Snapshot()
	assert.InDelta(t, whiteAfterMove.Milliseconds(), snapshot.WhiteTimeMs, 1)
	assert.InDelta(t, (initialTime - blackThinkTime/2).Milliseconds(), snapshot.BlackTimeMs, 1)
	assert.Equal(t, increment.Milliseconds(), snapshot.IncrementMs)

	clock.lastUpdated = time.Now().Add(-blackThinkTime / 2)
	loser, err = clock.BeforeMove()
	assert.Nil(t, err)
	assert.Nil(t, loser)

	blackAfterThink := clock.remaining[chess.Black]
	assert.InDelta(t, (initialTime - blackThinkTime).Milliseconds(), blackAfterThink.Milliseconds(), 1)

	// Black makes a move
	err = clock.AfterMove()
	assert.Nil(t, err)
	assert.Equal(t, chess.White, clock.activeColor)

	blackAfterMove := clock.remaining[chess.Black]
	assert.InDelta(t, (initialTime - blackThinkTime + increment).Milliseconds(), blackAfterMove.Milliseconds(), 1)

	// White runs out of time
	clock.lastUpdated = time.Now().Add(-initialTime)
	loser, err = clock.BeforeMove()
	assert.Nil(t, err)
	assert.Equal(t, chess.White, *loser)

	// Take snapshot again
	snapshot = clock.Snapshot()
	assert.Equal(t, int64(0), snapshot.WhiteTimeMs)
	assert.InDelta(t, blackAfterMove.Milliseconds(), snapshot.BlackTimeMs, 1)
	assert.Equal(t, increment.Milliseconds(), snapshot.IncrementMs)

	// Stop clock and verify that moves are not allowed
	clock.Stop()
	assert.False(t, clock.IsRunning())

	_, err = clock.BeforeMove()
	assert.NotNil(t, err)
	err = clock.AfterMove()
	assert.NotNil(t, err)
}
