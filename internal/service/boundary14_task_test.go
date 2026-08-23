package service
import ("context"; "errors"; "testing")
func TestCandidate14(t *testing.T) { s := NewStore(); svc := NewService(s); if err := svc.Process(context.Background(), "sample-14"); err != nil { t.Fatal(err) }; err := svc.Process(context.Background(), "sample-14"); if !errors.Is(err, ErrBoundaryConflict14) { t.Fatalf("conflict identity lost: %v", err) } }
