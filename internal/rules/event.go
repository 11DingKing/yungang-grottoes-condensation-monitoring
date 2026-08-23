package rules

import (
	"fmt"
	"sync"
	"time"
)

type Event struct {
	Type, Entity, Actor string
	At                  time.Time
	Payload             map[string]string
}
type Bus struct {
	mu       sync.RWMutex
	handlers map[string][]func(Event)
}

func NewBus() *Bus { return &Bus{handlers: map[string][]func(Event){}} }
func (b *Bus) On(kind string, h func(Event)) {
	b.mu.Lock()
	b.handlers[kind] = append(b.handlers[kind], h)
	b.mu.Unlock()
}
func (b *Bus) Emit(e Event) error {
	if e.Type == "" || e.Entity == "" {
		return fmt.Errorf("event requires type and entity")
	}
	b.mu.RLock()
	hs := append([]func(Event){}, b.handlers[e.Type]...)
	b.mu.RUnlock()
	for _, h := range hs {
		if h != nil {
			h(e)
		}
	}
	return nil
}
func (b *Bus) Count(kind string) int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.handlers[kind])
}
