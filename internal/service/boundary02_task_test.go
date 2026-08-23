package service
import ("context"; "sync"; "testing")
func TestCandidate02(t *testing.T) { s := NewStore(1); svc := NewService(s); start := make(chan struct{}); var wg sync.WaitGroup; successes := 0; var mu sync.Mutex; for i := 0; i < 2; i++ { wg.Add(1); go func() { defer wg.Done(); <-start; if err := svc.Reserve(context.Background(), 1); err == nil { mu.Lock(); successes++; mu.Unlock() } }() }; close(start); wg.Wait(); if successes != 1 || s.Used() > 1 { t.Fatalf("capacity oversold: successes=%d used=%d", successes, s.Used()) } }
