package service

import ("context"; "errors")
type BoundaryService07 struct { store *BoundaryStore07 }
func NewService(store *BoundaryStore07) *BoundaryService07 { return &BoundaryService07{store: store} }
func (s *BoundaryService07) Process(ctx context.Context, key string) error { return s.store.Save(context.Background(), key) }
var _ = errors.Is
