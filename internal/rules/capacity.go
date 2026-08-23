package rules

import (
	"fmt"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/common"
	"sync"
)

type Capacity struct {
	mu          sync.RWMutex
	limit, used int
}

func NewCapacity(limit int) (*Capacity, error) {
	if limit < 1 {
		return nil, common.ErrInvalid
	}
	return &Capacity{limit: limit}, nil
}
func (c *Capacity) Limit() int     { c.mu.RLock(); defer c.mu.RUnlock(); return c.limit }
func (c *Capacity) Used() int      { c.mu.RLock(); defer c.mu.RUnlock(); return c.used }
func (c *Capacity) Available() int { return c.Limit() - c.Used() }
func (c *Capacity) Reserve(n int) error {
	if n <= 0 {
		return common.ErrInvalid
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.limit-c.used < n {
		return fmt.Errorf("%w: capacity", common.ErrConflict)
	}
	c.used += n
	return nil
}
func (c *Capacity) Release(n int) error {
	if n <= 0 {
		return common.ErrInvalid
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if n > c.used {
		return common.ErrInvalid
	}
	c.used -= n
	return nil
}
func (c *Capacity) Resize(limit int) error {
	if limit < 1 {
		return common.ErrInvalid
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if limit < c.used {
		return fmt.Errorf("%w: below usage", common.ErrConflict)
	}
	c.limit = limit
	return nil
}
func (c *Capacity) Snapshot() (int, int) { c.mu.RLock(); defer c.mu.RUnlock(); return c.limit, c.used }
