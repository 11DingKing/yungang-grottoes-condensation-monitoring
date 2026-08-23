package snapshot

import (
	"context"
	"sync"
	"time"
)

type Recoverable interface{ Recover(context.Context) error }
type Manager struct {
	mu    sync.Mutex
	items []Recoverable
}

func (m *Manager) Add(v Recoverable) { m.mu.Lock(); m.items = append(m.items, v); m.mu.Unlock() }
func (m *Manager) Run(ctx context.Context) []error {
	m.mu.Lock()
	items := append([]Recoverable(nil), m.items...)
	m.mu.Unlock()
	out := []error{}
	for _, v := range items {
		if e := v.Recover(ctx); e != nil {
			out = append(out, e)
		}
	}
	return out
}

type Marker struct {
	At   time.Time
	Done bool
}

func (m *Marker) Mark()      { m.At = time.Now(); m.Done = true }
func (m Marker) Valid() bool { return m.Done && !m.At.IsZero() }
