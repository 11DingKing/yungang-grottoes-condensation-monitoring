package operations

import (
	"fmt"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/common"
	"sync"
	"time"
)

type Claim struct {
	ID, Owner string
	ExpiresAt time.Time
	mu        sync.Mutex
}

func (c *Claim) Acquire(owner string, ttl time.Duration, now time.Time) error {
	if owner == "" || ttl <= 0 {
		return common.ErrInvalid
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.Owner != "" && now.Before(c.ExpiresAt) {
		return fmt.Errorf("%w: claim held", common.ErrConflict)
	}
	c.Owner = owner
	c.ExpiresAt = now.Add(ttl)
	return nil
}
func (c *Claim) Extend(owner string, ttl time.Duration, now time.Time) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.Owner != owner || !now.Before(c.ExpiresAt) {
		return common.ErrConflict
	}
	c.ExpiresAt = now.Add(ttl)
	return nil
}
func (c *Claim) Release(owner string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.Owner != owner {
		return common.ErrForbidden
	}
	c.Owner = ""
	return nil
}
func (c *Claim) Held(now time.Time) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.Owner != "" && now.Before(c.ExpiresAt)
}
