package aggregate

import (
	"sort"
	"time"
)

type Event struct {
	Type, ID, Actor string
	At              time.Time
	Detail          string
}

func SortEvents(events []Event) []Event {
	out := append([]Event(nil), events...)
	sort.SliceStable(out, func(i, j int) bool { return out[i].At.Before(out[j].At) })
	return out
}
func GroupEvents(events []Event) map[string][]Event {
	out := map[string][]Event{}
	for _, e := range events {
		out[e.ID] = append(out[e.ID], e)
	}
	for k, v := range out {
		out[k] = SortEvents(v)
	}
	return out
}
func Latest(events []Event) (Event, bool) {
	if len(events) == 0 {
		return Event{}, false
	}
	out := events[0]
	for _, e := range events[1:] {
		if e.At.After(out.At) {
			out = e
		}
	}
	return out, true
}
