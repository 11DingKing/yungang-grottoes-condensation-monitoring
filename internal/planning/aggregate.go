package planning

import (
	"sort"
	"time"
)

type Aggregate struct {
	Zone      string
	Tasks     int
	Completed int
	Blocked   int
	At        time.Time
}

func (a Aggregate) Progress() float64 {
	if a.Tasks == 0 {
		return 0
	}
	return float64(a.Completed) / float64(a.Tasks)
}
func (a Aggregate) Actionable() int {
	n := a.Tasks - a.Completed - a.Blocked
	if n < 0 {
		return 0
	}
	return n
}
func Combine(items []Aggregate) Aggregate {
	out := Aggregate{}
	for _, a := range items {
		if out.Zone == "" {
			out.Zone = a.Zone
		}
		out.Tasks += a.Tasks
		out.Completed += a.Completed
		out.Blocked += a.Blocked
		if a.At.After(out.At) {
			out.At = a.At
		}
	}
	return out
}
func SortAggregates(items []Aggregate) []Aggregate {
	out := append([]Aggregate(nil), items...)
	sort.Slice(out, func(i, j int) bool { return out[i].Progress() > out[j].Progress() })
	return out
}
