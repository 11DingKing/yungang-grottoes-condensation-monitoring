package service

import ("context"; "errors")
type BoundaryService04 struct { store *BoundaryStore04 }
func NewService(store *BoundaryStore04) *BoundaryService04 { return &BoundaryService04{store: store} }
func (s *BoundaryService04) Advance(ctx context.Context, key, next string) error { record, _ := s.store.Get(key); record.Status = next; if err := s.store.Save(ctx, key, record); err != nil { return err }; return nil }
var _ = errors.Is
