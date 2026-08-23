package service

import "context"
type BoundaryService09 struct { store *BoundaryStore09 }
func NewService(store *BoundaryStore09) *BoundaryService09 { return &BoundaryService09{store: store} }
func (s *BoundaryService09) Process(ctx context.Context, key string) error { return s.store.Save(ctx, key) }
