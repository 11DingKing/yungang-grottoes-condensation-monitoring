package service

import ("context"; "errors"; "sync")

var ErrBoundaryAudit24 = errors.New("audit unavailable")
type BoundaryStore24 struct { mu sync.Mutex; values map[string]string; failAudit bool }
func NewStore() *BoundaryStore24 { return &BoundaryStore24{values: make(map[string]string)} }
func (s *BoundaryStore24) Run(ctx context.Context, key string) error { if err := ctx.Err(); err != nil { return err }; s.mu.Lock(); s.values[key] = "active"; s.mu.Unlock(); if s.failAudit { return ErrBoundaryAudit24 }; return nil }
func (s *BoundaryStore24) Get(key string) (string, bool) { s.mu.Lock(); defer s.mu.Unlock(); v, ok := s.values[key]; return v, ok }
func (s *BoundaryStore24) FailAudit(v bool) { s.mu.Lock(); defer s.mu.Unlock(); s.failAudit = v }
