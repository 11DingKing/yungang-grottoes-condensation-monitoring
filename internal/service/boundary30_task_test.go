package service
import ("context"; "testing")
func TestCandidate30(t *testing.T) { s := NewStore(); svc := NewService(s); if err := svc.Submit(context.Background(), "report-30"); err != nil { t.Fatal(err) }; if err := svc.Submit(context.Background(), "report-30"); err != nil { t.Fatal(err) }; if got := len(s.Events()); got != 1 { t.Fatalf("duplicate idempotency event count: %d", got) } }
