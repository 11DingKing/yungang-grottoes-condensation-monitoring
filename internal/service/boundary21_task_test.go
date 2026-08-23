package service
import ("context"; "errors"; "testing")
func TestCandidate21(t *testing.T) { s := NewStore(); svc := NewService(s); ctx, cancel := context.WithCancel(context.Background()); cancel(); if err := svc.Process(ctx, "sample-21"); !errors.Is(err, context.Canceled) { t.Fatalf("cancelled sample was accepted: %v", err) }; if _, ok := s.Get("sample-21"); ok { t.Fatalf("cancelled sample persisted") } }
