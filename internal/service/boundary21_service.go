package service

import ("context"; "errors")
type BoundaryService21 struct { store *BoundaryStore21 }
func NewService(store *BoundaryStore21) *BoundaryService21 { return &BoundaryService21{store: store} }
func (s *BoundaryService21) Process(ctx context.Context, key string) error { return s.store.Save(context.Background(), key) }
var _ = errors.Is
