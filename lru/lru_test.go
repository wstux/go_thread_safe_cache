package lru

import (
	"github.com/wstux/go_thread_safe_cache/internal/testsuite"
	"testing"
)

func getTestCaches[K comparable, V any](capacity int) []struct {
	name  string
	cache testsuite.GenericCache[K, V]
} {
	return []struct {
		name  string
		cache testsuite.GenericCache[K, V]
	}{
		{name: "Standard LRU Cache", cache: New[K, V](capacity)},
		{name: "Safe LRU Cache", cache: NewSafe[K, V](capacity*2, 2, testsuite.StandardHasher[K])},
	}
}

// TestCache_BasicOperations verifies standard Set and Get behavior for both cache types.
func TestCache_BasicOperations(t *testing.T) {
	for _, tc := range getTestCaches[int, int](2) {
		t.Run(tc.name, func(t *testing.T) {
			testsuite.RunBasicOperationsSuite(t, tc.cache)
		})
	}
}

// TestCache_Eviction verifies that the least recently used item is evicted when capacity is exceeded.
func TestCache_Eviction(t *testing.T) {
	for _, tc := range getTestCaches[int, int](2) {
		t.Run(tc.name, func(t *testing.T) {
			testsuite.RunEvictionSuite(t, tc.cache)
		})
	}
}

// TestCache_UpdateAndPromote verifies that accessing or updating an item marks it as recently used.
func TestCache_UpdateAndPromote(t *testing.T) {
	for _, tc := range getTestCaches[int, int](2) {
		t.Run(tc.name, func(t *testing.T) {
			testsuite.RunUpdateAndPromoteSuite(t, tc.cache)
		})
	}
}

// TestCache_SingleElementEdgeCases exercises branches where head and tail modify
// when there is only one element in the cache.
func TestCache_SingleElementEdgeCases(t *testing.T) {
	cache := New[string, int](1)
	testsuite.RunSingleElementEdgeCasesSuite(t, cache)
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
	cache := New[string, int](2)
	testsuite.RunRecycledNodeSequenceSuite(t, cache)
}

// TestSafeCache_CapacityDistribution verifies the splitting math across shards.
func TestSafeCache_CapacityDistribution(t *testing.T) {
	sc1 := NewSafe[string, int](10, 3, testsuite.StandardHasher[string])
	if sc1.shards[0].cache.capacity != 4 {
		t.Errorf("expected shard 0 to have capacity 4, got %d", sc1.shards[0].cache.capacity)
	}
	if sc1.shards[1].cache.capacity != 3 {
		t.Errorf("expected shard 1 to have capacity 3, got %d", sc1.shards[1].cache.capacity)
	}

	sc2 := NewSafe[string, int](2, 5, testsuite.StandardHasher[string])
	for i := 0; i < 5; i++ {
		if sc2.shards[i].cache.capacity != 1 {
			t.Errorf("expected shard %d to have fallback capacity 1, got %d", i, sc2.shards[i].cache.capacity)
		}
	}
}

// TestSafeCache_Concurrency runs multiple goroutines performing concurrent
// Set and Get operations specifically to test the SafeCache implementation.
func TestSafeCache_Concurrency(t *testing.T) {
	cache := NewSafe[int, int](128, 8, testsuite.StandardHasher[int])
	testsuite.RunConcurrencySuite(t, cache)
}

// BenchmarkLRUCache_Set measures the performance of sequential Set operations.
func BenchmarkLRUCache_Set(b *testing.B) {
	testsuite.BenchmarkSetSuite(b, New[int, int](1000))
}

// BenchmarkLRUCache_Get measures the performance of sequential Get operations.
func BenchmarkLRUCache_Get(b *testing.B) {
	testsuite.BenchmarkGetSuite(b, New[int, int](1000))
}

// BenchmarkLRUSafeCache_Set measures the sequential performance of SafeCache.Set.
func BenchmarkLRUSafeCache_Set(b *testing.B) {
	testsuite.BenchmarkSetSuite(b, NewSafe[int, int](1000, 16, testsuite.StandardHasher[int]))
}

// BenchmarkLRUSafeCache_Get measures the sequential performance of SafeCache.Get.
func BenchmarkLRUSafeCache_Get(b *testing.B) {
	testsuite.BenchmarkGetSuite(b, NewSafe[int, int](1000, 16, testsuite.StandardHasher[int]))
}

// BenchmarkLRUSafeCache_SetParallel measures the performance of SafeCache.Set under concurrent load.
func BenchmarkLRUSafeCache_SetParallel(b *testing.B) {
	testsuite.BenchmarkSetParallelSuite(b, NewSafe[int, int](1000, 16, testsuite.StandardHasher[int]))
}

// BenchmarkLRUSafeCache_GetParallel measures the performance of SafeCache.Get under concurrent load.
func BenchmarkLRUSafeCache_GetParallel(b *testing.B) {
	testsuite.BenchmarkGetParallelSuite(b, NewSafe[int, int](1000, 16, testsuite.StandardHasher[int]))
}
