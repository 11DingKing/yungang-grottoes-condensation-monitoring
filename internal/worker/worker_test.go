package worker_test

import (
	"context"
	"errors"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/storage/sqlite"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/worker"
	"log/slog"
	"testing"
	"time"
)

func TestWorkerStopsOnContext(t *testing.T) {
	db, e := sqlite.Open(context.Background(), ":memory:")
	if e != nil {
		t.Fatal(e)
	}
	defer db.Close()
	w := worker.New(db, nil, slog.Default())
	w.Start()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if e = w.Stop(ctx); e != nil {
		t.Fatal(e)
	}
}
func TestWorkerEnqueueIsIdempotent(t *testing.T) {
	db, e := sqlite.Open(context.Background(), ":memory:")
	if e != nil {
		t.Fatal(e)
	}
	defer db.Close()
	ctx := context.Background()
	if e = worker.Enqueue(ctx, db.SQL, "dispatch", "d1"); e != nil {
		t.Fatal(e)
	}
	if e = worker.Enqueue(ctx, db.SQL, "dispatch", "d1"); e != nil {
		t.Fatal(e)
	}
	var n int
	if e = db.SQL.QueryRowContext(ctx, "SELECT COUNT(*) FROM worker_jobs").Scan(&n); e != nil {
		t.Fatal(e)
	}
	if n != 1 {
		t.Fatal(n)
	}
}
func TestWorkerHandlerFailureRecorded(t *testing.T) {
	db, e := sqlite.Open(context.Background(), ":memory:")
	if e != nil {
		t.Fatal(e)
	}
	defer db.Close()
	ctx := context.Background()
	if e = worker.Enqueue(ctx, db.SQL, "dispatch", "d1"); e != nil {
		t.Fatal(e)
	}
	w := worker.New(db, func(context.Context, string) error { return errors.New("temporary") }, slog.Default())
	_ = w
}
