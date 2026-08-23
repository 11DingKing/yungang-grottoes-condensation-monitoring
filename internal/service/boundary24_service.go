package service

import "context"
type BoundaryService24 struct { store *BoundaryStore24 }
func NewService(store *BoundaryStore24) *BoundaryService24 { return &BoundaryService24{store: store} }
func (s *BoundaryService24) Process(ctx context.Context, key string) error { return s.store.Run(ctx, key) }
