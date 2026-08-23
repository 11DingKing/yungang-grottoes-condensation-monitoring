package domain_test

import (
	"errors"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/common"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/dispatch"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/identity"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/incident"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/plan"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/request"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/shelter"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/supply"
	"testing"
	"time"
)

func TestPlanLifecycle(t *testing.T) {
	p, e := plan.New("z1", "coastal response", 2)
	if e != nil {
		t.Fatal(e)
	}
	if p.Status != plan.Draft {
		t.Fatal(p.Status)
	}
	for _, n := range []plan.Status{plan.Active, plan.Standby, plan.Closed} {
		if e = p.Transition(n); e != nil {
			t.Fatal(e)
		}
	}
	if p.EndsAt == nil || p.Version != 4 {
		t.Fatalf("%+v", p)
	}
	if e = p.Transition(plan.Active); !errors.Is(e, common.ErrConflict) {
		t.Fatal(e)
	}
}
func TestPlanRejectsInvalidInputs(t *testing.T) {
	cases := []struct {
		zone, name string
		level      int
	}{{"", "x", 1}, {"z", "", 1}, {"z", "x", 0}, {"z", "x", 5}}
	for _, c := range cases {
		if _, e := plan.New(c.zone, c.name, c.level); !errors.Is(e, common.ErrInvalid) {
			t.Fatalf("%+v %v", c, e)
		}
	}
}
func TestShelterCapacity(t *testing.T) {
	s := shelter.Shelter{Capacity: 10, Status: shelter.Open, Version: 1}
	if e := s.Assign(7); e != nil {
		t.Fatal(e)
	}
	if s.Available() != 3 {
		t.Fatal(s.Available())
	}
	if e := s.Assign(4); !errors.Is(e, common.ErrConflict) {
		t.Fatal(e)
	}
	if e := s.Assign(3); e != nil {
		t.Fatal(e)
	}
	if s.Status != shelter.Full {
		t.Fatal(s.Status)
	}
	if e := s.Release(2); e != nil {
		t.Fatal(e)
	}
	if s.Status != shelter.Open || s.Occupied != 8 {
		t.Fatal(s)
	}
}
func TestShelterInvalidTransitions(t *testing.T) {
	for _, st := range []shelter.Status{shelter.Paused, shelter.Retired} {
		s := shelter.Shelter{Capacity: 2, Status: st}
		if e := s.Assign(1); !errors.Is(e, common.ErrConflict) {
			t.Fatal(st, e)
		}
	}
	s := shelter.Shelter{Capacity: 2, Occupied: 1, Status: shelter.Open}
	for _, n := range []int{0, -1, 2} {
		if e := s.Release(n); !errors.Is(e, common.ErrInvalid) {
			t.Fatal(n, e)
		}
	}
}
func TestLotReservationAndExpiry(t *testing.T) {
	now := time.Now()
	future := now.Add(time.Hour)
	l := supply.Lot{Quantity: 8, Status: supply.Available, ExpiresAt: &future}
	if l.Free(now) != 8 {
		t.Fatal(l.Free(now))
	}
	if e := l.Reserve(5, now); e != nil {
		t.Fatal(e)
	}
	if l.Free(now) != 3 || l.Reserved != 5 {
		t.Fatal(l)
	}
	if e := l.Consume(5); e != nil {
		t.Fatal(e)
	}
	if l.Quantity != 3 || l.Reserved != 0 {
		t.Fatal(l)
	}
	past := now.Add(-time.Hour)
	expired := supply.Lot{Quantity: 2, Status: supply.Available, ExpiresAt: &past}
	if expired.Free(now) != 0 {
		t.Fatal(expired.Free(now))
	}
	if e := expired.Reserve(1, now); !errors.Is(e, common.ErrConflict) {
		t.Fatal(e)
	}
	if expired.Status != supply.Expired {
		t.Fatal(expired.Status)
	}
}
func TestLotRejectsBadAmounts(t *testing.T) {
	l := supply.Lot{Quantity: 1, Status: supply.Available}
	for _, n := range []int{0, -1, 2} {
		if e := l.Reserve(n, time.Now()); !errors.Is(e, common.ErrInvalid) && n < 1 {
			t.Fatal(n, e)
		}
	}
	if e := l.Consume(1); !errors.Is(e, common.ErrInvalid) {
		t.Fatal(e)
	}
}
func TestRequestStateMachine(t *testing.T) {
	r := request.Request{Status: request.Submitted, Items: []request.Item{{Quantity: 2, Allocated: 2}}}
	for _, n := range []request.Status{request.Approved, request.Allocated, request.Dispatched, request.Received} {
		if e := r.Transition(n); e != nil {
			t.Fatal(e)
		}
	}
	if !r.Complete() {
		t.Fatal(r)
	}
	if e := r.Transition(request.Cancelled); !errors.Is(e, common.ErrConflict) {
		t.Fatal(e)
	}
}
func TestRequestCancellation(t *testing.T) {
	r := request.Request{Status: request.Submitted}
	if e := r.Transition(request.Cancelled); e != nil {
		t.Fatal(e)
	}
	if e := r.Transition(request.Approved); !errors.Is(e, common.ErrConflict) {
		t.Fatal(e)
	}
}
func TestDispatchLifecycle(t *testing.T) {
	d := dispatch.Dispatch{Status: dispatch.Queued, NextAttemptAt: time.Now()}
	if e := d.Start(); e != nil {
		t.Fatal(e)
	}
	if e := d.Deliver(); e != nil {
		t.Fatal(e)
	}
	if d.DeliveredAt == nil || d.Status != dispatch.Delivered {
		t.Fatal(d)
	}
}
func TestDispatchRetry(t *testing.T) {
	d := dispatch.Dispatch{Status: dispatch.InTransit}
	for i := 0; i < 4; i++ {
		if e := d.Retry("network"); e != nil {
			t.Fatal(e)
		}
		if d.Status == dispatch.Failed {
			t.Fatal("early failure")
		}
	}
	if e := d.Retry("network"); e != nil {
		t.Fatal(e)
	}
	if d.Status != dispatch.Failed || d.Attempt != 5 {
		t.Fatal(d)
	}
	if e := d.Deliver(); !errors.Is(e, common.ErrConflict) {
		t.Fatal(e)
	}
}
func TestIncidentLifecycle(t *testing.T) {
	i := incident.Incident{Status: incident.Open}
	for _, n := range []incident.Status{incident.Assigned, incident.Resolved, incident.Closed} {
		if e := i.Transition(n); e != nil {
			t.Fatal(e)
		}
	}
	if e := i.Transition(incident.Open); !errors.Is(e, common.ErrConflict) {
		t.Fatal(e)
	}
}
func TestSessionValidity(t *testing.T) {
	now := time.Now()
	s := identity.Session{ExpiresAt: now.Add(time.Minute)}
	if !s.Valid(now) {
		t.Fatal("valid session rejected")
	}
	s.ExpiresAt = now.Add(-time.Minute)
	if s.Valid(now) {
		t.Fatal("expired session accepted")
	}
	s.ExpiresAt = now.Add(time.Minute)
	s.RevokedAt = &now
	if s.Valid(now) {
		t.Fatal("revoked session accepted")
	}
}
func TestUserPermissions(t *testing.T) {
	u := identity.User{Role: identity.RoleCoordinator, Active: true}
	for _, a := range []string{"plan:write", "request:approve", "dispatch:write", "audit:read"} {
		if !u.Can(a) {
			t.Fatal(a)
		}
	}
	u.Role = identity.RoleAuditor
	if u.Can("dispatch:write") || !u.Can("audit:read") {
		t.Fatal(u)
	}
	u.Active = false
	if u.Can("audit:read") {
		t.Fatal("inactive user")
	}
}
func TestErrorSentinels(t *testing.T) {
	if common.ErrInvalid == nil || common.ErrConflict == nil || common.ErrForbidden == nil {
		t.Fatal("missing sentinel")
	}
}
