package cache

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// CleanupConfig defines the settings for periodic eviction of cache entries.
type CleanupConfig[V any] struct {
	Interval    time.Duration
	ShouldEvict func(value V, lastUsed time.Time) bool
	OnEvicted   func(count int)
}

// Options provides configuration settings for creating a new Cache instance.
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

// New creates and returns a new Cache instance with the specified options.
// Returns an error if the options are invalid.
func New[K comparable, V any](opts Options[V]) (*Cache[K, V], error) {
	if opts.Cleanup != nil {
		if opts.Cleanup.Interval <= 0 {
			return nil, fmt.Errorf("cache: cleanup interval must be positive, got %v", opts.Cleanup.Interval)
		}
		if opts.Cleanup.ShouldEvict == nil {
			return nil, fmt.Errorf("cache: ShouldEvict must not be nil")
		}
	}
	return &Cache[K, V]{
		items:         make(map[K]*cacheEntry[V]),
		cleanupConfig: opts.Cleanup,
	}, nil
}

// Get retrieves a value from the cache and updates its last accessed time.
// Returns the value and true if found, or a zero-value and false if not.
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

// GetOrCreate retrieves an existing value from the cache or atomically computes
// and stores a new one using the creator function if the key is missing.
func (c *Cache[K, V]) GetOrCreate(key K, creator func() V) V {
	c.mu.Lock()
	defer c.mu.Unlock()

	if entry, ok := c.items[key]; ok {
		entry.lastUsed.Store(time.Now().UnixNano())
		return entry.value
	}

	val := creator()
	c.set(key, val)
	return val
}

// Set inserts or updates a value in the cache with the given key.
func (c *Cache[K, V]) Set(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.set(key, value)
}

// SetIfAbsent inserts a value in the cache if the key does not already exist.
// Returns true if the value was inserted, or false if the key already exists.
func (c *Cache[K, V]) SetIfAbsent(key K, value V) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.items[key]; ok {
		return false
	}

	c.set(key, value)
	return true
}

// Touch updates the last accessed time of a key if it exists in the cache.
func (c *Cache[K, V]) Touch(key K) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if entry, ok := c.items[key]; ok {
		entry.lastUsed.Store(time.Now().UnixNano())
	}
}

// Delete removes a value from the cache by its key.
func (c *Cache[K, V]) Delete(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
}

// StartCleanup launches a background goroutine that periodically performs evictions
// of inactive or expired cache entries based on the CleanupConfig.
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
	var evictedCount int

	c.mu.Lock()
	for k, entry := range c.items {
		lastUsed := time.Unix(0, entry.lastUsed.Load())
		if c.cleanupConfig.ShouldEvict(entry.value, lastUsed) {
			delete(c.items, k)
			evictedCount++
		}
	}
	c.mu.Unlock()

	if evictedCount > 0 && c.cleanupConfig.OnEvicted != nil {
		c.cleanupConfig.OnEvicted(evictedCount)
	}
}

// set inserts or updates a value in the cache with the given key.
// This helper must only be called while holding the write lock.
func (c *Cache[K, V]) set(key K, value V) {
	entry := &cacheEntry[V]{
		value: value,
	}
	entry.lastUsed.Store(time.Now().UnixNano())
	c.items[key] = entry
}
