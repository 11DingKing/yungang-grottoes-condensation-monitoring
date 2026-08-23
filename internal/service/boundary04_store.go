package service

import ("context"; "errors"; "sync")

var ErrBoundaryTransition04 = errors.New("invalid transition")
type BoundaryRecord04 struct { Status string }
type BoundaryStore04 struct { mu sync.Mutex; records map[string]BoundaryRecord04 }
func NewStore() *BoundaryStore04 { return &BoundaryStore04{records: map[string]BoundaryRecord04{}} }
func (s *BoundaryStore04) Save(ctx context.Context, key string, r BoundaryRecord04) error { if err := ctx.Err(); err != nil { return err }; s.mu.Lock(); defer s.mu.Unlock(); s.records[key] = r; return nil }
func (s *BoundaryStore04) Get(key string) (BoundaryRecord04, bool) { s.mu.Lock(); defer s.mu.Unlock(); v, ok := s.records[key]; return v, ok }
