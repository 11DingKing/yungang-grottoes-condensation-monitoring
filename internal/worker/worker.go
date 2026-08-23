package worker

import (
	"context"
	"database/sql"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/common"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/storage/sqlite"
	"log/slog"
	"sync"
	"time"
)

type Handler func(context.Context, string) error
type Worker struct {
	DB       *sqlite.DB
	Logger   *slog.Logger
	Handler  Handler
	interval time.Duration
	stop     chan struct{}
	done     chan struct{}
	once     sync.Once
}

func New(db *sqlite.DB, h Handler, l *slog.Logger) *Worker {
	return &Worker{DB: db, Handler: h, Logger: l, interval: 250 * time.Millisecond, stop: make(chan struct{}), done: make(chan struct{})}
}
func (w *Worker) Start() { go w.loop() }
func (w *Worker) Stop(ctx context.Context) error {
	w.once.Do(func() { close(w.stop) })
	select {
	case <-w.done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
func (w *Worker) loop() {
	defer close(w.done)
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()
	for {
		select {
		case <-w.stop:
			return
		case <-ticker.C:
			w.tick(context.Background())
		}
	}
}
func (w *Worker) tick(ctx context.Context) {
	rows, e := w.DB.SQL.QueryContext(ctx, "SELECT id,entity_id FROM worker_jobs WHERE state='queued' AND run_after<=datetime('now') ORDER BY run_after LIMIT 16")
	if e != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var id, entity string
		if rows.Scan(&id, &entity) != nil {
			continue
		}
		if e = w.run(ctx, id, entity); e != nil && w.Logger != nil {
			w.Logger.Error("worker job failed", "job", id, "error", e)
		}
	}
}
func (w *Worker) run(ctx context.Context, id, entity string) error {
	tx, e := w.DB.SQL.BeginTx(ctx, nil)
	if e != nil {
		return e
	}
	var attempts int
	if e = tx.QueryRowContext(ctx, "SELECT attempts FROM worker_jobs WHERE id=?", id).Scan(&attempts); e != nil {
		tx.Rollback()
		return e
	}
	if _, e = tx.ExecContext(ctx, "UPDATE worker_jobs SET state='running',attempts=attempts+1,updated_at=datetime('now') WHERE id=? AND state='queued'", id); e != nil {
		tx.Rollback()
		return e
	}
	if e = tx.Commit(); e != nil {
		return e
	}
	if w.Handler == nil {
		e = common.ErrUnavailable
	} else {
		e = w.Handler(ctx, entity)
	}
	state := "done"
	last := ""
	if e != nil {
		state = "queued"
		last = e.Error()
		if attempts+1 >= 5 {
			state = "dead"
		}
	}
	_, u := w.DB.SQL.ExecContext(ctx, "UPDATE worker_jobs SET state=?,last_error=?,run_after=datetime('now', ?),updated_at=datetime('now') WHERE id=?", state, last, "+1 minute", id)
	return u
}
func Enqueue(ctx context.Context, q *sql.DB, kind, entity string) error {
	now := common.Now().UTC()
	_, e := q.ExecContext(ctx, "INSERT OR IGNORE INTO worker_jobs(id,kind,entity_id,state,run_after,created_at,updated_at) VALUES(?,?,?,?,?,?,?)", common.NewID("job"), kind, entity, "queued", now, now, now)
	return e
}
