package aggregate_test

import (
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/aggregate"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/common"
	"testing"
	"time"
)

func TestRequestAggregate(t *testing.T) {
	r := aggregate.NewRequest("r", "z", "hall")
	if e := r.AddItem("water", 3); e != nil {
		t.Fatal(e)
	}
	if e := r.AddItem("food", 2); e != nil {
		t.Fatal(e)
	}
	if r.Total() != 5 || len(r.SKUs()) != 2 {
		t.Fatal(r)
	}
	if e := r.Approve(); e != nil {
		t.Fatal(e)
	}
	if e := r.AddItem("blanket", 1); e == nil {
		t.Fatal("approved request modified")
	}
}
func TestDispatchAggregate(t *testing.T) {
	d := aggregate.Dispatch{ID: "d", RequestID: "r", Status: "queued"}
	if e := d.Assign("driver", "van"); e != nil {
		t.Fatal(e)
	}
	if e := d.Start(); e != nil {
		t.Fatal(e)
	}
	if e := d.Deliver(); e != nil {
		t.Fatal(e)
	}
	if len(d.Events) != 3 {
		t.Fatal(d.Events)
	}
	if e := d.Fail("late"); e == nil {
		t.Fatal("delivered failed")
	}
}
func TestIncidentAggregate(t *testing.T) {
	i := aggregate.Incident{ID: "i", PlanID: "p", Title: "water", Severity: 3, Status: "open"}
	if e := i.Validate(); e != nil {
		t.Fatal(e)
	}
	if e := i.Assign("team"); e != nil {
		t.Fatal(e)
	}
	if e := i.Resolve("fixed"); e != nil {
		t.Fatal(e)
	}
	if e := i.Close("verified"); e != nil {
		t.Fatal(e)
	}
	if i.Open() {
		t.Fatal(i)
	}
}
func TestPlanAggregate(t *testing.T) {
	now := time.Now()
	p := aggregate.Plan{ID: "p", Zone: "z", Name: "plan", Level: 2, Status: "draft", Version: 1}
	if e := p.Activate(now); e != nil {
		t.Fatal(e)
	}
	if !p.Running(now) {
		t.Fatal()
	}
	if e := p.Standby(); e != nil {
		t.Fatal(e)
	}
	if e := p.Close(now); e != nil {
		t.Fatal(e)
	}
	if p.Running(now) {
		t.Fatal()
	}
}
func TestEventGrouping(t *testing.T) {
	now := time.Now()
	events := []aggregate.Event{{ID: "r", At: now.Add(time.Minute)}, {ID: "r", At: now}, {ID: "d", At: now}}
	groups := aggregate.GroupEvents(events)
	if len(groups["r"]) != 2 || !groups["r"][0].At.Before(groups["r"][1].At) {
		t.Fatal(groups)
	}
	if _, ok := aggregate.Latest(nil); ok {
		t.Fatal()
	}
}
func TestValidationHelpers(t *testing.T) {
	if e := aggregate.Required("a", "b"); e != nil {
		t.Fatal(e)
	}
	if aggregate.ValidID("x") {
		t.Fatal()
	}
	if !aggregate.ValidStatus("open", "open", "closed") {
		t.Fatal()
	}
	a := map[string]int{"water": 2}
	b := aggregate.CloneItems(a)
	b["water"] = 4
	if a["water"] != 2 || aggregate.Sum(a) != 2 {
		t.Fatal(a, b)
	}
	_ = common.ErrInvalid
}
