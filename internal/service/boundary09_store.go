package service

import ("context"; "sync")

type BoundaryStore09 struct { mu sync.Mutex; values map[string]string; events []string }
func NewStore() *BoundaryStore09 { return &BoundaryStore09{values: make(map[string]string)} }
func (s *BoundaryStore09) Save(ctx context.Context, key string) error { if err := ctx.Err(); err != nil { return err }; s.mu.Lock(); defer s.mu.Unlock(); s.values[key] = "active"; s.events = append(s.events, key); return nil }
func (s *BoundaryStore09) Get(key string) (string, bool) { s.mu.Lock(); defer s.mu.Unlock(); v, ok := s.values[key]; return v, ok }
func (s *BoundaryStore09) Events() []string { s.mu.Lock(); defer s.mu.Unlock(); return s.events }
