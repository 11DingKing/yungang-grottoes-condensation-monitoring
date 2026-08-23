package planning

import (
	"fmt"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/common"
	"sync"
)

type Resource struct {
	ID, Kind       string
	Capacity, Used int
	mu             sync.Mutex
}

func (r *Resource) Allocate(n int) error {
	if n <= 0 {
		return common.ErrInvalid
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.Used+n > r.Capacity {
		return fmt.Errorf("%w: resource capacity", common.ErrConflict)
	}
	r.Used += n
	return nil
}
func (r *Resource) Free(n int) error {
	if n <= 0 {
		return common.ErrInvalid
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if n > r.Used {
		return common.ErrConflict
	}
	r.Used -= n
	return nil
}
func (r *Resource) Available() int       { r.mu.Lock(); defer r.mu.Unlock(); return r.Capacity - r.Used }
func (r *Resource) Snapshot() (int, int) { r.mu.Lock(); defer r.mu.Unlock(); return r.Capacity, r.Used }
func MergeResources(items []Resource) map[string]*Resource {
	out := map[string]*Resource{}
	for i := range items {
		out[items[i].ID] = &Resource{ID: items[i].ID, Kind: items[i].Kind, Capacity: items[i].Capacity, Used: items[i].Used}
	}
	return out
}
