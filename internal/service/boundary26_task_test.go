package service
import ("context"; "testing")
func TestCandidate26(t *testing.T) { s := NewStore(); s.Restore(); svc := NewService(s); defer func() { if recovered := recover(); recovered != nil { t.Fatalf("restored store panicked on first write: %v", recovered) } }(); if err := svc.Process(context.Background(), "sample-26"); err != nil { t.Fatal(err) }; if _, ok := s.Get("sample-26"); !ok { t.Fatalf("restored sample was not stored") } }
