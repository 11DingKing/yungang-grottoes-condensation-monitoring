package service

import ("context"; "sync")

type BoundaryStore26 struct { mu sync.Mutex; values map[string]string }
func NewStore() *BoundaryStore26 { return &BoundaryStore26{values: map[string]string{}} }
func (s *BoundaryStore26) Restore() { s.mu.Lock(); defer s.mu.Unlock(); s.values = nil }
func (s *BoundaryStore26) Put(ctx context.Context, key, value string) error { if err := ctx.Err(); err != nil { return err }; s.mu.Lock(); defer s.mu.Unlock(); s.values[key] = value; return nil }
func (s *BoundaryStore26) Get(key string) (string, bool) { s.mu.Lock(); defer s.mu.Unlock(); v, ok := s.values[key]; return v, ok }
