package planning

import (
	"fmt"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/common"
	"time"
)

type Constraint struct {
	Name  string
	Limit int
	Used  int
}

func (c Constraint) Remaining() int { return c.Limit - c.Used }
func (c *Constraint) Use(n int) error {
	if n <= 0 {
		return common.ErrInvalid
	}
	if c.Used+n > c.Limit {
		return fmt.Errorf("%w: %s", common.ErrConflict, c.Name)
	}
	c.Used += n
	return nil
}
func (c *Constraint) Release(n int) error {
	if n <= 0 || n > c.Used {
		return common.ErrInvalid
	}
	c.Used -= n
	return nil
}

type Interval struct {
	Start, End time.Time
	Capacity   int
}

func (i Interval) Valid() bool               { return i.End.After(i.Start) && i.Capacity > 0 }
func (i Interval) Includes(t time.Time) bool { return !t.Before(i.Start) && t.Before(i.End) }
func AvailableIntervals(items []Interval, t time.Time) []Interval {
	out := []Interval{}
	for _, i := range items {
		if i.Includes(t) {
			out = append(out, i)
		}
	}
	return out
}
