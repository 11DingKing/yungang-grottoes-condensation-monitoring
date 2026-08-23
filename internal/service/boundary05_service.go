package service

import "context"
type BoundaryService05 struct { store *BoundaryStore05 }
func NewService(store *BoundaryStore05) *BoundaryService05 { return &BoundaryService05{store: store} }
func (s *BoundaryService05) Process(ctx context.Context, key string) error { return s.store.Save(ctx, key) }
