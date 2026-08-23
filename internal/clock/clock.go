package clock

import "time"

type Clock interface {
	Now() time.Time
	After(time.Time) bool
}
type Real struct{}

func (Real) Now() time.Time         { return time.Now() }
func (Real) After(t time.Time) bool { return time.Now().After(t) }

type Fixed struct{ Value time.Time }

func (f Fixed) Now() time.Time         { return f.Value }
func (f Fixed) After(t time.Time) bool { return f.Value.After(t) }
func UTC(c Clock) time.Time            { return c.Now().UTC() }
func InWindow(c Clock, start, end time.Time) bool {
	n := c.Now()
	return !n.Before(start) && n.Before(end)
}
func Deadline(c Clock, d time.Duration) time.Time { return c.Now().Add(d) }
