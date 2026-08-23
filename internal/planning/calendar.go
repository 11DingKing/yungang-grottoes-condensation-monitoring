package planning

import (
	"sort"
	"time"
)

type Slot struct {
	Start, End time.Time
	Label      string
	Capacity   int
}

func (s Slot) Valid() bool               { return s.End.After(s.Start) && s.Label != "" && s.Capacity > 0 }
func (s Slot) Contains(t time.Time) bool { return !t.Before(s.Start) && t.Before(s.End) }
func SortSlots(slots []Slot) []Slot {
	out := append([]Slot(nil), slots...)
	sort.Slice(out, func(i, j int) bool { return out[i].Start.Before(out[j].Start) })
	return out
}
func Overlaps(a, b Slot) bool { return a.Start.Before(b.End) && b.Start.Before(a.End) }
func Available(slots []Slot, t time.Time) []Slot {
	out := []Slot{}
	for _, s := range slots {
		if s.Contains(t) {
			out = append(out, s)
		}
	}
	return out
}
