package service

import "context"
type BoundaryService03 struct { store *BoundaryStore03 }
func NewService(store *BoundaryStore03) *BoundaryService03 { return &BoundaryService03{store: store} }
func (s *BoundaryService03) Process(ctx context.Context, key string) error { return s.store.Run(ctx, key) }
