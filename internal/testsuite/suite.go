package testsuite

import (
	"hash/maphash"
	"sync"
	"testing"
)

var testSeed = maphash.MakeSeed()

// StandardHasher routes keys safely for concurrent sharding metrics.
func StandardHasher[K comparable](key K) uint32 {
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

// GenericCache defines the unified behavior expected from any cache implementation.
type GenericCache[K comparable, V any] interface {
	Get(key K) (V, bool)
	Set(key K, value V)
}

// RunBasicOperationsSuite verifies standard Set and Get behavior.
func RunBasicOperationsSuite(t *testing.T, cache GenericCache[int, int]) {
	cache.Set(0, 100)
	cache.Set(1, 200)

	if val, exists := cache.Get(0); !exists || val != 100 {
		t.Errorf("expected key 0 to hold 100, got %v (exists: %v)", val, exists)
	}

	if val, exists := cache.Get(1); !exists || val != 200 {
		t.Errorf("expected key 1 to hold 200, got %v (exists: %v)", val, exists)
	}

	if _, exists := cache.Get(999); exists {
		t.Errorf("expected non-existent key to be missing")
	}
}

// RunEvictionSuite verifies that the least recently used item is evicted when capacity is exceeded.
func RunEvictionSuite(t *testing.T, cache GenericCache[int, int]) {
	cache.Set(0, 1)
	cache.Set(2, 2)

	// This insert should evict key "0" since it is the oldest item at the tail.
	cache.Set(4, 3)

	if _, exists := cache.Get(0); exists {
		t.Errorf("expected key 0 to be evicted")
	}
	if _, exists := cache.Get(2); !exists {
		t.Errorf("expected key 2 to still exist")
	}
	if val, exists := cache.Get(4); !exists || val != 3 {
		t.Errorf("expected key 4 to exist, got %v", val)
	}
}

// RunUpdateAndPromoteSuite verifies that accessing or updating an item marks it as recently used.
func RunUpdateAndPromoteSuite(t *testing.T, cache GenericCache[int, int]) {
	cache.Set(0, 10)
	cache.Set(2, 20)

	// Accessing key 0 moves it to the head, making key 2 the oldest item at the tail.
	if _, exists := cache.Get(0); !exists {
		t.Fatalf("expected key 0 to exist")
	}

	cache.Set(4, 3) // This insert should now evict key "2" instead of key "0".

	if _, exists := cache.Get(2); exists {
		t.Errorf("expected key 1 to be evicted")
	}

	if _, exists := cache.Get(0); !exists {
		t.Errorf("expected key 0 to be preserved as recently used")
	}

	cache.Set(4, 30) // Updating key "4" value and promoting it.
	cache.Set(6, 4)  // This insert should now evict key "0".

	if _, exists := cache.Get(0); exists {
		t.Errorf("expected key 0 to be evicted after key 4 update")
	}
}

// RunSingleElementEdgeCasesSuite exercises branches where head and tail modify
// when there is only one element in the cache.
func RunSingleElementEdgeCasesSuite(t *testing.T, cache GenericCache[string, int]) {
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

// RunRecycledNodeSequenceSuite verifies that the cache correctly manages list
// pointers on low capacities during a tight sequence of element promotions,
// evictions with node recycling, and subsequent updates to the head element.
func RunRecycledNodeSequenceSuite(t *testing.T, cache GenericCache[string, int]) {
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

// RunConcurrencySuite validates thread-safety bounds under async load, runs
// multiple goroutines performing concurrent Set and Get operations specifically
// to test the SafeCache implementation.
func RunConcurrencySuite(t *testing.T, cache GenericCache[int, int]) {
	var wg sync.WaitGroup
	workers := 10
	opsPerWorker := 30

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

// BenchmarkSetSuite measures the performance of sequential Set operations.
func BenchmarkSetSuite(b *testing.B, cache GenericCache[int, int]) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Set(i%2000, i) // Mix of updates and evictions
	}
}

// BenchmarkGetSuite measures the performance of sequential Get operations.
func BenchmarkGetSuite(b *testing.B, cache GenericCache[int, int]) {
	for i := 0; i < 1000; i++ {
		cache.Set(i, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = cache.Get(i % 1000) // Always hits existing elements
	}
}

// BenchmarkSetParallelSuite measures the performance of Set under concurrent load.
func BenchmarkSetParallelSuite(b *testing.B, cache GenericCache[int, int]) {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			cache.Set(i%2000, i)
			i++
		}
	})
}

// BenchmarkGetParallelSuite measures the performance of Get under concurrent load.
func BenchmarkGetParallelSuite(b *testing.B, cache GenericCache[int, int]) {
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
