package service

import ("context"; "errors")
type BoundaryService25 struct { store *BoundaryStore25 }
func NewService(store *BoundaryStore25) *BoundaryService25 { return &BoundaryService25{store: store} }
func (s *BoundaryService25) Advance(ctx context.Context, key, next string) error { record, _ := s.store.Get(key); record.Status = next; if err := s.store.Save(ctx, key, record); err != nil { return err }; return nil }
var _ = errors.Is
