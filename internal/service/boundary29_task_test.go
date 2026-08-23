package service
import ("context"; "errors"; "testing")
func TestCandidate29(t *testing.T) { s := NewStore(); svc := NewService(s); if err := svc.Process(context.Background(), "sample-29"); err != nil { t.Fatal(err) }; err := svc.Process(context.Background(), "sample-29"); if !errors.Is(err, ErrBoundaryConflict29) { t.Fatalf("conflict identity lost: %v", err) } }
