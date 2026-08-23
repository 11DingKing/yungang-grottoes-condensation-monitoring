package operations

import (
	"sort"
	"time"
)

type Window struct {
	Start, End time.Time
	Resource   string
}

func (w Window) Valid() bool { return w.Resource != "" && w.End.After(w.Start) }
func (w Window) Overlaps(other Window) bool {
	return w.Resource == other.Resource && w.Start.Before(other.End) && other.Start.Before(w.End)
}
func SortWindows(items []Window) []Window {
	out := append([]Window(nil), items...)
	sort.Slice(out, func(i, j int) bool { return out[i].Start.Before(out[j].Start) })
	return out
}
func FindAvailable(items []Window, w Window) bool {
	for _, v := range items {
		if v.Overlaps(w) {
			return false
		}
	}
	return true
}
