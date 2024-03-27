package main

import (
	"sync"
	"sync/atomic"
	"testing"
)

// Atomic operations performs better than Mutexes in terms of concurrent programs
// sync.Map make use of atomic operations, that's what make it powerful

// In Golang docs, they say: 
// The Map type is optimized for two common use cases:
// (1) when the entry for a given key is only ever written once but read many times, as in caches that only grow, in caches systems
// (2) when multiple goroutines read, write, and overwrite entries for disjoint sets of keys, which is our case.

func BenchmarkMutexAdd(b *testing.B) {
	var a int32
	var mu sync.Mutex

	for i := 0; i < b.N; i++ {
		mu.Lock()
		a++
		mu.Unlock()
	}
}

func BenchmarkAtomicAdd(b *testing.B) {
	var a atomic.Int32
	for i := 0; i < b.N; i++ {
		a.Add(1)
	}
}