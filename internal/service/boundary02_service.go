package service

import "context"
type BoundaryService02 struct { store *BoundaryStore02 }
func NewService(store *BoundaryStore02) *BoundaryService02 { return &BoundaryService02{store: store} }
func (s *BoundaryService02) Reserve(ctx context.Context, n int) error { return s.store.Reserve(ctx, n) }
