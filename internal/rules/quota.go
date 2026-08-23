package rules

import (
	"fmt"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/common"
	"sync"
)

type Quota struct {
	mu     sync.Mutex
	values map[string]int
	limits map[string]int
}

func NewQuota(limits map[string]int) *Quota {
	cp := map[string]int{}
	for k, v := range limits {
		cp[k] = v
	}
	return &Quota{values: map[string]int{}, limits: cp}
}
func (q *Quota) Consume(key string, n int) error {
	if n <= 0 {
		return common.ErrInvalid
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	limit, ok := q.limits[key]
	if !ok {
		return fmt.Errorf("%w: unknown quota", common.ErrInvalid)
	}
	if q.values[key]+n > limit {
		return fmt.Errorf("%w: quota %s", common.ErrConflict, key)
	}
	q.values[key] += n
	return nil
}
func (q *Quota) Restore(key string, n int) error {
	if n <= 0 {
		return common.ErrInvalid
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.values[key] < n {
		return common.ErrInvalid
	}
	q.values[key] -= n
	return nil
}
func (q *Quota) Remaining(key string) int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.limits[key] - q.values[key]
}
func (q *Quota) Snapshot() map[string]int {
	q.mu.Lock()
	defer q.mu.Unlock()
	out := map[string]int{}
	for k, v := range q.values {
		out[k] = v
	}
	return out
}
