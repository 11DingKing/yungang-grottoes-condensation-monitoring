package service
import ("context"; "testing")
func TestCandidate20(t *testing.T) { s := NewStore(); svc := NewService(s); if err := svc.Process(context.Background(), "sample-20"); err != nil { t.Fatal(err) }; got := s.Events(); got[0] = "rewritten-20"; again := s.Events(); if again[0] != "sample-20" { t.Fatalf("event history was mutated through snapshot: %v", again) } }
