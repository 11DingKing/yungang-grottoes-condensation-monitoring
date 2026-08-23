package service
import ("context"; "testing")
func TestCandidate05(t *testing.T) { s := NewStore(); svc := NewService(s); if err := svc.Process(context.Background(), "sample-5"); err != nil { t.Fatal(err) }; got := s.Events(); got[0] = "rewritten-5"; again := s.Events(); if again[0] != "sample-5" { t.Fatalf("event history was mutated through snapshot: %v", again) } }
