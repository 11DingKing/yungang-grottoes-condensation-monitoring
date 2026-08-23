package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/common"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/identity"
)

func CreateUser(ctx context.Context, q Tx, u identity.User) error {
	_, err := q.ExecContext(ctx, "INSERT INTO users(id,username,password_hash,role,active,created_at) VALUES(?,?,?,?,?,?)", u.ID, u.Username, u.PasswordHash, u.Role, u.Active, u.CreatedAt.UTC())
	return err
}
func FindUser(ctx context.Context, q Tx, username string) (identity.User, error) {
	var u identity.User
	var role string
	var active int
	var created string
	err := q.QueryRowContext(ctx, "SELECT id,username,password_hash,role,active,created_at FROM users WHERE username=?", username).Scan(&u.ID, &u.Username, &u.PasswordHash, &role, &active, &created)
	if errors.Is(err, sql.ErrNoRows) {
		return u, common.ErrNotFound
	}
	u.Role = identity.Role(role)
	u.Active = active == 1
	u.CreatedAt = parseTime(created)
	return u, err
}
func SaveSession(ctx context.Context, q Tx, s identity.Session) error {
	_, err := q.ExecContext(ctx, "INSERT INTO sessions(id,user_id,expires_at,created_at) VALUES(?,?,?,?)", s.ID, s.UserID, s.ExpiresAt.UTC(), s.CreatedAt.UTC())
	return err
}
func RevokeSession(ctx context.Context, q Tx, id string) error {
	_, err := q.ExecContext(ctx, "UPDATE sessions SET revoked_at=? WHERE id=?", common.Now().UTC(), id)
	return err
}
func FindSession(ctx context.Context, q Tx, id string) (identity.Session, error) {
	var s identity.Session
	var exp, created string
	var revoked sql.NullString
	err := q.QueryRowContext(ctx, "SELECT id,user_id,expires_at,created_at,revoked_at FROM sessions WHERE id=?", id).Scan(&s.ID, &s.UserID, &exp, &created, &revoked)
	if errors.Is(err, sql.ErrNoRows) {
		return s, common.ErrNotFound
	}
	s.ExpiresAt = parseTime(exp)
	s.CreatedAt = parseTime(created)
	if revoked.Valid {
		t := parseTime(revoked.String)
		s.RevokedAt = &t
	}
	return s, err
}
