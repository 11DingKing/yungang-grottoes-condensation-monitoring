package service

import ("context"; "errors")
type BoundaryService29 struct { store *BoundaryStore29 }
func NewService(store *BoundaryStore29) *BoundaryService29 { return &BoundaryService29{store: store} }
func (s *BoundaryService29) Process(ctx context.Context, key string) error { if err := s.store.Save(ctx, key); err != nil { return errors.New(err.Error()) }; return nil }
