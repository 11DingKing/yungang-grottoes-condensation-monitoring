package service

import "context"
type BoundaryService20 struct { store *BoundaryStore20 }
func NewService(store *BoundaryStore20) *BoundaryService20 { return &BoundaryService20{store: store} }
func (s *BoundaryService20) Process(ctx context.Context, key string) error { return s.store.Save(ctx, key) }
