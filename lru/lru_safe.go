package lru

import "sync"

// shard represents a single concurrent partition of the cache.
type shard[K comparable, V any] struct {
	mu    sync.RWMutex
	cache *Cache[K, V]
}

// SafeCache implements a thread-safe LRU (Least Recently Used) cache.
// It wraps the non-thread-safe Cache with a sync.RWMutex to allow
// concurrent reads and synchronized writes across multiple goroutines.
type SafeCache[K comparable, V any] struct {
	shards     []shard[K, V]
	shardCount uint32
	hashFn     func(key K) uint32 // Pure, type-safe allocation-free routing function
}

// NewSafe creates a sharded thread-safe LRU cache.
// - capacity: total maximum elements distributed evenly among shards.
// - shardCount: ideal number of partitions (will be rounded up to the nearest power of 2).
// - hashFn: type-safe allocation-free routing function to map a key to a uint32 integer.
func NewSafe[K comparable, V any](capacity int, shardCount int, hashFn func(key K) uint32) *SafeCache[K, V] {
	if shardCount <= 0 {
		shardCount = 32
	}
	if capacity <= 0 {
		capacity = shardCount
	}
	if hashFn == nil {
		hashFn = func(key K) uint32 { return 0 }
	}

	shards := make([]shard[K, V], shardCount)
	// Case 1: Total capacity is less than the number of shards.
	// Every shard gets a minimum capacity of 1.
	if capacity < shardCount {
		for i := 0; i < shardCount; i++ {
			shards[i] = shard[K, V]{
				cache: New[K, V](1),
			}
		}
	} else {
		// Case 2: Distribute capacity evenly and add the remainder to the first shard.
		baseCapacity := capacity / shardCount
		remainder := capacity % shardCount

		for i := 0; i < shardCount; i++ {
			currentCapacity := baseCapacity
			if i == 0 {
				currentCapacity += remainder
			}
			shards[i] = shard[K, V]{
				cache: New[K, V](currentCapacity),
			}
		}
	}

	return &SafeCache[K, V]{
		shards:     shards,
		shardCount: uint32(shardCount),
		hashFn:     hashFn,
	}
}

// Get looks up a key's value in the cache safely using a read lock.
// It returns the value and true if the key exists.
// Accessing an existing key promotes its underlying node to the head of the list.
func (sc *SafeCache[K, V]) Get(key K) (V, bool) {
	shardIdx := uint32(sc.hashFn(key)) % sc.shardCount
	s := &sc.shards[shardIdx]

	s.mu.Lock()
	defer s.mu.Unlock()

	return s.cache.Get(key)
}

// Set adds a new key-value pair to the cache or updates an existing one safely using a write lock.
// If the cache reaches its capacity limit, the least recently used element is evicted.
func (sc *SafeCache[K, V]) Set(key K, value V) {
	shardIdx := uint32(sc.hashFn(key)) % sc.shardCount
	s := &sc.shards[shardIdx]

	s.mu.Lock()
	defer s.mu.Unlock()

	s.cache.Set(key, value)
}
