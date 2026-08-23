package service
import ("context"; "errors"; "testing")
func TestCandidate07(t *testing.T) { s := NewStore(); svc := NewService(s); ctx, cancel := context.WithCancel(context.Background()); cancel(); if err := svc.Process(ctx, "sample-7"); !errors.Is(err, context.Canceled) { t.Fatalf("cancelled sample was accepted: %v", err) }; if _, ok := s.Get("sample-7"); ok { t.Fatalf("cancelled sample persisted") } }
