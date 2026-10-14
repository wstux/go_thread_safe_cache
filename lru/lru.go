package lru

import (
	"github.com/wstux/go_thread_safe_cache/internal/list"
)

// entry represents a flat, self-contained slot in our pre-allocated array.
type entry[K comparable, V any] struct {
	key       K
	value     V
	isDefined bool // Tracks whether the slot contains active data
}

// Cache implements a non-thread-safe LRU (Least Recently Used) cache
// using an internal sync.Pool to eliminate heap allocations during inserts.
type Cache[K comparable, V any] struct {
	slots    []entry[K, V]
	links    []list.Link // pre-allocated slice managed by internal/list
	items    map[K]int32 // Maps user keys to flat indices in the slots array
	capacity int32
	size     int32
	freeHead int32 // Head index of the pre-allocated free slots chain
	lruList  list.List
}

// New creates and initializes a new Cache instance with the given capacity.
// If the provided capacity is less than or equal to zero, a default capacity of 10 is used.
func New[K comparable, V any](capacity int) *Cache[K, V] {
	if capacity <= 0 {
		capacity = 1
	}

	slots := make([]entry[K, V], capacity)
	links := make([]list.Link, capacity)

	// Pre-link all slots into a single-linked Free List for fast O(1) acquisition
	for i := 0; i < capacity; i++ {
		links[i] = list.Link{
			Prev: list.NullIdx,
			Next: int32(i + 1), // Points to the next free slot
		}
	}
	links[capacity-1].Next = list.NullIdx // Terminate the free list chain

	return &Cache[K, V]{
		slots:    slots,
		links:    links,
		items:    make(map[K]int32, capacity), // Pre-size the map to prevent runtime resizing
		capacity: int32(capacity),
		size:     0,
		freeHead: 0,
		lruList:  list.New(),
	}
}

// Get looks up a key's value in the cache.
// It returns the value and true if the key exists.
// If the key is found, the element is moved to the front of the list as the most recently used.
func (c *Cache[K, V]) Get(key K) (V, bool) {
	idx, found := c.items[key]
	if !found {
		return *new(V), false
	}
	c.lruList.MoveToHead(idx, c.links)
	return c.slots[idx].value, true
}

// Set adds a new key-value pair to the cache or updates an existing one.
// If the element already exists, its value is updated and it is moved to the front of the list.
// If the cache reaches its capacity limit, the least recently used element is evicted.
func (c *Cache[K, V]) Set(key K, value V) {
	// 1. If the key already exists, just update its value and promote it in-place
	if idx, exists := c.items[key]; exists {
		c.slots[idx].value = value
		c.lruList.MoveToHead(idx, c.links)
		return
	}

	// 2. Determine target index.
	var targetIdx int32
	if c.size >= c.capacity {
		// Evict the least recently used element from the tail
		targetIdx = c.lruList.Tail
		c.lruList.RemoveNode(targetIdx, c.links)

		// Clean up old key tracking
		delete(c.items, c.slots[targetIdx].key)
		c.slots[targetIdx].isDefined = false
	} else {
		// Fast O(1) retrieval from the Free List without loops
		targetIdx = c.freeHead
		c.freeHead = c.links[targetIdx].Next
		c.size++
	}

	// 3. Write data directly into the pre-allocated flat slot
	c.slots[targetIdx].key = key
	c.slots[targetIdx].value = value
	c.slots[targetIdx].isDefined = true

	// Track index inside the map indexer
	c.items[key] = targetIdx

	// Link into the active LRU doubly linked list
	c.lruList.AddToHead(targetIdx, c.links)
}
