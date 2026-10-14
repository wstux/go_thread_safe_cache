package ttl

import (
	"github.com/wstux/go_thread_safe_cache/internal/list"
	"time"
)

// entry represents a flat, self-contained slot in our pre-allocated array.
type entry[K comparable, V any] struct {
	key       K
	value     V
	expiresAt int64 // Expiration timestamp in UnixNano format (0 strictly means slot is empty)
}

// isDefined checks whether the slot currently holds active user data.
func (e *entry[K, V]) isDefined() bool {
	return e.expiresAt != 0
}

// isExpired checks whether the item has passed its time-to-live deadline.
func (e *entry[K, V]) isExpired(now int64) bool {
	return now > e.expiresAt
}

// Cache implements a highly optimized, single-array TLRU cache with a global TTL setting.
type Cache[K comparable, V any] struct {
	slots      []entry[K, V]
	links      []list.Link
	items      map[K]int32 // Maps user keys to flat indices in the slots array
	capacity   int32
	size       int32
	freeHead   int32 // Head index of the pre-allocated free slots chain
	ttlList    list.List
	defaultTTL time.Duration // Global time-to-live duration enforced on all elements
}

// New creates a completely pre-allocated, zero-allocation TTL cache instance with a fixed timeout.
func New[K comparable, V any](capacity int, defaultTTL time.Duration) *Cache[K, V] {
	if capacity <= 0 {
		capacity = 1
	}
	if defaultTTL <= 0 {
		defaultTTL = time.Minute // Safe and sane fallback if no timeout is specified
	}

	slots := make([]entry[K, V], capacity)
	links := make([]list.Link, capacity)

	// Pre-link all slots into a single-linked Free List for rapid O(1) acquisition
	for i := 0; i < capacity; i++ {
		links[i] = list.Link{
			Prev: list.NullIdx,
			Next: int32(i + 1), // Points to the next free slot
		}
	}
	links[capacity-1].Next = list.NullIdx // Terminate the free list chain
	//for i := 0; i < size; i++ {
	//	slots[i] = entry[K, V]{
	//		prevList:  nullIdx,
	//		nextList:  int32(i + 1),
	//		expiresAt: 0,            // 0 strictly marks the slot as undefined/free
	//	}
	//}

	return &Cache[K, V]{
		slots:      slots,
		links:      links,
		items:      make(map[K]int32, capacity), // Pre-size the map to prevent runtime resizing
		capacity:   int32(capacity),
		size:       0,
		freeHead:   0,
		ttlList:    list.New(),
		defaultTTL: defaultTTL,
	}
}

// Get looks up a key's value in the cache and promotes it to the front if it is alive.
func (c *Cache[K, V]) Get(key K) (V, bool) {
	idx, found := c.items[key]
	if !found {
		return *new(V), false
	}

	// Lazy eviction check using unified UnixNano timestamp
	if c.slots[idx].isExpired(time.Now().UnixNano()) {
		c.ttlList.RemoveNode(idx, c.links)
		c.evict(idx)
		return *new(V), false
	}

	c.ttlList.MoveToHead(idx, c.links)
	return c.slots[idx].value, true
}

// Set adds or updates a key-value pair in the cache.
// The expiration time is automatically calculated using the global defaultTTL.
func (c *Cache[K, V]) Set(key K, value V) {
	now := time.Now().UnixNano()

	// 1. If the key already exists, update its value, refresh expiration, and promote to head
	if idx, exists := c.items[key]; exists {
		c.slots[idx].value = value
		c.slots[idx].expiresAt = now + int64(c.defaultTTL)
		c.ttlList.MoveToHead(idx, c.links)
		return
	}

	// 2. Capacity management. If full, evict the oldest item from the tail
	var targetIdx int32
	if c.size >= c.capacity {
		targetIdx = c.ttlList.Tail
		c.ttlList.RemoveNode(targetIdx, c.links)
		delete(c.items, c.slots[targetIdx].key)
		c.slots[targetIdx].expiresAt = 0 // Return slot to raw empty state before resizing
	} else {
		// Fast O(1) retrieval from the pre-allocated Free List
		targetIdx = c.freeHead
		c.freeHead = c.links[targetIdx].Next
		c.size++
	}

	// 3. Write data directly into the pre-allocated flat slot
	c.slots[targetIdx].key = key
	c.slots[targetIdx].value = value
	c.slots[targetIdx].expiresAt = now + int64(c.defaultTTL) // Apply consistent global timeout

	// Track index inside the map indexer and push to head
	c.items[key] = targetIdx
	c.ttlList.AddToHead(targetIdx, c.links)
}

func (c *Cache[K, V]) evict(idx int32) {
	delete(c.items, c.slots[idx].key)

	// Reset data to prevent memory leaks if types contain pointers
	c.slots[idx].key = *new(K)
	c.slots[idx].value = *new(V)
	c.slots[idx].expiresAt = 0 // Returns slot to clean empty state
	c.size--

	// Return the recycled index back to the front of the Free List chain
	c.links[idx].Next = c.freeHead
	c.freeHead = idx
}
