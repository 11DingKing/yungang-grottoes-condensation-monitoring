package operations

import (
	"fmt"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/common"
	"sync"
)

type CheckItem struct {
	ID, Label          string
	Required, Complete bool
}
type Checklist struct {
	mu    sync.RWMutex
	Items map[string]CheckItem
}

func NewChecklist(items []CheckItem) *Checklist {
	m := map[string]CheckItem{}
	for _, i := range items {
		m[i.ID] = i
	}
	return &Checklist{Items: m}
}
func (c *Checklist) Complete(id string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	i, ok := c.Items[id]
	if !ok {
		return common.ErrNotFound
	}
	i.Complete = true
	c.Items[id] = i
	return nil
}
func (c *Checklist) Pending() []CheckItem {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := []CheckItem{}
	for _, i := range c.Items {
		if !i.Complete {
			out = append(out, i)
		}
	}
	return out
}
func (c *Checklist) Ready() bool {
	for _, i := range c.Items {
		if i.Required && !i.Complete {
			return false
		}
	}
	return true
}
func (c *Checklist) RequireReady() error {
	if !c.Ready() {
		return fmt.Errorf("%w: checklist incomplete", common.ErrConflict)
	}
	return nil
}
