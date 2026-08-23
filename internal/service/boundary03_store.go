package service

import ("context"; "errors"; "sync")

var ErrBoundaryAudit03 = errors.New("audit unavailable")
type BoundaryStore03 struct { mu sync.Mutex; values map[string]string; failAudit bool }
func NewStore() *BoundaryStore03 { return &BoundaryStore03{values: make(map[string]string)} }
func (s *BoundaryStore03) Run(ctx context.Context, key string) error { if err := ctx.Err(); err != nil { return err }; s.mu.Lock(); s.values[key] = "active"; s.mu.Unlock(); if s.failAudit { return ErrBoundaryAudit03 }; return nil }
func (s *BoundaryStore03) Get(key string) (string, bool) { s.mu.Lock(); defer s.mu.Unlock(); v, ok := s.values[key]; return v, ok }
func (s *BoundaryStore03) FailAudit(v bool) { s.mu.Lock(); defer s.mu.Unlock(); s.failAudit = v }
