package snapshot

import (
	"sync"
	"time"
)

type Event struct {
	Type, ID string
	At       time.Time
	Detail   string
}

type Log struct {
	mu    sync.RWMutex
	items []Event
}

func (l *Log) Append(e Event) { l.mu.Lock(); l.items = append(l.items, e); l.mu.Unlock() }
func (l *Log) List(kind, id string) []Event {
	l.mu.RLock()
	defer l.mu.RUnlock()
	out := []Event{}
	for _, e := range l.items {
		if (kind == "" || e.Type == kind) && (id == "" || e.ID == id) {
			out = append(out, e)
		}
	}
	return out
}
func (l *Log) Since(t time.Time) []Event {
	l.mu.RLock()
	defer l.mu.RUnlock()
	out := []Event{}
	for _, e := range l.items {
		if !e.At.Before(t) {
			out = append(out, e)
		}
	}
	return out
}
func (l *Log) Len() int { l.mu.RLock(); defer l.mu.RUnlock(); return len(l.items) }
