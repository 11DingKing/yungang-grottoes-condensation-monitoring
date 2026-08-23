package service

import ("context"; "errors"; "sync")

var ErrBoundaryAudit22 = errors.New("audit unavailable")
type BoundaryStore22 struct { mu sync.Mutex; values map[string]string; failAudit bool }
func NewStore() *BoundaryStore22 { return &BoundaryStore22{values: make(map[string]string)} }
func (s *BoundaryStore22) Run(ctx context.Context, key string) error { if err := ctx.Err(); err != nil { return err }; s.mu.Lock(); s.values[key] = "active"; s.mu.Unlock(); if s.failAudit { return ErrBoundaryAudit22 }; return nil }
func (s *BoundaryStore22) Get(key string) (string, bool) { s.mu.Lock(); defer s.mu.Unlock(); v, ok := s.values[key]; return v, ok }
func (s *BoundaryStore22) FailAudit(v bool) { s.mu.Lock(); defer s.mu.Unlock(); s.failAudit = v }
