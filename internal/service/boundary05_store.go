package service

import ("context"; "sync")

type BoundaryStore05 struct { mu sync.Mutex; values map[string]string; events []string }
func NewStore() *BoundaryStore05 { return &BoundaryStore05{values: make(map[string]string)} }
func (s *BoundaryStore05) Save(ctx context.Context, key string) error { if err := ctx.Err(); err != nil { return err }; s.mu.Lock(); defer s.mu.Unlock(); s.values[key] = "active"; s.events = append(s.events, key); return nil }
func (s *BoundaryStore05) Get(key string) (string, bool) { s.mu.Lock(); defer s.mu.Unlock(); v, ok := s.values[key]; return v, ok }
func (s *BoundaryStore05) Events() []string { s.mu.Lock(); defer s.mu.Unlock(); return s.events }
