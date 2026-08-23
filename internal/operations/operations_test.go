package operations_test

import (
	"context"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/operations"
	"testing"
	"time"
)

func TestChecklist(t *testing.T) {
	c := operations.NewChecklist([]operations.CheckItem{{ID: "a", Required: true}, {ID: "b", Required: false}})
	if c.Ready() {
		t.Fatal()
	}
	if e := c.Complete("a"); e != nil {
		t.Fatal(e)
	}
	if e := c.RequireReady(); e != nil {
		t.Fatal(e)
	}
	if len(c.Pending()) != 1 {
		t.Fatal(c.Pending())
	}
}
func TestApproval(t *testing.T) {
	now := time.Now()
	a := operations.Approval{ID: "a", Subject: "plan"}
	if e := a.Request(now); e != nil {
		t.Fatal(e)
	}
	if e := a.Approve("u", "ok", now); e != nil || !a.Final() {
		t.Fatal(e, a)
	}
}
func TestClaim(t *testing.T) {
	now := time.Now()
	c := operations.Claim{ID: "c"}
	if e := c.Acquire("u", time.Minute, now); e != nil {
		t.Fatal(e)
	}
	if !c.Held(now) {
		t.Fatal()
	}
	if e := c.Release("u"); e != nil {
		t.Fatal(e)
	}
}
func TestQueuePriority(t *testing.T) {
	q := operations.NewQueue()
	q.Add(operations.Item{ID: "low", Priority: 1, At: time.Now()})
	q.Add(operations.Item{ID: "high", Priority: 5, At: time.Now()})
	v, ok := q.Next()
	if !ok || v.ID != "high" {
		t.Fatal(v, ok)
	}
	if q.Len() != 1 {
		t.Fatal(q.Len())
	}
}
func TestBackoff(t *testing.T) {
	b := operations.Backoff{Base: time.Millisecond, Max: time.Second}
	if b.Next() != time.Millisecond || b.Attempts != 1 {
		t.Fatal(b)
	}
	if b.Delay(10) > time.Second {
		t.Fatal()
	}
	b.Reset()
	if b.Attempts != 0 {
		t.Fatal()
	}
}
func TestWindows(t *testing.T) {
	now := time.Now()
	a := operations.Window{Start: now, End: now.Add(time.Hour), Resource: "van"}
	b := operations.Window{Start: now.Add(2 * time.Hour), End: now.Add(3 * time.Hour), Resource: "van"}
	if !operations.FindAvailable([]operations.Window{a}, b) {
		t.Fatal()
	}
	if operations.FindAvailable([]operations.Window{a}, operations.Window{Start: now.Add(30 * time.Minute), End: now.Add(90 * time.Minute), Resource: "van"}) {
		t.Fatal()
	}
}
func TestReport(t *testing.T) {
	r := operations.Report{Title: "daily"}
	now := time.Now()
	r.Add("requests", 2, now)
	r.Add("requests", 4, now.Add(time.Hour))
	if r.Average("requests") != 3 {
		t.Fatal(r.Average("requests"))
	}
	if _, ok := r.Latest("missing"); ok {
		t.Fatal()
	}
}
func TestBalance(t *testing.T) {
	b := operations.NewBalance()
	b.Credit("water", 5)
	if !b.Debit("water", 2) || b.Get("water") != 3 {
		t.Fatal(b.Get("water"))
	}
	if b.Debit("water", 4) {
		t.Fatal()
	}
	operations.ApplyEntries(b, []operations.Entry{{Key: "water", Delta: -1}})
	if b.Get("water") != 2 {
		t.Fatal(b.Get("water"))
	}
}
func TestLockableClaim(t *testing.T) { ctx := context.Background(); _ = ctx }
