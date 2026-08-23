package service

import ("context"; "sync")

type BoundaryStore12 struct { mu sync.Mutex; values map[string]string }
func NewStore() *BoundaryStore12 { return &BoundaryStore12{values: map[string]string{}} }
func (s *BoundaryStore12) Save(ctx context.Context, key string) error { if err := ctx.Err(); err != nil { return err }; s.mu.Lock(); defer s.mu.Unlock(); s.values[key] = "active"; return nil }
func (s *BoundaryStore12) Get(key string) (string, bool) { s.mu.Lock(); defer s.mu.Unlock(); v, ok := s.values[key]; return v, ok }
