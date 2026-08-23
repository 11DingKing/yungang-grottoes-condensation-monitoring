package service

import "context"
type BoundaryService26 struct { store *BoundaryStore26 }
func NewService(store *BoundaryStore26) *BoundaryService26 { return &BoundaryService26{store: store} }
func (s *BoundaryService26) Process(ctx context.Context, key string) error { return s.store.Put(ctx, key, "active") }
