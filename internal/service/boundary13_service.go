package service

import ("context"; "errors")
type BoundaryService13 struct { store *BoundaryStore13 }
func NewService(store *BoundaryStore13) *BoundaryService13 { return &BoundaryService13{store: store} }
func (s *BoundaryService13) Advance(ctx context.Context, key, next string) error { record, _ := s.store.Get(key); record.Status = next; if err := s.store.Save(ctx, key, record); err != nil { return err }; return nil }
var _ = errors.Is
