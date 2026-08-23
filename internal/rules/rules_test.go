package rules_test

import (
	"errors"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/common"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/rules"
	"testing"
	"time"
)

func TestPriorityParsing(t *testing.T) {
	for i := 1; i <= 5; i++ {
		p, e := rules.ParsePriority(i)
		if e != nil || !p.Valid() {
			t.Fatal(i, e)
		}
	}
	for _, i := range []int{0, 6, -1, 99} {
		if _, e := rules.ParsePriority(i); !errors.Is(e, common.ErrInvalid) {
			t.Fatal(i, e)
		}
	}
}
func TestDeadline(t *testing.T) {
	now := time.Now()
	d := rules.NewDeadline(rules.Critical, now)
	if d.Expired(now) || d.Remaining(now) <= 0 {
		t.Fatal(d)
	}
	if !d.Expired(d.Due) {
		t.Fatal("due not expired")
	}
}
func TestSortUrgency(t *testing.T) {
	now := time.Now()
	items := []rules.Deadline{{Priority: rules.Routine, Due: now}, {Priority: rules.Emergency, Due: now.Add(time.Hour)}, {Priority: rules.Critical, Due: now}}
	got := rules.SortByUrgency(items)
	if got[0].Priority != rules.Emergency {
		t.Fatal(got)
	}
}
func TestCapacityLifecycle(t *testing.T) {
	c, e := rules.NewCapacity(10)
	if e != nil {
		t.Fatal(e)
	}
	if e = c.Reserve(4); e != nil {
		t.Fatal(e)
	}
	if c.Available() != 6 {
		t.Fatal(c.Available())
	}
	if e = c.Reserve(7); !errors.Is(e, common.ErrConflict) {
		t.Fatal(e)
	}
	if e = c.Release(2); e != nil {
		t.Fatal(e)
	}
	if e = c.Resize(1); !errors.Is(e, common.ErrConflict) {
		t.Fatal(e)
	}
	if e = c.Resize(20); e != nil {
		t.Fatal(e)
	}
}
func TestCapacityInvalid(t *testing.T) {
	if _, e := rules.NewCapacity(0); !errors.Is(e, common.ErrInvalid) {
		t.Fatal(e)
	}
	c, _ := rules.NewCapacity(2)
	for _, n := range []int{0, -1, 3} {
		if e := c.Reserve(n); !errors.Is(e, common.ErrInvalid) && n < 1 {
			t.Fatal(n, e)
		}
	}
}
func TestGraph(t *testing.T) {
	g := rules.NewGraph()
	g.Add("draft", "active")
	g.Add("active", "closed")
	if !g.Can("draft", "active") {
		t.Fatal()
	}
	if e := g.Move("draft", "closed"); !errors.Is(e, common.ErrConflict) {
		t.Fatal(e)
	}
	if len(g.Reachable("draft")) != 3 {
		t.Fatal(g.Reachable("draft"))
	}
	if e := g.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestQuota(t *testing.T) {
	q := rules.NewQuota(map[string]int{"water": 5})
	if e := q.Consume("water", 3); e != nil {
		t.Fatal(e)
	}
	if q.Remaining("water") != 2 {
		t.Fatal(q.Remaining("water"))
	}
	if e := q.Consume("water", 3); !errors.Is(e, common.ErrConflict) {
		t.Fatal(e)
	}
	if e := q.Restore("water", 2); e != nil {
		t.Fatal(e)
	}
	if e := q.Consume("unknown", 1); !errors.Is(e, common.ErrInvalid) {
		t.Fatal(e)
	}
}
func TestAssignments(t *testing.T) {
	a := rules.Assignment{ID: "a", Resource: "van", Owner: "team", Amount: 2}
	if e := a.Activate(); e != nil {
		t.Fatal(e)
	}
	if e := a.Activate(); !errors.Is(e, common.ErrConflict) {
		t.Fatal(e)
	}
	if e := a.Release(); e != nil {
		t.Fatal(e)
	}
	if e := a.Release(); !errors.Is(e, common.ErrConflict) {
		t.Fatal(e)
	}
	m := rules.Merge([]rules.Assignment{a})
	if len(m) != 1 {
		t.Fatal(m)
	}
}
func TestValidationHelpers(t *testing.T) {
	for _, v := range []string{"FJ-01", "CN-100"} {
		if e := rules.ValidateZoneCode(v); e != nil {
			t.Fatal(e)
		}
	}
	for _, v := range []string{"fj-1", "FJ", ""} {
		if e := rules.ValidateZoneCode(v); e == nil {
			t.Fatal(v)
		}
	}
	if !rules.ValidateUnique([]string{"a", "b"}) || rules.ValidateUnique([]string{"a", "a"}) {
		t.Fatal()
	}
	got := rules.NormalizeTags([]string{" Water ", "water", "Food"})
	if len(got) != 2 {
		t.Fatal(got)
	}
}
func TestEventBus(t *testing.T) {
	b := rules.NewBus()
	n := 0
	b.On("created", func(rules.Event) { n++ })
	if e := b.Emit(rules.Event{Type: "created", Entity: "request"}); e != nil {
		t.Fatal(e)
	}
	if n != 1 || b.Count("created") != 1 {
		t.Fatal(n)
	}
}
func TestSummary(t *testing.T) {
	s := rules.Summary{Zone: "z", Requests: 10, Allocated: 8, Delivered: 6}
	if s.Rate() != .6 || !s.Healthy() {
		t.Fatal(s)
	}
	s2 := rules.Summary{Zone: "z", Requests: 2, Allocated: 2, Delivered: 2}
	m := rules.MergeSummaries([]rules.Summary{s, s2})
	if m.Requests != 12 || m.Delivered != 8 {
		t.Fatal(m)
	}
}
