package service

import ("context"; "errors"; "sync")

var ErrBoundaryTransition13 = errors.New("invalid transition")
type BoundaryRecord13 struct { Status string }
type BoundaryStore13 struct { mu sync.Mutex; records map[string]BoundaryRecord13 }
func NewStore() *BoundaryStore13 { return &BoundaryStore13{records: map[string]BoundaryRecord13{}} }
func (s *BoundaryStore13) Save(ctx context.Context, key string, r BoundaryRecord13) error { if err := ctx.Err(); err != nil { return err }; s.mu.Lock(); defer s.mu.Unlock(); s.records[key] = r; return nil }
func (s *BoundaryStore13) Get(key string) (BoundaryRecord13, bool) { s.mu.Lock(); defer s.mu.Unlock(); v, ok := s.records[key]; return v, ok }
