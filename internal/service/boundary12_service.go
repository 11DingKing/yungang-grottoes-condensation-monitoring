package service

import ("context"; "errors")
type BoundaryService12 struct { store *BoundaryStore12 }
func NewService(store *BoundaryStore12) *BoundaryService12 { return &BoundaryService12{store: store} }
func (s *BoundaryService12) Process(ctx context.Context, key string) error { return s.store.Save(context.Background(), key) }
var _ = errors.Is
