package policy

import (
	"fmt"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/common"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/identity"
)

type Rule struct {
	Action             string
	Roles              map[identity.Role]bool
	RequiresActivePlan bool
}
type Engine struct{ rules map[string]Rule }

func New() *Engine {
	return &Engine{rules: map[string]Rule{
		"plan:create":     {Action: "plan:create", Roles: map[identity.Role]bool{identity.RoleCoordinator: true}},
		"request:approve": {Action: "request:approve", Roles: map[identity.Role]bool{identity.RoleCoordinator: true}},
		"dispatch:update": {Action: "dispatch:update", Roles: map[identity.Role]bool{identity.RoleCoordinator: true, identity.RoleDispatcher: true}},
		"audit:list":      {Action: "audit:list", Roles: map[identity.Role]bool{identity.RoleCoordinator: true, identity.RoleAuditor: true}},
	}}
}
func (e *Engine) Check(u identity.User, action string) error {
	r, ok := e.rules[action]
	if !ok {
		return fmt.Errorf("%w: unknown action", common.ErrForbidden)
	}
	if !u.Active || !r.Roles[u.Role] {
		return fmt.Errorf("%w: role %s cannot %s", common.ErrForbidden, u.Role, action)
	}
	return nil
}
func (e *Engine) Add(r Rule) {
	if e.rules == nil {
		e.rules = map[string]Rule{}
	}
	e.rules[r.Action] = r
}
