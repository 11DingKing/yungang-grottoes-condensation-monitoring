package service

import ("context"; "errors")
type BoundaryService06 struct { store *BoundaryStore06 }
func NewService(store *BoundaryStore06) *BoundaryService06 { return &BoundaryService06{store: store} }
func (s *BoundaryService06) Process(ctx context.Context, key string) error { if err := s.store.Save(ctx, key); err != nil { return errors.New(err.Error()) }; return nil }
