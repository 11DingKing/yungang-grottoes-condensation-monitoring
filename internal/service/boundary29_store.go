package service

import ("context"; "errors"; "sync")

var ErrBoundaryConflict29 = errors.New("version conflict")
type BoundaryStore29 struct { mu sync.Mutex; values map[string]string; events []string }
func NewStore() *BoundaryStore29 { return &BoundaryStore29{values: make(map[string]string)} }
func (s *BoundaryStore29) Save(ctx context.Context, key string) error { if err := ctx.Err(); err != nil { return err }; s.mu.Lock(); defer s.mu.Unlock(); if _, ok := s.values[key]; ok { return ErrBoundaryConflict29 }; s.values[key] = "active"; s.events = append(s.events, key); return nil }
func (s *BoundaryStore29) Get(key string) (string, bool) { s.mu.Lock(); defer s.mu.Unlock(); v, ok := s.values[key]; return v, ok }
