package service

import "context"
type BoundaryService11 struct { store *BoundaryStore11 }
func NewService(store *BoundaryStore11) *BoundaryService11 { return &BoundaryService11{store: store} }
func (s *BoundaryService11) Reserve(ctx context.Context, n int) error { return s.store.Reserve(ctx, n) }
