package ttl

import (
	"sync"
	"time"
)

// shard represents a single concurrent partition of the TTL cache.
type shard[K comparable, V any] struct {
	mu    sync.RWMutex
	cache *Cache[K, V]
}

// SafeCache implements a high-performance, sharded thread-safe TTL (Time-To-Live) cache.
// It distributes items across multiple independent shards to eliminate lock contention under heavy loads.
type SafeCache[K comparable, V any] struct {
	shards     []shard[K, V]
	shardCount uint32
	hashFn     func(key K) uint32
}

// NewSafe creates a sharded thread-safe TTL cache instance.
// - capacity: total maximum elements distributed among all shards.
// - shardCount: exact number of partitions (partitions will share the capacity load).
// - defaultTTL: global time-to-live duration enforced on all inserted items.
// - hashFn: type-safe allocation-free routing function to map a key K to a uint32 integer.
func NewSafe[K comparable, V any](capacity int, shardCount int, defaultTTL time.Duration, hashFn func(key K) uint32) *SafeCache[K, V] {
	if capacity <= 0 {
		capacity = 32
	}
	if shardCount <= 0 {
		shardCount = 32
	}
	if defaultTTL <= 0 {
		defaultTTL = time.Minute
	}
	if hashFn == nil {
		hashFn = func(key K) uint32 { return 0 }
	}

	shards := make([]shard[K, V], shardCount)

	// Capacity distribution math matching our strict routing architecture
	if capacity < shardCount {
		for i := 0; i < shardCount; i++ {
			shards[i] = shard[K, V]{
				cache: New[K, V](1, defaultTTL),
			}
		}
	} else {
		baseCapacity := capacity / shardCount
		remainder := capacity % shardCount

		for i := 0; i < shardCount; i++ {
			currentCapacity := baseCapacity
			if i == 0 {
				currentCapacity += remainder
			}
			shards[i] = shard[K, V]{
				cache: New[K, V](currentCapacity, defaultTTL),
			}
		}
	}

	return &SafeCache[K, V]{
		shards:     shards,
		shardCount: uint32(shardCount),
		hashFn:     hashFn,
	}
}

// Get looks up a key's value safely inside its designated shard partition.
// Note: Even though it is a lookup, we must use a full write Lock because
// expired items are lazily deleted (mutated) from the underlying cache structure.
func (sc *SafeCache[K, V]) Get(key K) (V, bool) {
	shardIdx := sc.hashFn(key) % sc.shardCount
	s := &sc.shards[shardIdx]

	s.mu.Lock()
	defer s.mu.Unlock()

	return s.cache.Get(key)
}

// Set inserts or updates a key-value pair safely inside its designated shard partition.
func (sc *SafeCache[K, V]) Set(key K, value V) {
	shardIdx := sc.hashFn(key) % sc.shardCount
	s := &sc.shards[shardIdx]

	s.mu.Lock()
	defer s.mu.Unlock()

	s.cache.Set(key, value)
}
