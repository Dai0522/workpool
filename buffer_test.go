package workpool

import (
	"math/rand"
	"sync"
	"testing"
)

// var procs = runtime.GOMAXPROCS(8)

var rb = func() *ringBuffer {
	r, _ := newRingBuffer(4096)
	return r
}()

func BenchmarkRingBufferParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			w := newWorker(rand.Uint64())
			rb.push(w)
		}
	})

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rb.pop()
		}
	})
}

var pool sync.Pool

func BenchmarkSyncPoolParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			switch rand.Uint64() % 2 {
			case 0:

			case 1:
				pool.Get()
			}
			w := newWorker(rand.Uint64())
			pool.Put(w)
		}
	})

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			pool.Get()
		}
	})
}
