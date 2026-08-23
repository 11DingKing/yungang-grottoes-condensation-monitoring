package service

import "context"
type BoundaryService19 struct { store *BoundaryStore19 }
func NewService(store *BoundaryStore19) *BoundaryService19 { return &BoundaryService19{store: store} }
func (s *BoundaryService19) Process(ctx context.Context, key string) error { return s.store.Put(ctx, key, "active") }
