package service

import "context"
type BoundaryService08 struct { store *BoundaryStore08 }
func NewService(store *BoundaryStore08) *BoundaryService08 { return &BoundaryService08{store: store} }
func (s *BoundaryService08) Submit(ctx context.Context, key string) error { return s.store.Save(ctx, key) }
