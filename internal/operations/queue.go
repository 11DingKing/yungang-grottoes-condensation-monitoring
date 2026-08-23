package operations

import (
	"container/heap"
	"sync"
	"time"
)

type Item struct {
	ID       string
	Priority int
	At       time.Time
	index    int
}
type queue []*Item

func (q queue) Len() int { return len(q) }
func (q queue) Less(i, j int) bool {
	if q[i].Priority == q[j].Priority {
		return q[i].At.Before(q[j].At)
	}
	return q[i].Priority > q[j].Priority
}
func (q queue) Swap(i, j int) { q[i], q[j] = q[j], q[i]; q[i].index = i; q[j].index = j }
func (q *queue) Push(x any)   { v := x.(*Item); v.index = len(*q); *q = append(*q, v) }
func (q *queue) Pop() any     { old := *q; n := len(old); v := old[n-1]; *q = old[:n-1]; return v }

type Queue struct {
	mu    sync.Mutex
	items queue
}

func NewQueue() *Queue      { q := &Queue{}; heap.Init(&q.items); return q }
func (q *Queue) Add(v Item) { q.mu.Lock(); heap.Push(&q.items, &v); q.mu.Unlock() }
func (q *Queue) Next() (Item, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.items) == 0 {
		return Item{}, false
	}
	return *heap.Pop(&q.items).(*Item), true
}
func (q *Queue) Len() int { q.mu.Lock(); defer q.mu.Unlock(); return len(q.items) }
