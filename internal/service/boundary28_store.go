package service

import ("context"; "errors"; "sync")

var ErrBoundaryProbe28 = errors.New("probe failed")
type BoundaryResource28 struct { mu sync.Mutex; closed bool; fail bool }
func NewResource() *BoundaryResource28 { return &BoundaryResource28{} }
func (r *BoundaryResource28) Close() { r.mu.Lock(); defer r.mu.Unlock(); r.closed = true }
func (r *BoundaryResource28) Closed() bool { r.mu.Lock(); defer r.mu.Unlock(); return r.closed }
func (r *BoundaryResource28) Check(ctx context.Context) error { if err := ctx.Err(); err != nil { return err }; r.mu.Lock(); defer r.mu.Unlock(); if r.fail { return ErrBoundaryProbe28 }; return nil }
func (r *BoundaryResource28) Fail(v bool) { r.mu.Lock(); defer r.mu.Unlock(); r.fail = v }
