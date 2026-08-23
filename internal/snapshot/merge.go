package snapshot

import (
	"sort"
	"time"
)

type Change struct {
	ID      string
	Version int
	At      time.Time
	Deleted bool
	Payload []byte
}

func SortChanges(changes []Change) []Change {
	out := append([]Change(nil), changes...)
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].At.Equal(out[j].At) {
			return out[i].Version < out[j].Version
		}
		return out[i].At.Before(out[j].At)
	})
	return out
}
func LatestChanges(changes []Change) map[string]Change {
	out := map[string]Change{}
	for _, c := range SortChanges(changes) {
		out[c.ID] = c
	}
	return out
}
func ActiveChanges(changes []Change) []Change {
	out := []Change{}
	for _, c := range LatestChanges(changes) {
		if !c.Deleted {
			out = append(out, c)
		}
	}
	return out
}
func CloneChanges(changes []Change) []Change {
	out := make([]Change, len(changes))
	for i, c := range changes {
		c.Payload = append([]byte(nil), c.Payload...)
		out[i] = c
	}
	return out
}
