package priority_queue

// Item 定义队列元素
type Item struct {
	value    interface{}
	priority uint64
	index    int
}
