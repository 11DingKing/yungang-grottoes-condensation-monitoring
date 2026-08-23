package service
import ("context"; "testing")
func TestCandidate28(t *testing.T) { r := NewResource(); r.Fail(true); svc := NewService(r); if err := svc.Use(context.Background()); err == nil { t.Fatalf("failed probe was accepted") }; if !r.Closed() { t.Fatalf("probe resource leaked after failure") } }
