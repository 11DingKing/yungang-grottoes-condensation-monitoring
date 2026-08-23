package service

import ("context"; "sync")

type BoundaryStore01 struct { mu sync.Mutex; values map[string]string }
func NewStore() *BoundaryStore01 { return &BoundaryStore01{values: map[string]string{}} }
func (s *BoundaryStore01) Save(ctx context.Context, key string) error { if err := ctx.Err(); err != nil { return err }; s.mu.Lock(); defer s.mu.Unlock(); s.values[key] = "active"; return nil }
func (s *BoundaryStore01) Get(key string) (string, bool) { s.mu.Lock(); defer s.mu.Unlock(); v, ok := s.values[key]; return v, ok }
