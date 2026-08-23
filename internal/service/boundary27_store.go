package service

import ("context"; "errors"; "sync")

var ErrBoundaryCapacity27 = errors.New("capacity exceeded")
type BoundaryStore27 struct { mu sync.Mutex; capacity, used int }
func NewStore(capacity int) *BoundaryStore27 { return &BoundaryStore27{capacity: capacity} }
func (s *BoundaryStore27) Reserve(ctx context.Context, n int) error { if err := ctx.Err(); err != nil { return err }; s.mu.Lock(); s.used += n; s.mu.Unlock(); return nil }
func (s *BoundaryStore27) Used() int { s.mu.Lock(); defer s.mu.Unlock(); return s.used }
