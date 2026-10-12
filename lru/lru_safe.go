package lru

import "sync"

// SafeCache implements a thread-safe LRU (Least Recently Used) cache.
// It wraps the non-thread-safe Cache with a sync.RWMutex to allow
// concurrent reads and synchronized writes across multiple goroutines.
type SafeCache[K comparable, V any] struct {
	mu    sync.RWMutex
	cache *Cache[K, V]
}

// NewSafe creates and initializes a new SafeCache instance with the given capacity.
// If the provided capacity is less than or equal to zero, a default capacity of 10 is used.
func NewSafe[K comparable, V any](capacity int) *SafeCache[K, V] {
	return &SafeCache[K, V]{
		cache: New[K, V](capacity),
	}
}

// Get looks up a key's value in the cache safely using a read lock.
// It returns the value and true if the key exists.
// Accessing an existing key promotes its underlying node to the head of the list.
func (sc *SafeCache[K, V]) Get(key K) (V, bool) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	return sc.cache.Get(key)
}

// Set adds a new key-value pair to the cache or updates an existing one safely using a write lock.
// If the cache reaches its capacity limit, the least recently used element is evicted.
func (sc *SafeCache[K, V]) Set(key K, value V) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.cache.Set(key, value)
}
