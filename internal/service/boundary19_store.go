package service

import ("context"; "sync")

type BoundaryStore19 struct { mu sync.Mutex; values map[string]string }
func NewStore() *BoundaryStore19 { return &BoundaryStore19{values: map[string]string{}} }
func (s *BoundaryStore19) Restore() { s.mu.Lock(); defer s.mu.Unlock(); s.values = nil }
func (s *BoundaryStore19) Put(ctx context.Context, key, value string) error { if err := ctx.Err(); err != nil { return err }; s.mu.Lock(); defer s.mu.Unlock(); s.values[key] = value; return nil }
func (s *BoundaryStore19) Get(key string) (string, bool) { s.mu.Lock(); defer s.mu.Unlock(); v, ok := s.values[key]; return v, ok }
