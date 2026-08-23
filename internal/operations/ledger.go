package operations

import (
	"sync"
	"time"
)

type Balance struct {
	mu     sync.RWMutex
	values map[string]int
}

func NewBalance() *Balance { return &Balance{values: map[string]int{}} }
func (b *Balance) Credit(key string, n int) {
	if n <= 0 {
		return
	}
	b.mu.Lock()
	b.values[key] += n
	b.mu.Unlock()
}
func (b *Balance) Debit(key string, n int) bool {
	if n <= 0 {
		return false
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.values[key] < n {
		return false
	}
	b.values[key] -= n
	return true
}
func (b *Balance) Get(key string) int { b.mu.RLock(); defer b.mu.RUnlock(); return b.values[key] }
func (b *Balance) Snapshot() map[string]int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	out := map[string]int{}
	for k, v := range b.values {
		out[k] = v
	}
	return out
}

type Entry struct {
	Key   string
	Delta int
	At    time.Time
}

func ApplyEntries(b *Balance, entries []Entry) {
	for _, e := range entries {
		if e.Delta >= 0 {
			b.Credit(e.Key, e.Delta)
		} else {
			b.Debit(e.Key, -e.Delta)
		}
	}
}
