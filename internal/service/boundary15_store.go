package service

import ("context"; "errors"; "sync")

var ErrBoundaryProbe15 = errors.New("probe failed")
type BoundaryResource15 struct { mu sync.Mutex; closed bool; fail bool }
func NewResource() *BoundaryResource15 { return &BoundaryResource15{} }
func (r *BoundaryResource15) Close() { r.mu.Lock(); defer r.mu.Unlock(); r.closed = true }
func (r *BoundaryResource15) Closed() bool { r.mu.Lock(); defer r.mu.Unlock(); return r.closed }
func (r *BoundaryResource15) Check(ctx context.Context) error { if err := ctx.Err(); err != nil { return err }; r.mu.Lock(); defer r.mu.Unlock(); if r.fail { return ErrBoundaryProbe15 }; return nil }
func (r *BoundaryResource15) Fail(v bool) { r.mu.Lock(); defer r.mu.Unlock(); r.fail = v }
