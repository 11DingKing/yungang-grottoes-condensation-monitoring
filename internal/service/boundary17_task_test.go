package service
import ("context"; "testing")
func TestCandidate17(t *testing.T) { s := NewStore(); svc := NewService(s); if err := svc.Process(context.Background(), "sample-17"); err != nil { t.Fatal(err) }; got := s.Events(); got[0] = "rewritten-17"; again := s.Events(); if again[0] != "sample-17" { t.Fatalf("event history was mutated through snapshot: %v", again) } }
