package operations

import (
	"sort"
	"time"
)

type Metric struct {
	Key   string
	Value float64
	At    time.Time
}
type Report struct {
	Title       string
	Metrics     []Metric
	GeneratedAt time.Time
}

func (r *Report) Add(key string, value float64, at time.Time) {
	r.Metrics = append(r.Metrics, Metric{key, value, at})
}
func (r Report) Sorted() []Metric {
	out := append([]Metric(nil), r.Metrics...)
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Key == out[j].Key {
			return out[i].At.Before(out[j].At)
		}
		return out[i].Key < out[j].Key
	})
	return out
}
func (r Report) Latest(key string) (Metric, bool) {
	var out Metric
	ok := false
	for _, m := range r.Metrics {
		if m.Key == key && (!ok || m.At.After(out.At)) {
			out = m
			ok = true
		}
	}
	return out, ok
}
func (r Report) Average(key string) float64 {
	var total float64
	n := 0
	for _, m := range r.Metrics {
		if m.Key == key {
			total += m.Value
			n++
		}
	}
	if n == 0 {
		return 0
	}
	return total / float64(n)
}
