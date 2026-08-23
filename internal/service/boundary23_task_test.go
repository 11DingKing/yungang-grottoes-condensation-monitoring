package service
import ("context"; "errors"; "testing")
func TestCandidate23(t *testing.T) { s := NewStore(); svc := NewService(s); err := svc.RunOnce(context.Background(), func(context.Context) error { return errors.New("temporary device failure") }); if err == nil { t.Fatalf("temporary failure was hidden") }; job := s.BoundaryJob23(); if job.State != "pending" || job.Attempts != 1 { t.Fatalf("job cannot be retried: %#v", job) } }
