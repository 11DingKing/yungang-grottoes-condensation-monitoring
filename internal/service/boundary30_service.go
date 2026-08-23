package service

import "context"
type BoundaryService30 struct { store *BoundaryStore30 }
func NewService(store *BoundaryStore30) *BoundaryService30 { return &BoundaryService30{store: store} }
func (s *BoundaryService30) Submit(ctx context.Context, key string) error { return s.store.Save(ctx, key) }
