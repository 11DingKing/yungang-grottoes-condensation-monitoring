package service

import ("context"; "errors"; "sync")

var ErrBoundaryCapacity02 = errors.New("capacity exceeded")
type BoundaryStore02 struct { mu sync.Mutex; capacity, used int }
func NewStore(capacity int) *BoundaryStore02 { return &BoundaryStore02{capacity: capacity} }
func (s *BoundaryStore02) Reserve(ctx context.Context, n int) error { if err := ctx.Err(); err != nil { return err }; s.mu.Lock(); s.used += n; s.mu.Unlock(); return nil }
func (s *BoundaryStore02) Used() int { s.mu.Lock(); defer s.mu.Unlock(); return s.used }
