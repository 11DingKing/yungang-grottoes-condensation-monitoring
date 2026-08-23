package idempotency

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/common"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/storage/sqlite"
	"time"
)

type Record struct {
	Key, Hash, Body      string
	Code                 int
	CreatedAt, ExpiresAt time.Time
}

func Find(ctx context.Context, q sqlite.Tx, key string) (Record, error) {
	var r Record
	e := q.QueryRowContext(ctx, "SELECT key,request_hash,response_body,status_code,created_at,expires_at FROM idempotency_keys WHERE key=?", key).Scan(&r.Key, &r.Hash, &r.Body, &r.Code, &r.CreatedAt, &r.ExpiresAt)
	if e == sql.ErrNoRows {
		return r, common.ErrNotFound
	}
	return r, e
}
func Save(ctx context.Context, q sqlite.Tx, key, hash string, code int, v any) error {
	b, e := json.Marshal(v)
	if e != nil {
		return e
	}
	now := common.Now()
	_, e = q.ExecContext(ctx, "INSERT INTO idempotency_keys(key,request_hash,response_body,status_code,created_at,expires_at) VALUES(?,?,?,?,?,?)", key, hash, string(b), code, now.UTC(), now.Add(24*time.Hour).UTC())
	return e
}
