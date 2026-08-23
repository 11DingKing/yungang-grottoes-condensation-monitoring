package service

import "context"
type BoundaryService10 struct { store *BoundaryStore10 }
func NewService(store *BoundaryStore10) *BoundaryService10 { return &BoundaryService10{store: store} }
func (s *BoundaryService10) RunOnce(ctx context.Context, handler func(context.Context) error) error { j := s.store.BoundaryJob10(); j.Attempts++; j.State = "done"; s.store.Save(j); return handler(ctx) }
