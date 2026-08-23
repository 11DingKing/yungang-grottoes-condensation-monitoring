package registry

import (
	"sync"
	"time"
)

type Change struct {
	Entity, ID, Action, Actor string
	At                        time.Time
	Before, After             map[string]string
}
type History struct {
	mu    sync.RWMutex
	items []Change
}

func (h *History) Append(c Change) { h.mu.Lock(); h.items = append(h.items, c); h.mu.Unlock() }
func (h *History) For(entity, id string) []Change {
	h.mu.RLock()
	defer h.mu.RUnlock()
	out := []Change{}
	for _, v := range h.items {
		if v.Entity == entity && v.ID == id {
			out = append(out, v)
		}
	}
	return out
}
func (h *History) Latest(entity, id string) (Change, bool) {
	items := h.For(entity, id)
	if len(items) == 0 {
		return Change{}, false
	}
	return items[len(items)-1], true
}
func copyMap(in map[string]string) map[string]string {
	out := map[string]string{}
	for k, v := range in {
		out[k] = v
	}
	return out
}
