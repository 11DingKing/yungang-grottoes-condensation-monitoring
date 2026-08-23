package service
import ("context"; "testing")
func TestCandidate24(t *testing.T) { s := NewStore(); svc := NewService(s); s.FailAudit(true); if err := svc.Process(context.Background(), "sample-24"); err == nil { t.Fatalf("audit failure was hidden") }; if _, ok := s.Get("sample-24"); ok { t.Fatalf("failed workflow left active state") } }
