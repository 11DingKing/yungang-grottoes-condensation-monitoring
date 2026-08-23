package rules

import (
	"fmt"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/common"
	"sort"
	"time"
)

type Priority int

const (
	Routine   Priority = 1
	Elevated  Priority = 2
	Urgent    Priority = 3
	Critical  Priority = 4
	Emergency Priority = 5
)

func (p Priority) Valid() bool { return p >= Routine && p <= Emergency }
func (p Priority) Name() string {
	switch p {
	case Routine:
		return "routine"
	case Elevated:
		return "elevated"
	case Urgent:
		return "urgent"
	case Critical:
		return "critical"
	case Emergency:
		return "emergency"
	default:
		return "unknown"
	}
}
func ParsePriority(v int) (Priority, error) {
	p := Priority(v)
	if !p.Valid() {
		return 0, fmt.Errorf("%w: priority %d", common.ErrInvalid, v)
	}
	return p, nil
}

type Deadline struct {
	Priority     Priority
	Started, Due time.Time
}

func NewDeadline(p Priority, start time.Time) Deadline {
	return Deadline{p, start, start.Add(time.Duration(6-p) * time.Hour)}
}
func (d Deadline) Expired(now time.Time) bool { return !now.Before(d.Due) }
func (d Deadline) Remaining(now time.Time) time.Duration {
	if d.Expired(now) {
		return 0
	}
	return d.Due.Sub(now)
}
func SortByUrgency(items []Deadline) []Deadline {
	out := append([]Deadline(nil), items...)
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Priority == out[j].Priority {
			return out[i].Due.Before(out[j].Due)
		}
		return out[i].Priority > out[j].Priority
	})
	return out
}

type Window struct{ Start, End time.Time }

func (w Window) Contains(t time.Time) bool { return !t.Before(w.Start) && t.Before(w.End) }
func (w Window) Duration() time.Duration   { return w.End.Sub(w.Start) }
func (w Window) Valid() bool               { return w.End.After(w.Start) }
