package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/common"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/dispatch"
)

func CreateDispatch(ctx context.Context, q Tx, d dispatch.Dispatch) error {
	_, e := q.ExecContext(ctx, "INSERT INTO dispatches(id,request_id,driver,vehicle,status,attempt,next_attempt_at,version,created_at) VALUES(?,?,?,?,?,?,?,?,?)", d.ID, d.RequestID, d.Driver, d.Vehicle, d.Status, d.Attempt, d.NextAttemptAt.UTC(), d.Version, common.Now().UTC())
	return e
}
func GetDispatch(ctx context.Context, q Tx, id string) (dispatch.Dispatch, error) {
	var d dispatch.Dispatch
	var st string
	var next, created string
	var delivered sql.NullString
	e := q.QueryRowContext(ctx, "SELECT id,request_id,driver,vehicle,status,attempt,next_attempt_at,delivered_at,version,created_at FROM dispatches WHERE id=?", id).Scan(&d.ID, &d.RequestID, &d.Driver, &d.Vehicle, &st, &d.Attempt, &next, &delivered, &d.Version, &created)
	if errors.Is(e, sql.ErrNoRows) {
		return d, common.ErrNotFound
	}
	d.Status = dispatch.Status(st)
	d.NextAttemptAt = parseTime(next)
	if delivered.Valid {
		t := parseTime(delivered.String)
		d.DeliveredAt = &t
	}
	return d, e
}
func UpdateDispatch(ctx context.Context, q Tx, d dispatch.Dispatch, expected int) error {
	res, e := q.ExecContext(ctx, "UPDATE dispatches SET status=?,attempt=?,next_attempt_at=?,delivered_at=?,version=? WHERE id=? AND version=?", d.Status, d.Attempt, d.NextAttemptAt.UTC(), d.DeliveredAt, d.Version, d.ID, expected)
	if e != nil {
		return e
	}
	n, _ := res.RowsAffected()
	if n != 1 {
		return common.ErrConflict
	}
	return nil
}
