package service
import ("context"; "errors"; "testing")
func TestCandidate06(t *testing.T) { s := NewStore(); svc := NewService(s); if err := svc.Process(context.Background(), "sample-6"); err != nil { t.Fatal(err) }; err := svc.Process(context.Background(), "sample-6"); if !errors.Is(err, ErrBoundaryConflict06) { t.Fatalf("conflict identity lost: %v", err) } }
