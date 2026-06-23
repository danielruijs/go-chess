package cache

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCache_BasicOps(t *testing.T) {
	c, err := New[string](Options[int]{})
	assert.Nil(t, err)

	// Get on empty cache
	val, ok := c.Get("foo")
	assert.False(t, ok)
	assert.Equal(t, 0, val)

	// Set and Get
	c.Set("foo", 42)
	val, ok = c.Get("foo")
	assert.True(t, ok)
	assert.Equal(t, 42, val)

	// Delete
	c.Delete("foo")
	val, ok = c.Get("foo")
	assert.False(t, ok)
}

func TestCache_GetOrCreate(t *testing.T) {
	c, err := New[string](Options[int]{})
	assert.Nil(t, err)

	// 1. GetOrCreate on non-existent key
	creatorCalls := 0
	val := c.GetOrCreate("foo", func() int {
		creatorCalls++
		return 42
	})
	assert.Equal(t, 42, val)
	assert.Equal(t, 1, creatorCalls)

	// 2. GetOrCreate on existing key (should not call creator)
	val = c.GetOrCreate("foo", func() int {
		creatorCalls++
		return 100
	})
	assert.Equal(t, 42, val)
	assert.Equal(t, 1, creatorCalls) // still 1

	// 3. Concurrent GetOrCreate on the same key
	c2, err := New[string](Options[int]{})
	assert.Nil(t, err)
	var wg sync.WaitGroup
	numGoroutines := 10
	results := make([]int, numGoroutines)
	concurrentCreatorCalls := int64(0)

	for i := range numGoroutines {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			val := c2.GetOrCreate("bar", func() int {
				atomic.AddInt64(&concurrentCreatorCalls, 1)
				time.Sleep(10 * time.Millisecond) // simulate work
				return 999
			})
			results[idx] = val
		}(i)
	}

	wg.Wait()

	// Verify all returned 999
	for _, res := range results {
		assert.Equal(t, 999, res)
	}

	// Creator should have been called exactly once
	assert.Equal(t, int64(1), concurrentCreatorCalls)
}

func TestCache_Cleanup(t *testing.T) {
	// Setup cache with evict condition: evict even numbers
	c, err := New[string](Options[int]{
		Cleanup: &CleanupConfig[int]{
			Interval: 10 * time.Millisecond,
			ShouldEvict: func(v int, lastUsed time.Time) bool {
				return v%2 == 0
			},
		},
	})
	assert.Nil(t, err)

	c.Set("odd", 1)
	c.Set("even", 2)

	ctx := t.Context()
	c.StartCleanup(ctx)

	// Wait for cleanup tick
	time.Sleep(25 * time.Millisecond)

	val, ok := c.Get("odd")
	assert.True(t, ok)
	assert.Equal(t, 1, val)

	_, ok = c.Get("even")
	assert.False(t, ok)
}

func TestCache_CleanupStopOnCancel(t *testing.T) {
	var mu sync.Mutex
	cleanupCalls := 0

	c, err := New[string](Options[int]{
		Cleanup: &CleanupConfig[int]{
			Interval: 5 * time.Millisecond,
			ShouldEvict: func(v int, lastUsed time.Time) bool {
				mu.Lock()
				cleanupCalls++
				mu.Unlock()
				return false
			},
		},
	})
	assert.Nil(t, err)

	c.Set("foo", 1)

	ctx, cancel := context.WithCancel(t.Context())
	c.StartCleanup(ctx)

	time.Sleep(20 * time.Millisecond)

	mu.Lock()
	callsBeforeCancel := cleanupCalls
	mu.Unlock()
	assert.Greater(t, callsBeforeCancel, 0)

	cancel()

	// Sleep longer, should not tick anymore
	time.Sleep(20 * time.Millisecond)

	mu.Lock()
	callsAfterCancel := cleanupCalls
	mu.Unlock()
	assert.Equal(t, callsBeforeCancel, callsAfterCancel)
}

func TestCache_CleanupTTL(t *testing.T) {
	// Evict entries unused for more than 15ms
	c, err := New[string](Options[int]{
		Cleanup: &CleanupConfig[int]{
			Interval: 15 * time.Millisecond,
			ShouldEvict: func(v int, lastUsed time.Time) bool {
				return time.Since(lastUsed) > 15*time.Millisecond
			},
		},
	})
	assert.Nil(t, err)

	c.Set("short-lived", 1)
	c.Set("long-lived", 2)

	c.items["short-lived"].lastUsed.Store(time.Now().Add(-20 * time.Millisecond).UnixNano())
	c.items["long-lived"].lastUsed.Store(time.Now().Add(-5 * time.Millisecond).UnixNano())

	c.cleanup()

	// short-lived has not been accessed and is > 15ms old. It should be evicted.
	_, ok := c.Get("short-lived")
	assert.False(t, ok)

	// long-lived was updated, so it is only 5ms old. It should still be present.
	_, ok = c.Get("long-lived")
	assert.True(t, ok)
}

func TestCache_NewErrors(t *testing.T) {
	// 1. Negative cleanup interval
	_, err := New[string](Options[int]{
		Cleanup: &CleanupConfig[int]{
			Interval: -1 * time.Second,
			ShouldEvict: func(v int, lastUsed time.Time) bool {
				return false
			},
		},
	})
	assert.Error(t, err)

	// 2. Nil ShouldEvict
	_, err = New[string](Options[int]{
		Cleanup: &CleanupConfig[int]{
			Interval:    1 * time.Second,
			ShouldEvict: nil,
		},
	})
	assert.Error(t, err)
}
