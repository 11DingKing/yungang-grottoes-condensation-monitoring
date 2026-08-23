package service
import ("context"; "testing")
func TestCandidate09(t *testing.T) { s := NewStore(); svc := NewService(s); if err := svc.Process(context.Background(), "sample-9"); err != nil { t.Fatal(err) }; got := s.Events(); got[0] = "rewritten-9"; again := s.Events(); if again[0] != "sample-9" { t.Fatalf("event history was mutated through snapshot: %v", again) } }
