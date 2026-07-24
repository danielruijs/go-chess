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

	// White thinks for 10 seconds
	whiteThinkTime := 10 * time.Second
	clock.lastUpdated = time.Now().Add(-whiteThinkTime)
	loser, err := clock.BeforeMove(chess.White)
	assert.Nil(t, err)
	assert.Nil(t, loser)

	whiteAfterThink := clock.remaining[chess.White]
	assert.InDelta(t, (initialTime - whiteThinkTime).Milliseconds(), whiteAfterThink.Milliseconds(), 1)

	// White makes a move
	err = clock.AfterMove(chess.White)
	assert.Nil(t, err)

	whiteAfterMove := clock.remaining[chess.White]
	assert.InDelta(t, (initialTime - whiteThinkTime + increment).Milliseconds(), whiteAfterMove.Milliseconds(), 1)

	// Black thinks for 5 seconds
	blackThinkTime := 5 * time.Second
	clock.lastUpdated = time.Now().Add(-blackThinkTime / 2)

	// Take snapshot in the middle of think
	snapshot := clock.Snapshot(chess.Black)
	assert.InDelta(t, whiteAfterMove.Milliseconds(), snapshot.WhiteTimeMs, 1)
	assert.InDelta(t, (initialTime - blackThinkTime/2).Milliseconds(), snapshot.BlackTimeMs, 1)
	assert.Equal(t, increment.Milliseconds(), snapshot.IncrementMs)

	clock.lastUpdated = time.Now().Add(-blackThinkTime / 2)
	loser, err = clock.BeforeMove(chess.Black)
	assert.Nil(t, err)
	assert.Nil(t, loser)

	blackAfterThink := clock.remaining[chess.Black]
	assert.InDelta(t, (initialTime - blackThinkTime).Milliseconds(), blackAfterThink.Milliseconds(), 1)

	// Black makes a move
	err = clock.AfterMove(chess.Black)
	assert.Nil(t, err)

	blackAfterMove := clock.remaining[chess.Black]
	assert.InDelta(t, (initialTime - blackThinkTime + increment).Milliseconds(), blackAfterMove.Milliseconds(), 1)

	// White runs out of time
	clock.lastUpdated = time.Now().Add(-initialTime)
	loser, err = clock.BeforeMove(chess.White)
	assert.Nil(t, err)
	assert.Equal(t, chess.White, *loser)

	// Take snapshot again
	snapshot = clock.Snapshot(chess.White)
	assert.Equal(t, int64(0), snapshot.WhiteTimeMs)
	assert.InDelta(t, blackAfterMove.Milliseconds(), snapshot.BlackTimeMs, 1)
	assert.Equal(t, increment.Milliseconds(), snapshot.IncrementMs)

	// Stop clock and verify that moves are not allowed
	clock.Stop(chess.White)
	assert.False(t, clock.IsRunning())

	_, err = clock.BeforeMove(chess.White)
	assert.NotNil(t, err)
	err = clock.AfterMove(chess.White)
	assert.NotNil(t, err)
}

func TestTimeFormatValidation(t *testing.T) {
	validFormats := []TimeFormatMs{
		{InitialMs: 60000, IncrementMs: 0},      // 1+0
		{InitialMs: 180000, IncrementMs: 2000},  // 3+2
		{InitialMs: 300000, IncrementMs: 5000},  // 5+5
		{InitialMs: 600000, IncrementMs: 0},     // 10+0
		{InitialMs: 900000, IncrementMs: 10000}, // 15+10
	}

	for _, tf := range validFormats {
		assert.NoError(t, tf.Validate(), "Expected %v to be valid", tf)
		assert.NoError(t, MsToTimeFormat(tf).Validate())
	}

	invalidFormats := []struct {
		tf  TimeFormatMs
		msg string
	}{
		{TimeFormatMs{InitialMs: -1000, IncrementMs: 0}, "negative initial time"},
		{TimeFormatMs{InitialMs: 0, IncrementMs: 10000}, "zero initial time"},
		{TimeFormatMs{InitialMs: 60000, IncrementMs: -1000}, "negative increment"},
		{TimeFormatMs{InitialMs: MaxInitialTime.Milliseconds() + 1, IncrementMs: 0}, "exceeds initial limit"},
		{TimeFormatMs{InitialMs: 60000, IncrementMs: MaxIncrement.Milliseconds() + 1}, "exceeds increment limit"},
	}

	for _, tc := range invalidFormats {
		assert.Error(t, tc.tf.Validate(), "Expected error for %s (%v)", tc.msg, tc.tf)
		assert.Error(t, MsToTimeFormat(tc.tf).Validate())
	}
}
