package service

import "context"
type BoundaryService17 struct { store *BoundaryStore17 }
func NewService(store *BoundaryStore17) *BoundaryService17 { return &BoundaryService17{store: store} }
func (s *BoundaryService17) Process(ctx context.Context, key string) error { return s.store.Save(ctx, key) }
