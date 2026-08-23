package snapshot

import (
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/common"
	"sort"
	"strings"
	"sync"
)

type IndexItem struct {
	ID, Kind, Zone, Status string
	Tags                   []string
	Version                int
}
type Index struct {
	mu    sync.RWMutex
	items map[string]IndexItem
}

func NewIndex() *Index { return &Index{items: map[string]IndexItem{}} }
func (i *Index) Put(v IndexItem) error {
	if v.ID == "" || v.Kind == "" {
		return common.ErrInvalid
	}
	i.mu.Lock()
	defer i.mu.Unlock()
	if old, ok := i.items[v.ID]; ok && v.Version != old.Version+1 {
		return common.ErrConflict
	}
	if v.Version == 0 {
		v.Version = 1
	}
	v.Tags = normalize(v.Tags)
	i.items[v.ID] = v
	return nil
}
func (i *Index) Get(id string) (IndexItem, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()
	v, ok := i.items[id]
	if !ok {
		return IndexItem{}, common.ErrNotFound
	}
	v.Tags = append([]string(nil), v.Tags...)
	return v, nil
}
func (i *Index) Delete(id string) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	if _, ok := i.items[id]; !ok {
		return common.ErrNotFound
	}
	delete(i.items, id)
	return nil
}
func (i *Index) Search(kind, zone, status, term string) []IndexItem {
	i.mu.RLock()
	defer i.mu.RUnlock()
	term = strings.ToLower(term)
	out := []IndexItem{}
	for _, v := range i.items {
		if kind != "" && v.Kind != kind || zone != "" && v.Zone != zone || status != "" && v.Status != status {
			continue
		}
		if term != "" && !strings.Contains(strings.ToLower(v.ID), term) {
			continue
		}
		v.Tags = append([]string(nil), v.Tags...)
		out = append(out, v)
	}
	sort.Slice(out, func(a, b int) bool { return out[a].ID < out[b].ID })
	return out
}
func (i *Index) Count() int { i.mu.RLock(); defer i.mu.RUnlock(); return len(i.items) }
func normalize(values []string) []string {
	seen := map[string]bool{}
	out := []string{}
	for _, v := range values {
		v = strings.ToLower(strings.TrimSpace(v))
		if v != "" && !seen[v] {
			seen[v] = true
			out = append(out, v)
		}
	}
	sort.Strings(out)
	return out
}
