package service
import ("context"; "errors"; "testing")
func TestCandidate18(t *testing.T) { s := NewStore(); svc := NewService(s); if err := svc.Process(context.Background(), "sample-18"); err != nil { t.Fatal(err) }; err := svc.Process(context.Background(), "sample-18"); if !errors.Is(err, ErrBoundaryConflict18) { t.Fatalf("conflict identity lost: %v", err) } }
