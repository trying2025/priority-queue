# Priority Queue 使用说明 / Usage Instructions

## 简介 / Introduction

这是一个基于 Go 实现的优先级队列。它支持自定义优先级顺序（升序或降序）以及在元素弹出时的回调函数。

This is a priority queue implemented in Go. It supports custom priority order (ascending or descending) and a callback function when elements are popped.

## 安装 / Installation

```sh
go get github.com/trying2025/priority-queue
```

## 使用示例 / Usage Example

```go
package main

import (
    "fmt"
    "github.com/trying2025/priority-queue"
    "sync"
    "time"
)

func main() {
    var wg sync.WaitGroup
    pq := priority_queue.NewPriorityQueue(func(data interface{}) {
    fmt.Printf("pop item: %v\n", data)
	wg.Done()
	time.Sleep(time.Millisecond * 200)
    }, priority_queue.OrderDesc)

    count := 10
    wg.Add(count)
    go func() {
        for i := 0; i < count; i++ {
            pq.PushItem(i, uint64(i))
        }
        fmt.Println("push done")
    }()

    wg.Add(count)
    go func() {
        for i := 0; i < count; i++ {
            pq.PushItem(i, uint64(i))
        }
        fmt.Println("push done")
    }()

    wg.Wait()
    pq.Stop()
}
```