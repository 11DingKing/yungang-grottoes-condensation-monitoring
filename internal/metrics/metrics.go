package metrics

import (
	"sync"
	"time"
)

type Counter struct {
	mu     sync.RWMutex
	values map[string]int64
}

func NewCounter() *Counter                  { return &Counter{values: map[string]int64{}} }
func (c *Counter) Add(name string, n int64) { c.mu.Lock(); c.values[name] += n; c.mu.Unlock() }
func (c *Counter) Get(name string) int64    { c.mu.RLock(); defer c.mu.RUnlock(); return c.values[name] }
func (c *Counter) Snapshot() map[string]int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	m := make(map[string]int64, len(c.values))
	for k, v := range c.values {
		m[k] = v
	}
	return m
}

type Timer struct{ started time.Time }

func Start() Timer                     { return Timer{time.Now()} }
func (t Timer) Elapsed() time.Duration { return time.Since(t.started) }
