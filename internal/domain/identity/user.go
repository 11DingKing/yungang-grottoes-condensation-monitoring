package identity

import (
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/common"
	"strings"
	"time"
)

type Role string

const (
	RoleCoordinator Role = "coordinator"
	RoleDispatcher  Role = "dispatcher"
	RoleAuditor     Role = "auditor"
)

type User struct {
	ID, Username, PasswordHash string
	Role                       Role
	Active                     bool
	CreatedAt                  time.Time
}

func (u User) Can(action string) bool {
	if !u.Active {
		return false
	}
	switch strings.ToLower(action) {
	case "plan:write", "request:approve", "incident:write":
		return u.Role == RoleCoordinator
	case "dispatch:write":
		return u.Role == RoleDispatcher || u.Role == RoleCoordinator
	case "audit:read":
		return u.Role == RoleAuditor || u.Role == RoleCoordinator
	default:
		return false
	}
}
func NewUser(username string, role Role) User {
	return User{ID: common.NewID("usr"), Username: username, Role: role, Active: true, CreatedAt: common.Now()}
}
