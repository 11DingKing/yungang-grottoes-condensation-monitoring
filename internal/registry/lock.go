package registry

import (
	"context"
	"sync"
	"time"
)

type LockTable struct {
	mu    sync.Mutex
	locks map[string]chan struct{}
}

func NewLocks() *LockTable { return &LockTable{locks: map[string]chan struct{}{}} }
func (l *LockTable) Lock(ctx context.Context, key string) (func(), error) {
	l.mu.Lock()
	ch := l.locks[key]
	if ch == nil {
		ch = make(chan struct{}, 1)
		l.locks[key] = ch
	}
	l.mu.Unlock()
	select {
	case ch <- struct{}{}:
		return func() { <-ch }, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
func (l *LockTable) Try(key string, wait time.Duration) (func(), bool) {
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	unlock, e := l.Lock(ctx, key)
	return unlock, e == nil
}
