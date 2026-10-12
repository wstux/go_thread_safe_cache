package lru

import (
	"hash/maphash"
	"sync"
	"testing"
)

var testSeed = maphash.MakeSeed()

func standardHasher[K comparable](key K) uint32 {
	switch k := any(key).(type) {
	case string:
		return uint32(maphash.String(testSeed, k))
	case int:
		return uint32(k)
	case int32:
		return uint32(k)
	case int64:
		return uint32(k)
	default:
		return 0
	}
}

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
		{name: "Safe Cache", cache: NewSafe[K, V](capacity * 2, 2, standardHasher[K])},
	}
}

// TestCache_BasicOperations verifies standard Set and Get behavior for both cache types.
func TestCache_BasicOperations(t *testing.T) {
	for _, tc := range getTestCaches[int, int](2) {
		t.Run(tc.name, func(t *testing.T) {
			tc.cache.Set(0, 100)
			tc.cache.Set(1, 200)

			if val, exists := tc.cache.Get(0); !exists || val != 100 {
				t.Errorf("expected key 0 to hold 100, got %v (exists: %v)", val, exists)
			}

			if val, exists := tc.cache.Get(1); !exists || val != 200 {
				t.Errorf("expected key 1 to hold 200, got %v (exists: %v)", val, exists)
			}

			if _, exists := tc.cache.Get(999); exists {
				t.Errorf("expected non_existent key to be missing")
			}
		})
	}
}

// TestCache_Eviction verifies that the least recently used item is evicted when capacity is exceeded.
func TestCache_Eviction(t *testing.T) {
	for _, tc := range getTestCaches[int, int](2) {
		t.Run(tc.name, func(t *testing.T) {
			tc.cache.Set(0, 1)
			tc.cache.Set(2, 2)

			// This insert should evict key "0" since it is the oldest item at the tail.
			tc.cache.Set(4, 3)

			if _, exists := tc.cache.Get(0); exists {
				t.Errorf("expected key 0 to be evicted")
			}
			if _, exists := tc.cache.Get(2); !exists {
				t.Errorf("expected key 2 to still exist")
			}
			if val, exists := tc.cache.Get(4); !exists || val != 3 {
				t.Errorf("expected key 4 to exist, got %v", val)
			}
		})
	}
}

// TestCache_UpdateAndPromote verifies that accessing or updating an item marks it as recently used.
func TestCache_UpdateAndPromote(t *testing.T) {
	for _, tc := range getTestCaches[int, int](2) {
		t.Run(tc.name, func(t *testing.T) {
	        tc.cache.Set(0, 10)
	        tc.cache.Set(2, 20)

			// Accessing key 0 moves it to the head, making key 2 the oldest item at the tail.
			if _, exists := tc.cache.Get(0); !exists {
				t.Fatalf("expected key 0 to exist")
			}

			// This insert should now evict key "2" instead of key "0".
			tc.cache.Set(4, 3)

			if _, exists := tc.cache.Get(2); exists {
				t.Errorf("expected key 1 to be evicted")
			}

			if _, exists := tc.cache.Get(0); !exists {
				t.Errorf("expected key 0 to be preserved as recently used")
			}

			// Updating key "4" value and promoting it.
			tc.cache.Set(4, 30)

			// This insert should now evict key "0".
			tc.cache.Set(6, 4)

			if _, exists := tc.cache.Get(0); exists {
				t.Errorf("expected key 0 to be evicted after key 4 update")
			}
		})
	}
}

// TestCache_SingleElementEdgeCases exercises branches where head and tail modify
// when there is only one element in the cache.
func TestCache_SingleElementEdgeCases(t *testing.T) {
	cache := New[string, int](1)

	// Set on empty cache (initializes head and tail to the same node)
	cache.Set("key1", 10)

	// Update the exact same single element
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
	// Standard Cache validation
	cache := New[string, int](0)
	if cache.capacity != 1 {
		t.Errorf("expected capacity to fall back to 1, got %d", cache.capacity)
	}
}

// TestCache_RecycledNodeSequence verifies that the cache correctly manages list pointers
// on low capacities during a tight sequence of element promotions, evictions with node recycling,
// and subsequent updates to the head element.
func TestCache_RecycledNodeSequence(t *testing.T) {
	// Initialize cache with a strict capacity of 2
	cache := New[string, int](2)

	// Step 1: Fill up the cache
	cache.Set("key1", 1)
	cache.Set("key2", 2) // List state: key2 -> key1 (tail)

	// Step 2: Access key1 to promote it to the head
	_, _ = cache.Get("key1") // List state: key1 -> key2 (tail)

	// Step 3: Evict key2 by inserting key3.
	// Under the hood, key2 node is recycled to become key3.
	cache.Set("key3", 3) // List state should be: key3 -> key1 (tail)

	// Step 4: Perform an inline update on key3.
	// This exercises moveToHead branch on the current head node.
	cache.Set("key3", 30) // List state should remain: key3 -> key1 (tail)

	// Step 5: Insert key4, which must evict key1 from the tail.
	cache.Set("key4", 4) // List state should be: key4 -> key3 (tail)

	// Verification
	if _, exists := cache.Get("key1"); exists {
		t.Errorf("expected key1 to be evicted from the tail on step 5")
	}

	if val, exists := cache.Get("key3"); !exists || val != 30 {
		t.Errorf("expected key3 to exist with value 30, got %v (exists: %v)", val, exists)
	}

	if val, exists := cache.Get("key4"); !exists || val != 4 {
		t.Errorf("expected key4 to exist with value 4, got %v (exists: %v)", val, exists)
	}
}

// TestSafeCache_CapacityDistribution verifies the splitting math across shards.
func TestSafeCache_CapacityDistribution(t *testing.T) {
	sc1 := NewSafe[string, int](10, 3, standardHasher[string])
	if sc1.shards[0].cache.capacity != 4 {
		t.Errorf("expected shard 0 to have capacity 4, got %d", sc1.shards[0].cache.capacity)
	}
	if sc1.shards[1].cache.capacity != 3 {
		t.Errorf("expected shard 1 to have capacity 3, got %d", sc1.shards[1].cache.capacity)
	}

	sc2 := NewSafe[string, int](2, 5, standardHasher[string])
	for i := 0; i < 5; i++ {
		if sc2.shards[i].cache.capacity != 1 {
			t.Errorf("expected shard %d to have fallback capacity 1, got %d", i, sc2.shards[i].cache.capacity)
		}
	}
}

// TestSafeCache_Concurrency runs multiple goroutines performing concurrent
// Set and Get operations specifically to test the SafeCache implementation.
func TestSafeCache_Concurrency(t *testing.T) {
	cache := NewSafe[int, int](128, 8, standardHasher[int])
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
	cache := NewSafe[int, int](1000, 16, standardHasher[int])
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Set(i%2000, i)
	}
}

// BenchmarkSafeCache_Get measures the sequential performance of SafeCache.Get.
func BenchmarkSafeCache_Get(b *testing.B) {
	cache := NewSafe[int, int](1000, 16, standardHasher[int])
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
	cache := NewSafe[int, int](1000, 16, standardHasher[int])
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
	cache := NewSafe[int, int](1000, 16, standardHasher[int])
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
