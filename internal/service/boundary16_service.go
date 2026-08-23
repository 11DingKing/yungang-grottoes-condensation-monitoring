package service

import "context"
type BoundaryService16 struct { store *BoundaryStore16 }
func NewService(store *BoundaryStore16) *BoundaryService16 { return &BoundaryService16{store: store} }
func (s *BoundaryService16) Reserve(ctx context.Context, n int) error { return s.store.Reserve(ctx, n) }
