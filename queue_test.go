package priority_queue

import (
	"sync"
	"testing"
	"time"
)

func TestQueue(t *testing.T) {
	var wg sync.WaitGroup
	pq := NewPriorityQueue(func(data interface{}) {
		t.Logf("pop item: %v", data)
		wg.Done()
	}, OrderDesc)

	count := 10
	wg.Add(count)
	go func() {
		for i := 0; i < count; i++ {
			pq.PushItem(i, uint64(i))
			time.Sleep(time.Millisecond * 200)

		}
		t.Log("push done")
	}()

	wg.Add(count)
	go func() {
		for i := 0; i < count; i++ {
			pq.PushItem(i, uint64(i))
			time.Sleep(time.Millisecond * 200)

		}
		t.Log("push done")
	}()

	wg.Wait()
	pq.Stop()
}
