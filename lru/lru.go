package lru

// node represents a private node in the doubly linked list.
// It is used to track the access order of the cache items.
type node[K comparable, V any] struct {
	key   K
	value V
	prev  *node[K, V]
	next  *node[K, V]
}

// Cache implements a non-thread-safe LRU (Least Recently Used) cache.
// It is not protected by mutexes and is intended for single-threaded use.
type Cache[K comparable, V any] struct {
	items    map[K]*node[K, V]
	capacity int
	head     *node[K, V]
	tail     *node[K, V]
}

// New creates and initializes a new Cache instance with the given capacity.
// If the provided capacity is less than or equal to zero, a default capacity of 10 is used.
func New[K comparable, V any](capacity int) *Cache[K, V] {
	if capacity <= 0 {
		capacity = 10
	}
	return &Cache[K, V]{
		items:    make(map[K]*node[K, V]),
		capacity: capacity,
	}
}

// Get looks up a key's value in the cache.
// It returns the value and true if the key exists.
// If the key is found, the element is moved to the front of the list as the most recently used.
func (c *Cache[K, V]) Get(key K) (V, bool) {
	n, exists := c.items[key]
	if !exists {
		var zero V
		return zero, false
	}
	c.moveToHead(n)
	return n.value, true
}

// Set adds a new key-value pair to the cache or updates an existing one.
// If the element already exists, its value is updated and it is moved to the front of the list.
// If the cache reaches its capacity limit, the least recently used element is evicted.
func (c *Cache[K, V]) Set(key K, value V) {
	if n, exists := c.items[key]; exists {
		n.value = value
		c.moveToHead(n)
		return
	}

	if len(c.items) >= c.capacity {
		c.evict()
	}

	newNode := &node[K, V]{key: key, value: value}
	c.addToHead(newNode)
	c.items[key] = newNode
}

// addToHead inserts a node at the front of the doubly linked list.
func (c *Cache[K, V]) addToHead(n *node[K, V]) {
	n.next = c.head
	n.prev = nil
	if c.head != nil {
		c.head.prev = n
	}
	c.head = n
	if c.tail == nil {
		c.tail = n
	}
}

// removeNode detaches a node by updating its neighbors' references.
func (c *Cache[K, V]) removeNode(n *node[K, V]) {
	if n.prev != nil {
		n.prev.next = n.next
	} else {
		c.head = n.next
	}
	if n.next != nil {
		n.next.prev = n.prev
	} else {
		c.tail = n.prev
	}
}

// moveToHead promotes a node to the head position, marking it as the most recently used.
func (c *Cache[K, V]) moveToHead(n *node[K, V]) {
	c.removeNode(n)
	c.addToHead(n)
}

// evict removes the least recently used element from the tail of the list to free up space.
func (c *Cache[K, V]) evict() {
	if c.tail == nil {
		return
	}
	oldTail := c.tail
	c.removeNode(oldTail)
	delete(c.items, oldTail.key)
}
