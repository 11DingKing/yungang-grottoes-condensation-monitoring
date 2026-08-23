package service

import "context"
type BoundaryService27 struct { store *BoundaryStore27 }
func NewService(store *BoundaryStore27) *BoundaryService27 { return &BoundaryService27{store: store} }
func (s *BoundaryService27) Reserve(ctx context.Context, n int) error { return s.store.Reserve(ctx, n) }
