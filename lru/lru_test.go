package lru

import (
	"testing"
)

// TestCache_BasicOperations verifies standard Set and Get behavior.
func TestCache_BasicOperations(t *testing.T) {
	cache := New[string, int](2)

	cache.Set("key1", 100)
	cache.Set("key2", 200)

	if val, exists := cache.Get("key1"); !exists || val != 100 {
		t.Errorf("expected key1 to hold 100, got %v (exists: %v)", val, exists)
	}

	if val, exists := cache.Get("key2"); !exists || val != 200 {
		t.Errorf("expected key2 to hold 200, got %v (exists: %v)", val, exists)
	}

	if _, exists := cache.Get("non_existent"); exists {
		t.Errorf("expected non_existent key to be missing")
	}
}

// TestCache_Eviction verifies that the least recently used item is evicted when capacity is exceeded.
func TestCache_Eviction(t *testing.T) {
	cache := New[string, int](2)

	cache.Set("key1", 1)
	cache.Set("key2", 2)

	// This insert should evict "key1" since it is the oldest item at the tail.
	cache.Set("key3", 3)

	if _, exists := cache.Get("key1"); exists {
		t.Errorf("expected key1 to be evicted from the cache")
	}

	if _, exists := cache.Get("key2"); !exists {
		t.Errorf("expected key2 to still exist in the cache")
	}

	if _, exists := cache.Get("key3"); !exists {
		t.Errorf("expected key3 to still exist in the cache")
	}
}

// TestCache_UpdateAndPromote verifies that accessing or updating an item marks it as recently used.
func TestCache_UpdateAndPromote(t *testing.T) {
	cache := New[string, int](2)

	cache.Set("key1", 1)
	cache.Set("key2", 2)

	// Accessing "key1" moves it to the head, making "key2" the oldest item.
	if _, exists := cache.Get("key1"); !exists {
		t.Fatalf("expected key1 to exist")
	}

	// This insert should now evict "key2" instead of "key1".
	cache.Set("key3", 3)

	if _, exists := cache.Get("key2"); exists {
		t.Errorf("expected key2 to be evicted")
	}

	if _, exists := cache.Get("key1"); !exists {
		t.Errorf("expected key1 to be preserved as recently used")
	}

	// Updating "key3" value and promoting it.
	cache.Set("key3", 30)

	// This insert should now evict "key1".
	cache.Set("key4", 4)

	if _, exists := cache.Get("key1"); exists {
		t.Errorf("expected key1 to be evicted after key3 update")
	}
}

// TestCache_SingleElementEdgeCases exercises branches where head and tail modify
// when there is only one element in the cache.
func TestCache_SingleElementEdgeCases(t *testing.T) {
	cache := New[string, int](1)

	// Set on empty cache (initializes head and tail to the same node)
	cache.Set("key1", 10)

	// Update the exact same single element (exercises moveToHead where node.prev and node.next are nil)
	cache.Set("key1", 20)
	if val, exists := cache.Get("key1"); !exists || val != 20 {
		t.Errorf("expected key1 to be updated to 20, got %v", val)
	}

	// Evict the single element by adding a new one
	cache.Set("key2", 30)
	if _, exists := cache.Get("key1"); exists {
		t.Errorf("expected key1 to be evicted from single-capacity cache")
	}
}

// TestCache_DefaultCapacity verifies that zero or negative capacity falls back to a default value
// and functions correctly with standard operations.
func TestCache_DefaultCapacity(t *testing.T) {
	cache := New[string, int](0)
	if cache.capacity != 10 {
		t.Errorf("expected capacity to fall back to 10, got %d", cache.capacity)
	}

	// Ensure the fallback cache can accept values normally
	cache.Set("fallback_key", 999)
	if val, exists := cache.Get("fallback_key"); !exists || val != 999 {
		t.Errorf("expected fallback cache to store and retrieve data correctly")
	}
}

// BenchmarkCache_Set measures the performance of sequential Set operations.
func BenchmarkCache_Set(b *testing.B) {
	cache := New[int, int](1000)

	// ResetTimer ignores the time spent during the setup above
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cache.Set(i%2000, i) // Mix of updates and evictions
	}
}

// BenchmarkCache_Get measures the performance of sequential Get operations.
func BenchmarkCache_Get(b *testing.B) {
	cache := New[int, int](1000)
	// Pre-fill the cache
	for i := 0; i < 1000; i++ {
		cache.Set(i, i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = cache.Get(i % 1000) // Always hits existing elements
	}
}
