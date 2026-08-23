package service

import ("context"; "sync")

type BoundaryStore20 struct { mu sync.Mutex; values map[string]string; events []string }
func NewStore() *BoundaryStore20 { return &BoundaryStore20{values: make(map[string]string)} }
func (s *BoundaryStore20) Save(ctx context.Context, key string) error { if err := ctx.Err(); err != nil { return err }; s.mu.Lock(); defer s.mu.Unlock(); s.values[key] = "active"; s.events = append(s.events, key); return nil }
func (s *BoundaryStore20) Get(key string) (string, bool) { s.mu.Lock(); defer s.mu.Unlock(); v, ok := s.values[key]; return v, ok }
func (s *BoundaryStore20) Events() []string { s.mu.Lock(); defer s.mu.Unlock(); return s.events }
