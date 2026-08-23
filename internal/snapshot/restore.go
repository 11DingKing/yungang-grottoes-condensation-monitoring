package snapshot

import (
	"context"
	"errors"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/common"
	"sync"
)

type Restorer struct {
	mu     sync.Mutex
	steps  []string
	failed error
}

func (r *Restorer) Step(name string, fn func(context.Context) error) error {
	if name == "" || fn == nil {
		return common.ErrInvalid
	}
	r.mu.Lock()
	if r.failed != nil {
		e := r.failed
		r.mu.Unlock()
		return e
	}
	r.mu.Unlock()
	e := fn(context.Background())
	r.mu.Lock()
	defer r.mu.Unlock()
	if e != nil {
		r.failed = e
		return e
	}
	r.steps = append(r.steps, name)
	return nil
}
func (r *Restorer) Steps() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]string(nil), r.steps...)
}
func (r *Restorer) Error() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.failed == nil {
		return nil
	}
	return r.failed
}
func (r *Restorer) Reset() { r.mu.Lock(); r.failed = nil; r.steps = nil; r.mu.Unlock() }

var _ = errors.Is
