package rules

import (
	"sort"
	"time"
)

type Summary struct {
	Zone                           string
	Requests, Allocated, Delivered int
	UpdatedAt                      time.Time
}

func (s Summary) Rate() float64 {
	if s.Requests == 0 {
		return 0
	}
	return float64(s.Delivered) / float64(s.Requests)
}
func (s Summary) Healthy() bool {
	return s.Requests > 0 && s.Allocated <= s.Requests && s.Delivered <= s.Allocated
}
func SortSummaries(items []Summary) []Summary {
	out := append([]Summary(nil), items...)
	sort.Slice(out, func(i, j int) bool { return out[i].UpdatedAt.After(out[j].UpdatedAt) })
	return out
}
func MergeSummaries(items []Summary) Summary {
	out := Summary{}
	for _, s := range items {
		if out.Zone == "" {
			out.Zone = s.Zone
		}
		out.Requests += s.Requests
		out.Allocated += s.Allocated
		out.Delivered += s.Delivered
		if s.UpdatedAt.After(out.UpdatedAt) {
			out.UpdatedAt = s.UpdatedAt
		}
	}
	return out
}
