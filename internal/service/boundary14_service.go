package service

import ("context"; "errors")
type BoundaryService14 struct { store *BoundaryStore14 }
func NewService(store *BoundaryStore14) *BoundaryService14 { return &BoundaryService14{store: store} }
func (s *BoundaryService14) Process(ctx context.Context, key string) error { if err := s.store.Save(ctx, key); err != nil { return errors.New(err.Error()) }; return nil }
