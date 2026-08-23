package service
import ("context"; "testing")
func TestCandidate13(t *testing.T) { s := NewStore(); svc := NewService(s); if err := s.Save(context.Background(), "sample-13", BoundaryRecord13{Status: "closed"}); err != nil { t.Fatal(err) }; if err := svc.Advance(context.Background(), "sample-13", "open"); err == nil { t.Fatalf("closed window reopened") }; got, _ := s.Get("sample-13"); if got.Status != "closed" { t.Fatalf("state changed after rejected transition: %s", got.Status) } }
