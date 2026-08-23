package service

import ("context"; "errors"; "sync")

var ErrBoundaryCapacity16 = errors.New("capacity exceeded")
type BoundaryStore16 struct { mu sync.Mutex; capacity, used int }
func NewStore(capacity int) *BoundaryStore16 { return &BoundaryStore16{capacity: capacity} }
func (s *BoundaryStore16) Reserve(ctx context.Context, n int) error { if err := ctx.Err(); err != nil { return err }; s.mu.Lock(); s.used += n; s.mu.Unlock(); return nil }
func (s *BoundaryStore16) Used() int { s.mu.Lock(); defer s.mu.Unlock(); return s.used }
