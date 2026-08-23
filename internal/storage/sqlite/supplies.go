package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/common"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/supply"
	"time"
)

func CreateSupply(ctx context.Context, q Tx, id, zone, sku, name, unit string, reorder int) error {
	_, e := q.ExecContext(ctx, "INSERT INTO supplies(id,zone_id,sku,name,unit,reorder_point,created_at) VALUES(?,?,?,?,?,?,?)", id, zone, sku, name, unit, reorder, common.Now().UTC())
	return e
}
func CreateLot(ctx context.Context, q Tx, l supply.Lot) error {
	_, e := q.ExecContext(ctx, "INSERT INTO supply_lots(id,supply_id,lot_code,quantity,reserved,status,version,created_at) VALUES(?,?,?,?,?,?,?,?)", l.ID, l.SupplyID, l.LotCode, l.Quantity, l.Reserved, l.Status, l.Version, common.Now().UTC())
	return e
}
func GetLot(ctx context.Context, q Tx, id string) (supply.Lot, error) {
	var l supply.Lot
	var st string
	var exp sql.NullTime
	e := q.QueryRowContext(ctx, "SELECT id,supply_id,lot_code,quantity,reserved,expires_at,status,version FROM supply_lots WHERE id=?", id).Scan(&l.ID, &l.SupplyID, &l.LotCode, &l.Quantity, &l.Reserved, &exp, &st, &l.Version)
	if errors.Is(e, sql.ErrNoRows) {
		return l, common.ErrNotFound
	}
	l.Status = supply.LotStatus(st)
	if exp.Valid {
		l.ExpiresAt = &exp.Time
	}
	return l, e
}
func UpdateLot(ctx context.Context, q Tx, l supply.Lot, expected int) error {
	r, e := q.ExecContext(ctx, "UPDATE supply_lots SET quantity=?,reserved=?,status=?,version=? WHERE id=? AND version=?", l.Quantity, l.Reserved, l.Status, l.Version, l.ID, expected)
	if e != nil {
		return e
	}
	n, _ := r.RowsAffected()
	if n != 1 {
		return common.ErrConflict
	}
	return nil
}
func ListLots(ctx context.Context, q Tx, supplyID string) ([]supply.Lot, error) {
	rows, e := q.QueryContext(ctx, "SELECT id,supply_id,lot_code,quantity,reserved,expires_at,status,version FROM supply_lots WHERE supply_id=? ORDER BY expires_at,id", supplyID)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := make([]supply.Lot, 0)
	for rows.Next() {
		var l supply.Lot
		var st string
		var exp sql.NullTime
		if e = rows.Scan(&l.ID, &l.SupplyID, &l.LotCode, &l.Quantity, &l.Reserved, &exp, &st, &l.Version); e != nil {
			return nil, e
		}
		l.Status = supply.LotStatus(st)
		if exp.Valid {
			l.ExpiresAt = &exp.Time
		}
		out = append(out, l)
	}
	return out, rows.Err()
}

var _ = time.Time{}
