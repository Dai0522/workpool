package workpool

import (
	"encoding/binary"
	"testing"
	"time"
)

func TestPool(t *testing.T) {
	wp, err := NewWorkerPool(1024, nil)
	if err != nil {
		t.Errorf("TestPool NewWorkerPool error [%v]", err)
	}
	wp.Start()

	ftasks := make([]*FutureTask, 10000)
	for i := 0; i < 10000; i++ {
		ft := NewFutureTask(&SampleTask{
			ID: uint64(i),
		})
		for ac := wp.Submit(ft); ac != nil; ac = wp.Submit(ft) {
			// t.Errorf("TestPool submit [%v], actually: [%v]", nil, ac)
		}
		ftasks[i] = ft
	}

	for i, v := range ftasks {
		res, err := v.Wait(100 * time.Millisecond)
		if err != nil {
			t.Errorf("TestPool wait [%v], actually: [%v]", nil, err)
		}
		id := binary.BigEndian.Uint64(*res)
		if id != uint64(i) {
			t.Errorf("TestPool want[%v], actually: [%v]", i, id)
		}
	}
}

func BenchmarkPoolParallel(b *testing.B) {
	wp, _ := NewWorkerPool(1024, nil)
	wp.Start()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ft := NewFutureTask(new(SampleTask))
			for ac := wp.Submit(ft); ac != nil; ac = wp.Submit(ft) {
				// t.Errorf("TestPool submit [%v], actually: [%v]", nil, ac)
			}
		}
	})
}
