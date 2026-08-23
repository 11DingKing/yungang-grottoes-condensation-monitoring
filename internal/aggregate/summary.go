package aggregate

import (
	"math"
	"time"
)

type ZoneSummary struct {
	Zone                                      string
	OpenPlans, OpenIncidents, PendingRequests int
	UpdatedAt                                 time.Time
}

func (s ZoneSummary) Risk() float64 {
	return math.Min(1, float64(s.OpenIncidents*2+s.PendingRequests)/10)
}
func (s ZoneSummary) Stable() bool { return s.OpenPlans > 0 && s.OpenIncidents == 0 }
func MergeSummaries(items []ZoneSummary) ZoneSummary {
	out := ZoneSummary{}
	for _, s := range items {
		if out.Zone == "" {
			out.Zone = s.Zone
		}
		out.OpenPlans += s.OpenPlans
		out.OpenIncidents += s.OpenIncidents
		out.PendingRequests += s.PendingRequests
		if s.UpdatedAt.After(out.UpdatedAt) {
			out.UpdatedAt = s.UpdatedAt
		}
	}
	return out
}
