package service

import "context"
type BoundaryService15 struct { resource *BoundaryResource15 }
func NewService(resource *BoundaryResource15) *BoundaryService15 { return &BoundaryService15{resource: resource} }
func (s *BoundaryService15) Use(ctx context.Context) error { if err := s.resource.Check(ctx); err != nil { return err }; s.resource.Close(); return nil }
