package snapshot

import (
	"encoding/json"
	"fmt"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/common"
	"sync"
	"time"
)

type Header struct {
	ID, Kind  string
	Version   int
	CreatedAt time.Time
}
type Store struct {
	mu      sync.RWMutex
	items   map[string][]byte
	headers map[string]Header
}

func New() *Store { return &Store{items: map[string][]byte{}, headers: map[string]Header{}} }
func (s *Store) Put(kind, id string, value any) error {
	if kind == "" || id == "" {
		return common.ErrInvalid
	}
	b, e := json.Marshal(value)
	if e != nil {
		return e
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	h := s.headers[id]
	h.ID = id
	h.Kind = kind
	h.Version++
	h.CreatedAt = common.Now()
	s.headers[id] = h
	s.items[id] = append([]byte(nil), b...)
	return nil
}
func (s *Store) Get(id string) (Header, []byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	h, ok := s.headers[id]
	if !ok {
		return Header{}, nil, common.ErrNotFound
	}
	return h, append([]byte(nil), s.items[id]...), nil
}
func (s *Store) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.headers[id]; !ok {
		return common.ErrNotFound
	}
	delete(s.headers, id)
	delete(s.items, id)
	return nil
}
func (s *Store) Count() int             { s.mu.RLock(); defer s.mu.RUnlock(); return len(s.items) }
func Decode[T any](b []byte) (T, error) { var out T; e := json.Unmarshal(b, &out); return out, e }
func ValidateHeader(h Header) error {
	if h.ID == "" || h.Kind == "" || h.Version < 1 {
		return fmt.Errorf("%w: snapshot header", common.ErrInvalid)
	}
	return nil
}
