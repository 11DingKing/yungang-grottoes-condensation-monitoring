package snapshot

import (
	"fmt"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/common"
	"sync"
)

type Versioned[T any] struct {
	Value   T
	Version int
}
type Versions[T any] struct {
	mu     sync.RWMutex
	values map[string]Versioned[T]
}

func NewVersions[T any]() *Versions[T] { return &Versions[T]{values: map[string]Versioned[T]{}} }
func (v *Versions[T]) Put(id string, value T, expected int) error {
	v.mu.Lock()
	defer v.mu.Unlock()
	old, ok := v.values[id]
	if ok && old.Version != expected {
		return fmt.Errorf("%w: version", common.ErrConflict)
	}
	next := old.Version + 1
	v.values[id] = Versioned[T]{value, next}
	return nil
}
func (v *Versions[T]) Get(id string) (Versioned[T], error) {
	v.mu.RLock()
	defer v.mu.RUnlock()
	out, ok := v.values[id]
	if !ok {
		return Versioned[T]{}, common.ErrNotFound
	}
	return out, nil
}
func (v *Versions[T]) Delete(id string, expected int) error {
	v.mu.Lock()
	defer v.mu.Unlock()
	old, ok := v.values[id]
	if !ok {
		return common.ErrNotFound
	}
	if old.Version != expected {
		return common.ErrConflict
	}
	delete(v.values, id)
	return nil
}
