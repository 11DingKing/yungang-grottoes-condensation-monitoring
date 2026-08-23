package service_test

import (
	"context"
	"errors"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/common"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/dispatch"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/identity"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/plan"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/request"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/supply"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/service"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/storage/sqlite"
	"testing"
	"time"
)

func setup(t *testing.T) (context.Context, *service.Service, identity.User) {
	ctx := context.Background()
	db, e := sqlite.Open(ctx, ":memory:")
	if e != nil {
		t.Fatal(e)
	}
	t.Cleanup(func() { db.Close() })
	s := service.New(db)
	u := identity.NewUser("coordinator", identity.RoleCoordinator)
	u.PasswordHash = "pw"
	if e = s.CreateUser(ctx, u); e != nil {
		t.Fatal(e)
	}
	if e = s.CreateZone(ctx, "z1", "FJ-01", "coastal", "east"); e != nil {
		t.Fatal(e)
	}
	return ctx, s, u
}
func TestLoginAndLogout(t *testing.T) {
	ctx, s, _ := setup(t)
	ses, e := s.Login(ctx, "coordinator", "pw")
	if e != nil {
		t.Fatal(e)
	}
	if _, e = s.Authenticate(ctx, ses.ID); e != nil {
		t.Fatal(e)
	}
	if e = s.Logout(ctx, ses.ID); e != nil {
		t.Fatal(e)
	}
	if _, e = s.Authenticate(ctx, ses.ID); !errors.Is(e, common.ErrForbidden) {
		t.Fatal(e)
	}
}
func TestLoginRejectsWrongPassword(t *testing.T) {
	ctx, s, _ := setup(t)
	if _, e := s.Login(ctx, "coordinator", "bad"); !errors.Is(e, common.ErrForbidden) {
		t.Fatal(e)
	}
}
func TestCreateAndActivatePlan(t *testing.T) {
	ctx, s, u := setup(t)
	p, e := plan.New("z1", "level four", 4)
	if e != nil {
		t.Fatal(e)
	}
	if e = s.CreatePlan(ctx, u, p); e != nil {
		t.Fatal(e)
	}
	if e = s.TransitionPlan(ctx, u, p.ID, plan.Active); e != nil {
		t.Fatal(e)
	}
	if e = s.TransitionPlan(ctx, u, p.ID, plan.Closed); e != nil {
		t.Fatal(e)
	}
}
func TestUnauthorizedPlan(t *testing.T) {
	ctx, s, _ := setup(t)
	u := identity.NewUser("auditor", identity.RoleAuditor)
	p, _ := plan.New("z1", "plan", 1)
	if e := s.CreatePlan(ctx, u, p); !errors.Is(e, common.ErrForbidden) {
		t.Fatal(e)
	}
}
func TestApprovalReservesStock(t *testing.T) {
	ctx, s, u := setup(t)
	if e := s.CreateSupply(ctx, "sp1", "z1", "WATER", "Water", "box", 2); e != nil {
		t.Fatal(e)
	}
	if e := s.CreateLot(ctx, supply.Lot{ID: "lot1", SupplyID: "sp1", LotCode: "A", Quantity: 10, Status: supply.Available, Version: 1}); e != nil {
		t.Fatal(e)
	}
	r := request.Request{ID: "req1", ZoneID: "z1", RequesterID: u.ID, IdempotencyKey: "idem1", Status: request.Submitted, Priority: 1, Destination: "town", Version: 1, Items: []request.Item{{ID: "item1", SupplyID: "sp1", Quantity: 6}}}
	if e := s.CreateRequest(ctx, r); e != nil {
		t.Fatal(e)
	}
	if e := s.ApproveRequest(ctx, u, r.ID, "http-1"); e != nil {
		t.Fatal(e)
	}
	if e := s.ApproveRequest(ctx, u, r.ID, "http-2"); !errors.Is(e, common.ErrConflict) {
		t.Fatal(e)
	}
}
func TestApprovalRollsBackWhenShort(t *testing.T) {
	ctx, s, u := setup(t)
	if e := s.CreateSupply(ctx, "sp1", "z1", "WATER", "Water", "box", 2); e != nil {
		t.Fatal(e)
	}
	if e := s.CreateLot(ctx, supply.Lot{ID: "lot1", SupplyID: "sp1", LotCode: "A", Quantity: 2, Status: supply.Available, Version: 1}); e != nil {
		t.Fatal(e)
	}
	r := request.Request{ID: "req1", ZoneID: "z1", RequesterID: u.ID, IdempotencyKey: "idem1", Status: request.Submitted, Priority: 1, Destination: "town", Version: 1, Items: []request.Item{{ID: "item1", SupplyID: "sp1", Quantity: 5}}}
	if e := s.CreateRequest(ctx, r); e != nil {
		t.Fatal(e)
	}
	if e := s.ApproveRequest(ctx, u, r.ID, "http-1"); !errors.Is(e, common.ErrConflict) {
		t.Fatal(e)
	}
	var reserved int
	if e := s.DB.SQL.QueryRowContext(ctx, "SELECT reserved FROM supply_lots WHERE id='lot1'").Scan(&reserved); e != nil {
		t.Fatal(e)
	}
	if reserved != 0 {
		t.Fatal("rollback leaked", reserved)
	}
}
func TestDispatchOperations(t *testing.T) {
	ctx, s, u := setup(t)
	if e := s.CreateRequest(ctx, request.Request{ID: "r1", ZoneID: "z1", RequesterID: u.ID, IdempotencyKey: "dispatch-request", Status: request.Submitted, Priority: 1, Destination: "town", Version: 1}); e != nil {
		t.Fatal(e)
	}
	d := dispatch.Dispatch{ID: "d1", RequestID: "r1", Driver: "driver", Vehicle: "van", Status: dispatch.Queued, NextAttemptAt: time.Now(), Version: 1}
	if e := s.CreateDispatch(ctx, u, d); e != nil {
		t.Fatal(e)
	}
	for _, op := range []string{"start", "deliver"} {
		if e := s.UpdateDispatch(ctx, u, d.ID, op); e != nil {
			t.Fatal(e)
		}
	}
}
func TestIdempotentRequest(t *testing.T) {
	ctx, s, u := setup(t)
	r := request.Request{ID: "r1", ZoneID: "z1", RequesterID: u.ID, IdempotencyKey: "same", Status: request.Submitted, Priority: 1, Destination: "town", Version: 1}
	if e := s.CreateRequest(ctx, r); e != nil {
		t.Fatal(e)
	}
	r.ID = "r2"
	if e := s.CreateRequest(ctx, r); e == nil {
		t.Fatal("duplicate accepted")
	}
}
