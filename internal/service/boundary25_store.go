package service

import ("context"; "errors"; "sync")

var ErrBoundaryTransition25 = errors.New("invalid transition")
type BoundaryRecord25 struct { Status string }
type BoundaryStore25 struct { mu sync.Mutex; records map[string]BoundaryRecord25 }
func NewStore() *BoundaryStore25 { return &BoundaryStore25{records: map[string]BoundaryRecord25{}} }
func (s *BoundaryStore25) Save(ctx context.Context, key string, r BoundaryRecord25) error { if err := ctx.Err(); err != nil { return err }; s.mu.Lock(); defer s.mu.Unlock(); s.records[key] = r; return nil }
func (s *BoundaryStore25) Get(key string) (BoundaryRecord25, bool) { s.mu.Lock(); defer s.mu.Unlock(); v, ok := s.records[key]; return v, ok }
