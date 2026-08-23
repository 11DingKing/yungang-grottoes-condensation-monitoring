package service

import ("context"; "errors"; "sync")

var ErrBoundaryCapacity11 = errors.New("capacity exceeded")
type BoundaryStore11 struct { mu sync.Mutex; capacity, used int }
func NewStore(capacity int) *BoundaryStore11 { return &BoundaryStore11{capacity: capacity} }
func (s *BoundaryStore11) Reserve(ctx context.Context, n int) error { if err := ctx.Err(); err != nil { return err }; s.mu.Lock(); s.used += n; s.mu.Unlock(); return nil }
func (s *BoundaryStore11) Used() int { s.mu.Lock(); defer s.mu.Unlock(); return s.used }
