package planning_test

import (
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/common"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/planning"
	"testing"
	"time"
)

func TestTaskFlow(t *testing.T) {
	now := time.Now()
	task := planning.Task{ID: "t1", Name: "evacuate", Zone: "z", State: planning.Queued, Priority: 5}
	if e := task.Start(now); e != nil {
		t.Fatal(e)
	}
	if task.Attempts != 1 || task.State != planning.Running {
		t.Fatal(task)
	}
	if e := task.Finish(now.Add(time.Minute)); e != nil {
		t.Fatal(e)
	}
	if task.FinishedAt == nil {
		t.Fatal(task)
	}
}
func TestTaskDependencies(t *testing.T) {
	task := planning.Task{ID: "t", Name: "load", Zone: "z", Dependencies: []string{"a", "b"}}
	if task.Ready(map[string]bool{"a": true}) {
		t.Fatal("partial ready")
	}
	if !task.Ready(map[string]bool{"a": true, "b": true}) {
		t.Fatal("ready rejected")
	}
}
func TestSchedule(t *testing.T) {
	items := []planning.Task{{ID: "b", Name: "b", Zone: "z", Priority: 2}, {ID: "a", Name: "a", Zone: "z", Priority: 5}}
	got := planning.Schedule(items)
	if got[0].ID != "a" {
		t.Fatal(got)
	}
}
func TestRouteDistance(t *testing.T) {
	a := planning.Point{Lat: 24.5, Lon: 118.1}
	b := planning.Point{Lat: 24.6, Lon: 118.2}
	if planning.Distance(a, b) <= 0 {
		t.Fatal()
	}
	if planning.Distance(planning.Point{Lat: 100}, b) != 0 {
		t.Fatal()
	}
}
func TestStops(t *testing.T) {
	stops := []planning.Stop{{ID: "a", Point: planning.Point{Lat: 24, Lon: 118}, Load: 2, Priority: 1}, {ID: "b", Point: planning.Point{Lat: 25, Lon: 119}, Load: 3, Priority: 5}}
	if e := planning.ValidateStops(stops); e != nil {
		t.Fatal(e)
	}
	if planning.TotalLoad(stops) != 5 || planning.SortStops(stops)[0].ID != "b" {
		t.Fatal(stops)
	}
}
func TestResource(t *testing.T) {
	r := planning.Resource{ID: "r", Kind: "truck", Capacity: 10}
	if e := r.Allocate(4); e != nil {
		t.Fatal(e)
	}
	if r.Available() != 6 {
		t.Fatal(r.Available())
	}
	if e := r.Allocate(7); e == nil {
		t.Fatal("overallocated")
	}
	if e := r.Free(2); e != nil {
		t.Fatal(e)
	}
}
func TestSlots(t *testing.T) {
	now := time.Now()
	a := planning.Slot{Start: now, End: now.Add(time.Hour), Label: "morning", Capacity: 2}
	if !a.Valid() || !a.Contains(now.Add(time.Minute)) {
		t.Fatal(a)
	}
	b := planning.Slot{Start: now.Add(30 * time.Minute), End: now.Add(2 * time.Hour), Label: "overlap", Capacity: 1}
	if !planning.Overlaps(a, b) {
		t.Fatal()
	}
	if len(planning.Available([]planning.Slot{a, b}, now.Add(45*time.Minute))) != 2 {
		t.Fatal()
	}
}
func TestAggregate(t *testing.T) {
	items := []planning.Aggregate{{Zone: "a", Tasks: 10, Completed: 5, Blocked: 1}, {Zone: "b", Tasks: 2, Completed: 2}}
	m := planning.Combine(items)
	if m.Tasks != 12 || m.Completed != 7 {
		t.Fatal(m)
	}
	if m.Actionable() != 4 {
		t.Fatal(m.Actionable())
	}
	if m.Progress() <= 0 {
		t.Fatal()
	}
	_ = common.ErrInvalid
}
