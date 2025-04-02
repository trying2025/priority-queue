package priority_queue

import (
	"container/heap"
	"sync"
)

type OrderType int

const (
	OrderAsc  = OrderType(0)
	OrderDesc = OrderType(1)
)

type CallbackFun func(data interface{})
type LessFunc func(a, b interface{}) bool

type PriorityQueue struct {
	items    []*Item
	mu       sync.Mutex
	cond     *sync.Cond
	stopChan chan struct{}
	callback CallbackFun // 弹出时的回调函数
	compare  LessFunc
	order    OrderType
}

func NewPriorityQueue(callback CallbackFun, order OrderType) *PriorityQueue {
	pq := &PriorityQueue{
		items:    make([]*Item, 0),
		stopChan: make(chan struct{}),
		callback: callback,
		order:    order,
	}
	pq.cond = sync.NewCond(&pq.mu)
	heap.Init(pq)
	pq.StartPopping()
	return pq
}

// 实现 heap.Interface 的方法（与之前相同）
func (pq *PriorityQueue) Len() int { return len(pq.items) }
func (pq *PriorityQueue) Less(i, j int) bool {
	if pq.order == OrderDesc {
		return pq.items[i].priority > pq.items[j].priority // 最大堆
	} else {
		return pq.items[i].priority < pq.items[j].priority // 最小堆
	}
}
func (pq *PriorityQueue) Swap(i, j int) {
	pq.items[i], pq.items[j] = pq.items[j], pq.items[i]
	pq.items[i].index = i
	pq.items[j].index = j
}
func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*Item)
	item.index = len(pq.items)
	pq.items = append(pq.items, item)
}
func (pq *PriorityQueue) Pop() interface{} {
	old := pq.items
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	pq.items = old[0 : n-1]
	return item
}

// PushItem 线程安全操作封装
func (pq *PriorityQueue) PushItem(v interface{}, priority ...uint64) {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	if len(priority) > 0 {
		heap.Push(pq, &Item{
			value:    v,
			priority: priority[0],
		})
	} else {
		heap.Push(pq, &Item{
			value: v,
		})
	}
	pq.cond.Broadcast() // 通知等待的消费者
}

// StartPopping 启动持续弹出协程
func (pq *PriorityQueue) StartPopping() {
	go func() {
		for {
			select {
			case <-pq.stopChan:
				return
			default:
				pq.mu.Lock()

				// 等待直到有元素或收到停止信号
				for pq.Len() == 0 {
					pq.cond.Wait()
					select {
					case <-pq.stopChan:
						pq.mu.Unlock()
						return
					default:
					}
				}

				item := heap.Pop(pq).(*Item)
				pq.mu.Unlock() // 提前解锁，回调处理时不持有锁

				if pq.callback != nil {
					pq.callback(item.value)
				}
			}
		}
	}()
}

// Stop 停止队列操作
func (pq *PriorityQueue) Stop() {
	close(pq.stopChan)
	pq.cond.Broadcast() // 唤醒所有等待的协程
}
