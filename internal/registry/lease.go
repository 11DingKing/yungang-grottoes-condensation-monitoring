package registry

import (
	"fmt"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/common"
	"sync"
	"time"
)

type Lease struct {
	ID, Holder string
	ExpiresAt  time.Time
}
type Leases struct {
	mu    sync.Mutex
	items map[string]Lease
}

func NewLeases() *Leases { return &Leases{items: map[string]Lease{}} }
func (l *Leases) Acquire(id, holder string, ttl time.Duration, now time.Time) error {
	if id == "" || holder == "" || ttl <= 0 {
		return common.ErrInvalid
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if old, ok := l.items[id]; ok && now.Before(old.ExpiresAt) {
		return fmt.Errorf("%w: lease held", common.ErrConflict)
	}
	l.items[id] = Lease{id, holder, now.Add(ttl)}
	return nil
}
func (l *Leases) Renew(id, holder string, ttl time.Duration, now time.Time) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	v, ok := l.items[id]
	if !ok || v.Holder != holder || !now.Before(v.ExpiresAt) {
		return common.ErrConflict
	}
	v.ExpiresAt = now.Add(ttl)
	l.items[id] = v
	return nil
}
func (l *Leases) Release(id, holder string) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	v, ok := l.items[id]
	if !ok {
		return common.ErrNotFound
	}
	if v.Holder != holder {
		return common.ErrForbidden
	}
	delete(l.items, id)
	return nil
}
func (l *Leases) Active(id string, now time.Time) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	v, ok := l.items[id]
	return ok && now.Before(v.ExpiresAt)
}
