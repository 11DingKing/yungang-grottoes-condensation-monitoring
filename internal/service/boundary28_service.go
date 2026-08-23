package service

import "context"
type BoundaryService28 struct { resource *BoundaryResource28 }
func NewService(resource *BoundaryResource28) *BoundaryService28 { return &BoundaryService28{resource: resource} }
func (s *BoundaryService28) Use(ctx context.Context) error { if err := s.resource.Check(ctx); err != nil { return err }; s.resource.Close(); return nil }
