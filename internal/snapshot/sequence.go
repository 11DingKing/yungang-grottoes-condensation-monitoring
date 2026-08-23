package snapshot

import (
	"fmt"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/common"
	"sync"
)

type Sequence struct {
	mu   sync.Mutex
	next map[string]int64
}

func NewSequence() *Sequence { return &Sequence{next: map[string]int64{}} }
func (s *Sequence) Next(key string) (int64, error) {
	if key == "" {
		return 0, common.ErrInvalid
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.next[key]++
	return s.next[key], nil
}
func (s *Sequence) Current(key string) int64 { s.mu.Lock(); defer s.mu.Unlock(); return s.next[key] }
func (s *Sequence) Reserve(key string, n int64) (int64, error) {
	if key == "" || n < 1 {
		return 0, common.ErrInvalid
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	start := s.next[key] + 1
	s.next[key] += n
	return start, nil
}
func (s *Sequence) Require(key string, want int64) error {
	if s.Current(key) != want {
		return fmt.Errorf("%w: sequence", common.ErrConflict)
	}
	return nil
}
