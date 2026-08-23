package service

import ("context"; "errors")
type BoundaryService01 struct { store *BoundaryStore01 }
func NewService(store *BoundaryStore01) *BoundaryService01 { return &BoundaryService01{store: store} }
func (s *BoundaryService01) Process(ctx context.Context, key string) error { return s.store.Save(context.Background(), key) }
var _ = errors.Is
