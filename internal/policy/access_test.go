package policy_test

import (
	"errors"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/common"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/identity"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/policy"
	"testing"
)

func TestRoleMatrix(t *testing.T) {
	e := policy.New()
	cases := []struct {
		role   identity.Role
		action string
		ok     bool
	}{{identity.RoleCoordinator, "plan:create", true}, {identity.RoleCoordinator, "request:approve", true}, {identity.RoleDispatcher, "dispatch:update", true}, {identity.RoleAuditor, "audit:list", true}, {identity.RoleAuditor, "dispatch:update", false}, {identity.RoleDispatcher, "audit:list", false}}
	for _, c := range cases {
		u := identity.User{Role: c.role, Active: true}
		err := e.Check(u, c.action)
		if (c.ok && err != nil) || (!c.ok && !errors.Is(err, common.ErrForbidden)) {
			t.Fatalf("%+v %v", c, err)
		}
	}
}
func TestInactiveRoleDenied(t *testing.T) {
	e := policy.New()
	if err := e.Check(identity.User{Role: identity.RoleCoordinator}, "plan:create"); !errors.Is(err, common.ErrForbidden) {
		t.Fatal(err)
	}
}
