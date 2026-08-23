package service

import ("context"; "errors"; "sync")

var ErrBoundaryConflict14 = errors.New("version conflict")
type BoundaryStore14 struct { mu sync.Mutex; values map[string]string; events []string }
func NewStore() *BoundaryStore14 { return &BoundaryStore14{values: make(map[string]string)} }
func (s *BoundaryStore14) Save(ctx context.Context, key string) error { if err := ctx.Err(); err != nil { return err }; s.mu.Lock(); defer s.mu.Unlock(); if _, ok := s.values[key]; ok { return ErrBoundaryConflict14 }; s.values[key] = "active"; s.events = append(s.events, key); return nil }
func (s *BoundaryStore14) Get(key string) (string, bool) { s.mu.Lock(); defer s.mu.Unlock(); v, ok := s.values[key]; return v, ok }
