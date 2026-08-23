package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/common"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/plan"
)

func CreateZone(ctx context.Context, q Tx, id, code, name, coast string) error {
	_, e := q.ExecContext(ctx, "INSERT INTO zones(id,code,name,coastline,created_at) VALUES(?,?,?,?,?)", id, code, name, coast, common.Now().UTC())
	return e
}
func CreatePlan(ctx context.Context, q Tx, p plan.Plan) error {
	_, e := q.ExecContext(ctx, "INSERT INTO response_plans(id,zone_id,name,level,status,version,starts_at,created_at,updated_at) VALUES(?,?,?,?,?,?,?,?,?)", p.ID, p.ZoneID, p.Name, p.Level, p.Status, p.Version, p.StartsAt.UTC(), p.StartsAt.UTC(), p.UpdatedAt.UTC())
	return e
}
func GetPlan(ctx context.Context, q Tx, id string) (plan.Plan, error) {
	var p plan.Plan
	var status string
	var start, updated string
	var end sql.NullString
	e := q.QueryRowContext(ctx, "SELECT id,zone_id,name,level,status,version,starts_at,ends_at,updated_at FROM response_plans WHERE id=?", id).Scan(&p.ID, &p.ZoneID, &p.Name, &p.Level, &status, &p.Version, &start, &end, &updated)
	if errors.Is(e, sql.ErrNoRows) {
		return p, common.ErrNotFound
	}
	p.Status = plan.Status(status)
	p.StartsAt = parseTime(start)
	p.UpdatedAt = parseTime(updated)
	if end.Valid {
		t := parseTime(end.String)
		p.EndsAt = &t
	}
	return p, e
}
func UpdatePlan(ctx context.Context, q Tx, p plan.Plan, expected int) error {
	r, e := q.ExecContext(ctx, "UPDATE response_plans SET status=?,version=?,ends_at=?,updated_at=? WHERE id=? AND version=?", p.Status, p.Version, p.EndsAt, p.UpdatedAt.UTC(), p.ID, expected)
	if e != nil {
		return e
	}
	n, _ := r.RowsAffected()
	if n != 1 {
		return common.ErrConflict
	}
	return nil
}
