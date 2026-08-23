package utility_test

import (
	"context"
	"errors"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/clock"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/filter"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/metrics"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/notification"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/pagination"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/retry"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/tenant"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/transport"
	"net/url"
	"testing"
	"time"
)

func TestPagination(t *testing.T) {
	q := pagination.Parse(url.Values{"page": {"2"}, "size": {"200"}, "direction": {"DESC"}})
	if q.Page != 2 || q.Size != 100 || q.Direction != "desc" || q.Offset() != 100 {
		t.Fatal(q)
	}
	if !q.ValidSort(map[string]bool{"created_at": true}) {
		t.Fatal(q)
	}
}
func TestFiltersCopyInput(t *testing.T) {
	in := []int{4, 1, 3, 2}
	out := filter.Page(filter.SortBy(in, func(a, b int) bool { return a < b }), 2, 2)
	if len(out) != 2 || out[0] != 3 {
		t.Fatal(out)
	}
	if in[0] != 4 {
		t.Fatal("input mutated")
	}
}
func TestClockWindow(t *testing.T) {
	now := time.Date(2026, 8, 23, 10, 0, 0, 0, time.UTC)
	c := clock.Fixed{Value: now}
	if !clock.InWindow(c, now.Add(-time.Minute), now.Add(time.Minute)) {
		t.Fatal("outside")
	}
	if clock.InWindow(c, now.Add(time.Minute), now.Add(2*time.Minute)) {
		t.Fatal("inside")
	}
}
func TestCounterSnapshotIsolated(t *testing.T) {
	c := metrics.NewCounter()
	c.Add("requests", 2)
	m := c.Snapshot()
	m["requests"] = 9
	if c.Get("requests") != 2 {
		t.Fatal(c.Get("requests"))
	}
}
func TestNotificationContext(t *testing.T) {
	sink := &notification.MemorySink{}
	s := notification.Service{Sink: sink}
	if e := s.Notify(context.Background(), "zone", "alert", "one", "two"); e != nil {
		t.Fatal(e)
	}
	if len(sink.Snapshot()) != 1 {
		t.Fatal(sink.Snapshot())
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if e := sink.Send(ctx, notification.Message{}); !errors.Is(e, context.Canceled) {
		t.Fatal(e)
	}
}
func TestRetryContext(t *testing.T) {
	p := retry.Policy{Max: 3, Base: time.Millisecond, Cap: time.Millisecond}
	n := 0
	e := p.Run(context.Background(), func(context.Context, int) error {
		n++
		if n == 3 {
			return nil
		}
		return errors.New("try")
	})
	if e != nil || n != 3 {
		t.Fatal(e, n)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if e = p.Run(ctx, func(context.Context, int) error { return errors.New("no") }); !errors.Is(e, context.Canceled) {
		t.Fatal(e)
	}
}
func TestTransportValidation(t *testing.T) {
	if e := (transport.LoginRequest{Username: "ab", Password: "x"}).Validate(); e == nil {
		t.Fatal("accepted")
	}
	if e := (transport.PlanRequest{Name: "plan", ZoneID: "z", Level: 2}).Validate(); e != nil {
		t.Fatal(e)
	}
	if e := (transport.ReliefRequest{ZoneID: "z", Destination: "hall", Priority: 2, Items: []transport.RequestItem{{SupplyID: "s", Quantity: 1}}}).Validate(); e != nil {
		t.Fatal(e)
	}
	if transport.SafeNext(url.Values{"next": {"https://evil"}}) != "/" {
		t.Fatal("unsafe redirect")
	}
}
func TestTenantScope(t *testing.T) {
	ctx := tenant.With(context.Background(), "z1")
	if e := tenant.Require(ctx, "z1"); e != nil {
		t.Fatal(e)
	}
	if e := tenant.Require(ctx, "z2"); e == nil {
		t.Fatal("cross zone")
	}
}
