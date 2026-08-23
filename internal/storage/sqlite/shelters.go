package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/common"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/shelter"
)

func CreateShelter(ctx context.Context, q Tx, s shelter.Shelter) error {
	_, e := q.ExecContext(ctx, "INSERT INTO shelters(id,zone_id,name,capacity,occupied,status,version,created_at) VALUES(?,?,?,?,?,?,?,?)", s.ID, s.ZoneID, s.Name, s.Capacity, s.Occupied, s.Status, s.Version, common.Now().UTC())
	return e
}
func GetShelter(ctx context.Context, q Tx, id string) (shelter.Shelter, error) {
	var s shelter.Shelter
	var st string
	e := q.QueryRowContext(ctx, "SELECT id,zone_id,name,capacity,occupied,status,version FROM shelters WHERE id=?", id).Scan(&s.ID, &s.ZoneID, &s.Name, &s.Capacity, &s.Occupied, &st, &s.Version)
	if errors.Is(e, sql.ErrNoRows) {
		return s, common.ErrNotFound
	}
	s.Status = shelter.Status(st)
	return s, e
}
func UpdateShelter(ctx context.Context, q Tx, s shelter.Shelter, expected int) error {
	r, e := q.ExecContext(ctx, "UPDATE shelters SET occupied=?,status=?,version=? WHERE id=? AND version=?", s.Occupied, s.Status, s.Version, s.ID, expected)
	if e != nil {
		return e
	}
	n, _ := r.RowsAffected()
	if n != 1 {
		return common.ErrConflict
	}
	return nil
}
