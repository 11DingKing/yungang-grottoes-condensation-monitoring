package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/common"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/request"
	"time"
)

func CreateRequest(ctx context.Context, q Tx, r request.Request) error {
	_, e := q.ExecContext(ctx, "INSERT INTO relief_requests(id,zone_id,requester_id,idempotency_key,status,priority,destination,version,submitted_at,updated_at) VALUES(?,?,?,?,?,?,?,?,?,?)", r.ID, r.ZoneID, r.RequesterID, r.IdempotencyKey, r.Status, r.Priority, r.Destination, r.Version, common.Now().UTC(), common.Now().UTC())
	if e != nil {
		return e
	}
	for _, i := range r.Items {
		if _, e = q.ExecContext(ctx, "INSERT INTO request_items(id,request_id,supply_id,quantity,allocated) VALUES(?,?,?,?,?)", i.ID, r.ID, i.SupplyID, i.Quantity, 0); e != nil {
			return e
		}
	}
	return nil
}
func GetRequest(ctx context.Context, q Tx, id string) (request.Request, error) {
	var r request.Request
	var st string
	var submitted, updated time.Time
	e := q.QueryRowContext(ctx, "SELECT id,zone_id,requester_id,idempotency_key,status,priority,destination,version,submitted_at,updated_at FROM relief_requests WHERE id=?", id).Scan(&r.ID, &r.ZoneID, &r.RequesterID, &r.IdempotencyKey, &st, &r.Priority, &r.Destination, &r.Version, &submitted, &updated)
	if errors.Is(e, sql.ErrNoRows) {
		return r, common.ErrNotFound
	}
	r.Status = request.Status(st)
	rows, e := q.QueryContext(ctx, "SELECT id,request_id,supply_id,quantity,allocated FROM request_items WHERE request_id=? ORDER BY id", id)
	if e != nil {
		return r, e
	}
	defer rows.Close()
	for rows.Next() {
		var i request.Item
		if e = rows.Scan(&i.ID, &i.RequestID, &i.SupplyID, &i.Quantity, &i.Allocated); e != nil {
			return r, e
		}
		r.Items = append(r.Items, i)
	}
	return r, rows.Err()
}
func UpdateRequest(ctx context.Context, q Tx, r request.Request, expected int) error {
	res, e := q.ExecContext(ctx, "UPDATE relief_requests SET status=?,version=?,updated_at=? WHERE id=? AND version=?", r.Status, r.Version, common.Now().UTC(), r.ID, expected)
	if e != nil {
		return e
	}
	n, _ := res.RowsAffected()
	if n != 1 {
		return common.ErrConflict
	}
	return nil
}
func SetAllocated(ctx context.Context, q Tx, requestID, supplyID string, n int) error {
	_, e := q.ExecContext(ctx, "UPDATE request_items SET allocated=allocated+? WHERE request_id=? AND supply_id=?", n, requestID, supplyID)
	return e
}
