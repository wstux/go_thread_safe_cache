window.BENCHMARK_DATA = {
  "lastUpdate": 1791838460334,
  "repoUrl": "https://github.com/wstux/go_thread_safe_cache",
  "entries": {
    "Benchmark": [
      {
        "commit": {
          "author": {
            "email": "wstux1@gmail.com",
            "name": "wstux",
            "username": "wstux"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "26a3507fea4245374b2b41cdc2d415920a36168d",
          "message": "* all: #2 Merge remote-tracking branch 'wstux/bugfix/ci' into master;",
          "timestamp": "2026-10-12T19:55:28+03:00",
          "tree_id": "0c4cb575ac3b41014bdd5604762bbc721bfd048d",
          "url": "https://github.com/wstux/go_thread_safe_cache/commit/26a3507fea4245374b2b41cdc2d415920a36168d"
        },
        "date": 1791824160627,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkCache_Set",
            "value": 100.2,
            "unit": "ns/op\t      32 B/op\t       1 allocs/op",
            "extra": "12085729 times\n4 procs"
          },
          {
            "name": "BenchmarkCache_Set - ns/op",
            "value": 100.2,
            "unit": "ns/op",
            "extra": "12085729 times\n4 procs"
          },
          {
            "name": "BenchmarkCache_Set - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "12085729 times\n4 procs"
          },
          {
            "name": "BenchmarkCache_Set - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "12085729 times\n4 procs"
          },
          {
            "name": "BenchmarkCache_Get",
            "value": 11.02,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCache_Get - ns/op",
            "value": 11.02,
            "unit": "ns/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCache_Get - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCache_Get - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkSafeCache_Set",
            "value": 109.7,
            "unit": "ns/op\t      32 B/op\t       1 allocs/op",
            "extra": "10960821 times\n4 procs"
          },
          {
            "name": "BenchmarkSafeCache_Set - ns/op",
            "value": 109.7,
            "unit": "ns/op",
            "extra": "10960821 times\n4 procs"
          },
          {
            "name": "BenchmarkSafeCache_Set - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "10960821 times\n4 procs"
          },
          {
            "name": "BenchmarkSafeCache_Set - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "10960821 times\n4 procs"
          },
          {
            "name": "BenchmarkSafeCache_Get",
            "value": 20.9,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "58383906 times\n4 procs"
          },
          {
            "name": "BenchmarkSafeCache_Get - ns/op",
            "value": 20.9,
            "unit": "ns/op",
            "extra": "58383906 times\n4 procs"
          },
          {
            "name": "BenchmarkSafeCache_Get - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "58383906 times\n4 procs"
          },
          {
            "name": "BenchmarkSafeCache_Get - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "58383906 times\n4 procs"
          },
          {
            "name": "BenchmarkSafeCache_SetParallel",
            "value": 143.7,
            "unit": "ns/op\t      28 B/op\t       0 allocs/op",
            "extra": "8256404 times\n4 procs"
          },
          {
            "name": "BenchmarkSafeCache_SetParallel - ns/op",
            "value": 143.7,
            "unit": "ns/op",
            "extra": "8256404 times\n4 procs"
          },
          {
            "name": "BenchmarkSafeCache_SetParallel - B/op",
            "value": 28,
            "unit": "B/op",
            "extra": "8256404 times\n4 procs"
          },
          {
            "name": "BenchmarkSafeCache_SetParallel - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "8256404 times\n4 procs"
          },
          {
            "name": "BenchmarkSafeCache_GetParallel",
            "value": 43.83,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "27460159 times\n4 procs"
          },
          {
            "name": "BenchmarkSafeCache_GetParallel - ns/op",
            "value": 43.83,
            "unit": "ns/op",
            "extra": "27460159 times\n4 procs"
          },
          {
            "name": "BenchmarkSafeCache_GetParallel - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "27460159 times\n4 procs"
          },
          {
            "name": "BenchmarkSafeCache_GetParallel - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "27460159 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "wstux1@gmail.com",
            "name": "wstux",
            "username": "wstux"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "8666f5b4529665214b1c457694628df0ab13c560",
          "message": "* all: #3 Merge remote-tracking branch 'wstux/feature/lru_performance' into master;",
          "timestamp": "2026-10-12T23:53:49+03:00",
          "tree_id": "0286bf03cd4aeb45e3199d40c72270c574f5b4eb",
          "url": "https://github.com/wstux/go_thread_safe_cache/commit/8666f5b4529665214b1c457694628df0ab13c560"
        },
        "date": 1791838459574,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkCache_Set",
            "value": 51.46,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "22118391 times\n4 procs"
          },
          {
            "name": "BenchmarkCache_Set - ns/op",
            "value": 51.46,
            "unit": "ns/op",
            "extra": "22118391 times\n4 procs"
          },
          {
            "name": "BenchmarkCache_Set - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "22118391 times\n4 procs"
          },
          {
            "name": "BenchmarkCache_Set - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "22118391 times\n4 procs"
          },
          {
            "name": "BenchmarkCache_Get",
            "value": 9.128,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "132602853 times\n4 procs"
          },
          {
            "name": "BenchmarkCache_Get - ns/op",
            "value": 9.128,
            "unit": "ns/op",
            "extra": "132602853 times\n4 procs"
          },
          {
            "name": "BenchmarkCache_Get - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "132602853 times\n4 procs"
          },
          {
            "name": "BenchmarkCache_Get - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "132602853 times\n4 procs"
          },
          {
            "name": "BenchmarkSafeCache_Set",
            "value": 66.13,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "18370921 times\n4 procs"
          },
          {
            "name": "BenchmarkSafeCache_Set - ns/op",
            "value": 66.13,
            "unit": "ns/op",
            "extra": "18370921 times\n4 procs"
          },
          {
            "name": "BenchmarkSafeCache_Set - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "18370921 times\n4 procs"
          },
          {
            "name": "BenchmarkSafeCache_Set - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "18370921 times\n4 procs"
          },
          {
            "name": "BenchmarkSafeCache_Get",
            "value": 30.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "38834911 times\n4 procs"
          },
          {
            "name": "BenchmarkSafeCache_Get - ns/op",
            "value": 30.7,
            "unit": "ns/op",
            "extra": "38834911 times\n4 procs"
          },
          {
            "name": "BenchmarkSafeCache_Get - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "38834911 times\n4 procs"
          },
          {
            "name": "BenchmarkSafeCache_Get - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "38834911 times\n4 procs"
          },
          {
            "name": "BenchmarkSafeCache_SetParallel",
            "value": 47.88,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "25117828 times\n4 procs"
          },
          {
            "name": "BenchmarkSafeCache_SetParallel - ns/op",
            "value": 47.88,
            "unit": "ns/op",
            "extra": "25117828 times\n4 procs"
          },
          {
            "name": "BenchmarkSafeCache_SetParallel - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "25117828 times\n4 procs"
          },
          {
            "name": "BenchmarkSafeCache_SetParallel - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "25117828 times\n4 procs"
          },
          {
            "name": "BenchmarkSafeCache_GetParallel",
            "value": 26.12,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "44533515 times\n4 procs"
          },
          {
            "name": "BenchmarkSafeCache_GetParallel - ns/op",
            "value": 26.12,
            "unit": "ns/op",
            "extra": "44533515 times\n4 procs"
          },
          {
            "name": "BenchmarkSafeCache_GetParallel - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "44533515 times\n4 procs"
          },
          {
            "name": "BenchmarkSafeCache_GetParallel - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "44533515 times\n4 procs"
          }
        ]
      }
    ]
  }
}