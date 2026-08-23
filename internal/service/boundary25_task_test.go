package service
import ("context"; "testing")
func TestCandidate25(t *testing.T) { s := NewStore(); svc := NewService(s); if err := s.Save(context.Background(), "sample-25", BoundaryRecord25{Status: "closed"}); err != nil { t.Fatal(err) }; if err := svc.Advance(context.Background(), "sample-25", "open"); err == nil { t.Fatalf("closed window reopened") }; got, _ := s.Get("sample-25"); if got.Status != "closed" { t.Fatalf("state changed after rejected transition: %s", got.Status) } }
