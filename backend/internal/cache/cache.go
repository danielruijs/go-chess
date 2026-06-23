package cache

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type CleanupConfig[V any] struct {
	Interval    time.Duration
	ShouldEvict func(value V, lastUsed time.Time) bool
}

type Options[V any] struct {
	Cleanup *CleanupConfig[V]
}

type cacheEntry[V any] struct {
	value    V
	lastUsed atomic.Int64 // UnixNano timestamp
}

type Cache[K comparable, V any] struct {
	mu            sync.RWMutex
	items         map[K]*cacheEntry[V]
	cleanupConfig *CleanupConfig[V]
}

func New[K comparable, V any](opts Options[V]) *Cache[K, V] {
	if opts.Cleanup != nil {
		if opts.Cleanup.Interval <= 0 {
			panic(fmt.Sprintf("cache: cleanup interval must be positive, got %v", opts.Cleanup.Interval))
		}
		if opts.Cleanup.ShouldEvict == nil {
			panic("cache: ShouldEvict must not be nil")
		}
	}
	return &Cache[K, V]{
		items:         make(map[K]*cacheEntry[V]),
		cleanupConfig: opts.Cleanup,
	}
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, ok := c.items[key]
	if !ok {
		var zero V
		return zero, false
	}

	entry.lastUsed.Store(time.Now().UnixNano())
	return entry.value, true
}

func (c *Cache[K, V]) GetOrCreate(key K, creator func() V) V {
	c.mu.Lock()
	defer c.mu.Unlock()

	if entry, ok := c.items[key]; ok {
		entry.lastUsed.Store(time.Now().UnixNano())
		return entry.value
	}

	val := creator()
	entry := &cacheEntry[V]{
		value: val,
	}
	entry.lastUsed.Store(time.Now().UnixNano())
	c.items[key] = entry
	return val
}

func (c *Cache[K, V]) Set(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry := &cacheEntry[V]{
		value: value,
	}
	entry.lastUsed.Store(time.Now().UnixNano())
	c.items[key] = entry
}

func (c *Cache[K, V]) Delete(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
}

func (c *Cache[K, V]) StartCleanup(ctx context.Context) {
	if c.cleanupConfig == nil {
		return
	}

	go func() {
		ticker := time.NewTicker(c.cleanupConfig.Interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				c.cleanup()
			}
		}
	}()
}

func (c *Cache[K, V]) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for k, entry := range c.items {
		lastUsed := time.Unix(0, entry.lastUsed.Load())
		if c.cleanupConfig.ShouldEvict(entry.value, lastUsed) {
			delete(c.items, k)
		}
	}
}
