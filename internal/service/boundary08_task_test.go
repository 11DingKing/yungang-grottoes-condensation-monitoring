package service
import ("context"; "testing")
func TestCandidate08(t *testing.T) { s := NewStore(); svc := NewService(s); if err := svc.Submit(context.Background(), "report-8"); err != nil { t.Fatal(err) }; if err := svc.Submit(context.Background(), "report-8"); err != nil { t.Fatal(err) }; if got := len(s.Events()); got != 1 { t.Fatalf("duplicate idempotency event count: %d", got) } }
