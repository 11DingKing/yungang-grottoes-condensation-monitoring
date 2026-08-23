package registry

import (
	"fmt"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/common"
	"sync"
)

type Entry struct {
	ID, Name, Zone string
	Metadata       map[string]string
	Active         bool
}
type Registry struct {
	mu      sync.RWMutex
	entries map[string]Entry
}

func New() *Registry { return &Registry{entries: map[string]Entry{}} }
func (r *Registry) Put(e Entry) error {
	if e.ID == "" || e.Name == "" || e.Zone == "" {
		return common.ErrInvalid
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if old, ok := r.entries[e.ID]; ok && old.Active {
		return fmt.Errorf("%w: entry exists", common.ErrConflict)
	}
	e.Active = true
	r.entries[e.ID] = e
	return nil
}
func (r *Registry) Get(id string) (Entry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	e, ok := r.entries[id]
	if !ok {
		return Entry{}, common.ErrNotFound
	}
	return clone(e), nil
}
func (r *Registry) Deactivate(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	e, ok := r.entries[id]
	if !ok {
		return common.ErrNotFound
	}
	if !e.Active {
		return common.ErrConflict
	}
	e.Active = false
	r.entries[id] = e
	return nil
}
func (r *Registry) List(zone string) []Entry {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := []Entry{}
	for _, e := range r.entries {
		if e.Active && (zone == "" || e.Zone == zone) {
			out = append(out, clone(e))
		}
	}
	return out
}
func clone(e Entry) Entry {
	m := map[string]string{}
	for k, v := range e.Metadata {
		m[k] = v
	}
	e.Metadata = m
	return e
}
