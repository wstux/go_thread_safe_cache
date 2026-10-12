package lru

// Index constants to represent a null pointer in our array-based list
const nullIdx int32 = -1

// entry represents a flat, self-contained slot in our pre-allocated array.
type entry[K comparable, V any] struct {
	key       K
	value     V
	prevList  int32 // Doubly linked list pointers using array indices
	nextList  int32
	isDefined bool // Tracks whether the slot contains active data
}

// Cache implements a non-thread-safe LRU (Least Recently Used) cache
// using an internal sync.Pool to eliminate heap allocations during inserts.
type Cache[K comparable, V any] struct {
	slots    []entry[K, V]
	items    map[K]int32 // Maps user keys to flat indices in the slots array
	capacity int32
	size     int32
	head     int32
	tail     int32
	freeHead int32 // Head index of the pre-allocated free slots chain
}

// New creates and initializes a new Cache instance with the given capacity.
// If the provided capacity is less than or equal to zero, a default capacity of 10 is used.
func New[K comparable, V any](capacity int) *Cache[K, V] {
	if capacity <= 0 {
		capacity = 1
	}

	size := capacity
	slots := make([]entry[K, V], size)

	// Pre-link all slots into a single-linked Free List for fast O(1) acquisition
	for i := 0; i < size; i++ {
		slots[i] = entry[K, V]{
			prevList:  nullIdx,
			nextList:  int32(i + 1), // Points to the next free slot
			isDefined: false,
		}
	}
	slots[size-1].nextList = nullIdx // Terminate the free list chain

	return &Cache[K, V]{
		slots:    slots,
		items:    make(map[K]int32, size), // Pre-size the map to prevent runtime resizing
		capacity: int32(capacity),
		size:     0,
		head:     nullIdx,
		tail:     nullIdx,
		freeHead: 0,
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
	c.moveToHead(idx)
	return c.slots[idx].value, true
}

// Set adds a new key-value pair to the cache or updates an existing one.
// If the element already exists, its value is updated and it is moved to the front of the list.
// If the cache reaches its capacity limit, the least recently used element is evicted.
func (c *Cache[K, V]) Set(key K, value V) {
	// 1. If the key already exists, just update its value and promote it in-place
	if idx, exists := c.items[key]; exists {
		c.slots[idx].value = value
		c.moveToHead(idx)
		return
	}

	// 2. Determine target index.
	var targetIdx int32
	if c.size >= c.capacity {
		// Evict the least recently used element from the tail
		targetIdx = c.tail
		c.removeNode(targetIdx)

		// Clean up old key tracking
		delete(c.items, c.slots[targetIdx].key)
		c.slots[targetIdx].isDefined = false
		c.size--
	} else {
		// Fast O(1) retrieval from the Free List without loops
		targetIdx = c.freeHead
		c.freeHead = c.slots[targetIdx].nextList
	}

	// 3. Write data directly into the pre-allocated flat slot
	c.slots[targetIdx].key = key
	c.slots[targetIdx].value = value
	c.slots[targetIdx].isDefined = true
	c.size++

	// Track index inside the map indexer
	c.items[key] = targetIdx

	// Link into the active LRU doubly linked list
	c.addToHead(targetIdx)
}

// addToHead inserts a node at the front of the doubly linked list.
func (c *Cache[K, V]) addToHead(idx int32) {
	c.slots[idx].nextList = c.head
	c.slots[idx].prevList = nullIdx

	if c.head != nullIdx {
		c.slots[c.head].prevList = idx
	}
	c.head = idx
	if c.tail == nullIdx {
		c.tail = idx
	}
}

// removeNode detaches a node by updating its neighbors' references.
func (c *Cache[K, V]) removeNode(idx int32) {
	prev := c.slots[idx].prevList
	next := c.slots[idx].nextList

	if prev != nullIdx {
		c.slots[prev].nextList = next
	} else {
		c.head = next
	}

	if next != nullIdx {
		c.slots[next].prevList = prev
	} else {
		c.tail = prev
	}
}

// moveToHead promotes a node to the head position, marking it as the most recently used.
func (c *Cache[K, V]) moveToHead(idx int32) {
	if c.head == idx {
		return
	}
	c.removeNode(idx)
	c.addToHead(idx)
}
