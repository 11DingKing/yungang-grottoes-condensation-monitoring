package service

import "context"
type BoundaryService23 struct { store *BoundaryStore23 }
func NewService(store *BoundaryStore23) *BoundaryService23 { return &BoundaryService23{store: store} }
func (s *BoundaryService23) RunOnce(ctx context.Context, handler func(context.Context) error) error { j := s.store.BoundaryJob23(); j.Attempts++; j.State = "done"; s.store.Save(j); return handler(ctx) }
