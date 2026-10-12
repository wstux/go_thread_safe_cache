package lru

import (
	"sync"
	"testing"
)

// cacheInterface defines the common behavior of both standard and safe caches.
// This allows us to run the same set of tests against both implementations.
type cacheInterface[K comparable, V any] interface {
	Get(key K) (V, bool)
	Set(key K, value V)
}

// getTestCaches returns a slice of both cache implementations initialized with the given capacity.
type testCacheInstance[K comparable, V any] struct {
	name  string
	cache cacheInterface[K, V]
}

func getTestCaches[K comparable, V any](capacity int) []testCacheInstance[K, V] {
	return []testCacheInstance[K, V]{
		{name: "Standard Cache", cache: New[K, V](capacity)},
		{name: "Safe Cache", cache: NewSafe[K, V](capacity)},
	}
}

// TestCache_BasicOperations verifies standard Set and Get behavior for both cache types.
func TestCache_BasicOperations(t *testing.T) {
	for _, tc := range getTestCaches[string, int](2) {
		t.Run(tc.name, func(t *testing.T) {
			tc.cache.Set("key1", 100)
			tc.cache.Set("key2", 200)

			if val, exists := tc.cache.Get("key1"); !exists || val != 100 {
				t.Errorf("expected key1 to hold 100, got %v (exists: %v)", val, exists)
			}

			if val, exists := tc.cache.Get("key2"); !exists || val != 200 {
				t.Errorf("expected key2 to hold 200, got %v (exists: %v)", val, exists)
			}

			if _, exists := tc.cache.Get("non_existent"); exists {
				t.Errorf("expected non_existent key to be missing")
			}
		})
	}
}

// TestCache_Eviction verifies that the least recently used item is evicted when capacity is exceeded.
func TestCache_Eviction(t *testing.T) {
	for _, tc := range getTestCaches[string, int](2) {
		t.Run(tc.name, func(t *testing.T) {
			tc.cache.Set("key1", 1)
			tc.cache.Set("key2", 2)

			// This insert should evict "key1" since it is the oldest item at the tail.
			tc.cache.Set("key3", 3)

			if _, exists := tc.cache.Get("key1"); exists {
				t.Errorf("expected key1 to be evicted from the cache")
			}

			if _, exists := tc.cache.Get("key2"); !exists {
				t.Errorf("expected key2 to still exist in the cache")
			}

			if _, exists := tc.cache.Get("key3"); !exists {
				t.Errorf("expected key3 to still exist in the cache")
			}
		})
	}
}

// TestCache_UpdateAndPromote verifies that accessing or updating an item marks it as recently used.
func TestCache_UpdateAndPromote(t *testing.T) {
	for _, tc := range getTestCaches[string, int](2) {
		t.Run(tc.name, func(t *testing.T) {
			tc.cache.Set("key1", 1)
			tc.cache.Set("key2", 2)

			// Accessing "key1" moves it to the head, making "key2" the oldest item.
			if _, exists := tc.cache.Get("key1"); !exists {
				t.Fatalf("expected key1 to exist")
			}

			// This insert should now evict "key2" instead of "key1".
			tc.cache.Set("key3", 3)

			if _, exists := tc.cache.Get("key2"); exists {
				t.Errorf("expected key2 to be evicted")
			}

			if _, exists := tc.cache.Get("key1"); !exists {
				t.Errorf("expected key1 to be preserved as recently used")
			}

			// Updating "key3" value and promoting it.
			tc.cache.Set("key3", 30)

			// This insert should now evict "key1".
			tc.cache.Set("key4", 4)

			if _, exists := tc.cache.Get("key1"); exists {
				t.Errorf("expected key1 to be evicted after key3 update")
			}
		})
	}
}

// TestCache_SingleElementEdgeCases exercises branches where head and tail modify
// when there is only one element in the cache.
func TestCache_SingleElementEdgeCases(t *testing.T) {
	for _, tc := range getTestCaches[string, int](1) {
		t.Run(tc.name, func(t *testing.T) {
			// Set on empty cache (initializes head and tail to the same node)
			tc.cache.Set("key1", 10)

			// Update the exact same single element (exercises moveToHead where node.prev and node.next are nil)
			tc.cache.Set("key1", 20)
			if val, exists := tc.cache.Get("key1"); !exists || val != 20 {
				t.Errorf("expected key1 to be updated to 20, got %v", val)
			}

			// Evict the single element by adding a new one
			tc.cache.Set("key2", 30)
			if _, exists := tc.cache.Get("key1"); exists {
				t.Errorf("expected key1 to be evicted from single-capacity cache")
			}
		})
	}
}

// TestCache_DefaultCapacity verifies that zero or negative capacity falls back to a default value
// and functions correctly with standard operations.
func TestCache_DefaultCapacity(t *testing.T) {
	// Standard Cache validation
	cache := New[string, int](0)
	if cache.capacity != 10 {
		t.Errorf("expected capacity to fall back to 10, got %d", cache.capacity)
	}

	// Safe Cache validation
	safeCache := NewSafe[string, int](0)
	if safeCache.cache.capacity != 10 {
		t.Errorf("expected safe cache underlying capacity to fall back to 10, got %d", safeCache.cache.capacity)
	}
}

// TestSafeCache_Concurrency runs multiple goroutines performing concurrent
// Set and Get operations specifically to test the SafeCache implementation.
func TestSafeCache_Concurrency(t *testing.T) {
	cache := NewSafe[int, int](50)
	var wg sync.WaitGroup
	workers := 100
	opsPerWorker := 100

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for j := 0; j < opsPerWorker; j++ {
				key := (workerID * opsPerWorker) + j
				cache.Set(key, j)
				_, _ = cache.Get(key)
			}
		}(i)
	}

	wg.Wait()
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

// BenchmarkSafeCache_Set measures the sequential performance of SafeCache.Set.
func BenchmarkSafeCache_Set(b *testing.B) {
	cache := NewSafe[int, int](1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Set(i%2000, i)
	}
}

// BenchmarkSafeCache_Get measures the sequential performance of SafeCache.Get.
func BenchmarkSafeCache_Get(b *testing.B) {
	cache := NewSafe[int, int](1000)
	for i := 0; i < 1000; i++ {
		cache.Set(i, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = cache.Get(i % 1000)
	}
}

// BenchmarkSafeCache_SetParallel measures the performance of SafeCache.Set under concurrent load.
func BenchmarkSafeCache_SetParallel(b *testing.B) {
	cache := NewSafe[int, int](1000)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			cache.Set(i%2000, i)
			i++
		}
	})
}

// BenchmarkSafeCache_GetParallel measures the performance of SafeCache.Get under concurrent load.
func BenchmarkSafeCache_GetParallel(b *testing.B) {
	cache := NewSafe[int, int](1000)
	for i := 0; i < 1000; i++ {
		cache.Set(i, i)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			_, _ = cache.Get(i % 1000)
			i++
		}
	})
}
