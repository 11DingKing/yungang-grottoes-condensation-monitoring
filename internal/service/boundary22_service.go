package service

import "context"
type BoundaryService22 struct { store *BoundaryStore22 }
func NewService(store *BoundaryStore22) *BoundaryService22 { return &BoundaryService22{store: store} }
func (s *BoundaryService22) Process(ctx context.Context, key string) error { return s.store.Run(ctx, key) }
