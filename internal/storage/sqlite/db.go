package sqlite

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"time"
)

//go:embed migration.sql
var migrationFS embed.FS

//go:embed monitoring.sql
var monitoringMigrationFS embed.FS

type DB struct{ SQL *sql.DB }

func parseTime(value string) time.Time {
	for _, layout := range []string{time.RFC3339Nano, time.RFC3339, "2006-01-02 15:04:05-07:00", "2006-01-02 15:04:05"} {
		if t, e := time.Parse(layout, value); e == nil {
			return t
		}
	}
	return time.Time{}
}

func Open(ctx context.Context, path string) (*DB, error) {
	if path == "" {
		path = ":memory:"
	}
	d, err := sql.Open("sqlite3", path+"?_foreign_keys=on&_busy_timeout=5000")
	if err != nil {
		return nil, err
	}
	d.SetMaxOpenConns(1)
	if err = d.PingContext(ctx); err != nil {
		d.Close()
		return nil, err
	}
	out := &DB{SQL: d}
	if err = out.Migrate(ctx); err != nil {
		d.Close()
		return nil, err
	}
	return out, nil
}
func (d *DB) Close() error { return d.SQL.Close() }
func (d *DB) Migrate(ctx context.Context) error {
	b, err := migrationFS.ReadFile("migration.sql")
	if err != nil {
		return fmt.Errorf("read migration: %w", err)
	}
	if _, err = d.SQL.ExecContext(ctx, string(b)); err != nil {
		return err
	}
	if _, err = d.SQL.ExecContext(ctx, "INSERT OR IGNORE INTO schema_migrations(version, applied_at) VALUES(1, datetime('now'))"); err != nil {
		return err
	}
	monitoringSQL, err := monitoringMigrationFS.ReadFile("monitoring.sql")
	if err != nil {
		return fmt.Errorf("read monitoring migration: %w", err)
	}
	if _, err = d.SQL.ExecContext(ctx, string(monitoringSQL)); err != nil {
		return fmt.Errorf("apply monitoring migration: %w", err)
	}
	if _, err = d.SQL.ExecContext(ctx, "INSERT OR IGNORE INTO schema_migrations(version, applied_at) VALUES(2, datetime('now'))"); err != nil {
		return err
	}
	return seed(ctx, d.SQL)
}
func seed(ctx context.Context, q *sql.DB) error {
	if os.Getenv("SEED_DEMO") == "" {
		return nil
	}
	_, err := q.ExecContext(ctx, "INSERT OR IGNORE INTO users(id,username,password_hash,role,created_at) VALUES('usr-demo','coordinator','demo','coordinator',datetime('now'))")
	return err
}

type Tx interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
}
