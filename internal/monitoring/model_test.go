package monitoring

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

func validSample() Sample {
	return Sample{ID: "s1", CaveCode: "Y9", SensorID: "m-001", WallTemperature: 18, RelativeHumidity: 100, OutsideTemperature: 24, OutsideHumidity: 60, RecordedAt: time.Date(2026, 8, 23, 8, 0, 0, 0, time.UTC)}
}

func TestEvaluateCondensationAndStableBranches(t *testing.T) {
	now := validSample().RecordedAt.Add(time.Minute)
	d, err := Evaluate(context.Background(), validSample(), now, time.Hour)
	if err != nil || d.State != StateCondensation || d.Action == "observe" || d.DewPoint <= 0 {
		t.Fatalf("unexpected risk decision: %#v %v", d, err)
	}
	s := validSample()
	s.RelativeHumidity = 55
	d, err = Evaluate(context.Background(), s, now, time.Hour)
	if err != nil || d.State != StateStable || d.Action != "observe" {
		t.Fatalf("unexpected stable decision: %#v %v", d, err)
	}
}

func TestEvaluateRejectsInvalidStaleAndCancelled(t *testing.T) {
	s := validSample()
	if _, err := Evaluate(context.Background(), s, s.RecordedAt.Add(2*time.Hour), time.Hour); !errors.Is(err, ErrStaleSample) {
		t.Fatalf("expected stale error, got %v", err)
	}
	s.RelativeHumidity = 101
	if _, err := Evaluate(context.Background(), s, s.RecordedAt, time.Hour); !errors.Is(err, ErrInvalidSample) {
		t.Fatalf("expected invalid error, got %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := Evaluate(ctx, validSample(), s.RecordedAt, time.Hour); !errors.Is(err, context.Canceled) {
		t.Fatalf("expected cancellation, got %v", err)
	}
}

func TestLedgerRejectsOutOfOrderAndReturnsIsolatedSnapshot(t *testing.T) {
	l := NewLedger()
	s := validSample()
	if _, err := l.Record(context.Background(), s, s.RecordedAt, time.Hour); err != nil {
		t.Fatal(err)
	}
	old := s
	old.RecordedAt = old.RecordedAt.Add(-time.Second)
	if _, err := l.Record(context.Background(), old, s.RecordedAt, time.Hour); err == nil {
		t.Fatal("expected out-of-order rejection")
	}
	snapshot := l.Snapshot()
	snapshot[0].SensorID = "mutated"
	if got := l.Snapshot()[0].SensorID; got != "m-001" {
		t.Fatalf("snapshot leaked internal state: %s", got)
	}
}

func TestLedgerConcurrentSensors(t *testing.T) {
	l := NewLedger()
	base := validSample()
	var wg sync.WaitGroup
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			s := base
			s.SensorID = "m-" + string(rune('a'+i))
			s.RecordedAt = base.RecordedAt.Add(time.Duration(i) * time.Second)
			if _, err := l.Record(context.Background(), s, s.RecordedAt, time.Hour); err != nil {
				t.Errorf("record %d: %v", i, err)
			}
		}(i)
	}
	wg.Wait()
	if got := len(l.Snapshot()); got != 32 {
		t.Fatalf("expected 32 sensors, got %d", got)
	}
}
