package workpool

import (
	"runtime"
	"sync"
	"testing"
)

var procs = runtime.GOMAXPROCS(8)

var rb = func() *ringBuffer {
	r, _ := newRingBuffer(4096)
	return r
}()

var wg sync.WaitGroup

func BenchmarkRingBufferParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			w := newWorker(1)
			rb.push(w)
		}
	})

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rb.pop()
		}
	})
}

func TestRingBuffer(t *testing.T) {
	w1 := newWorker(1)
	w2 := newWorker(2)
	w3 := newWorker(3)

	for _, unit := range []struct {
		x        *worker
		expected error
	}{
		{w1, nil},
		{w2, nil},
		{w3, nil},
	} {
		if ac := rb.push(unit.x); ac != unit.expected {
			t.Errorf("TestRingBuffer [%v], actually: [%v]", unit.expected, ac)
		}
	}

	for _, unit := range []struct {
		expected *worker
	}{
		{w1},
		{w2},
		{w3},
		{nil},
		{nil},
	} {
		if ac := rb.pop(); ac != unit.expected {
			t.Errorf("TestRingBuffer [%v], actually: [%v]", unit.expected, ac)
		}
	}

	for _, unit := range []struct {
		x   *worker
		y   *worker
		err error
	}{
		{w1, w1, nil},
		{w2, w2, nil},
		{w3, w3, nil},
	} {
		if ac1 := rb.push(unit.x); ac1 != unit.err {
			t.Errorf("TestRingBuffer [%v], actually: [%v]", unit.err, ac1)
		}
		if ac2 := rb.pop(); ac2 != unit.y {
			t.Errorf("TestRingBuffer [%v], actually: [%v]", unit.y, ac2)
		}
	}
}

func TestRingBufferParallel(t *testing.T) {
	t.Parallel()
	wg.Add(1)
	for i := 0; i <= 10000; i++ {
		if i&1 == 0 {
			if err := rb.push(newWorker(uint64(i))); err != nil {
				t.Errorf("TestRingBufferParallel push error [%v]", err)
			}
		} else {
			if w := rb.pop(); w == nil {
				t.Errorf("TestRingBufferParallel pop nil")
			}
		}
	}
	wg.Done()
}
