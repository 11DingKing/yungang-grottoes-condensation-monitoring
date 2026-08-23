package service

import ("context"; "errors")
type BoundaryService18 struct { store *BoundaryStore18 }
func NewService(store *BoundaryStore18) *BoundaryService18 { return &BoundaryService18{store: store} }
func (s *BoundaryService18) Process(ctx context.Context, key string) error { if err := s.store.Save(ctx, key); err != nil { return errors.New(err.Error()) }; return nil }
