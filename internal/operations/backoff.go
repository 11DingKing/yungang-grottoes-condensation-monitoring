package operations

import (
	"math"
	"time"
)

type Backoff struct {
	Base, Max time.Duration
	Attempts  int
}

func (b Backoff) Delay(attempt int) time.Duration {
	if attempt < 1 {
		return 0
	}
	d := time.Duration(float64(b.Base) * math.Pow(2, float64(attempt-1)))
	if d > b.Max {
		return b.Max
	}
	return d
}
func (b *Backoff) Next() time.Duration     { b.Attempts++; return b.Delay(b.Attempts) }
func (b *Backoff) Reset()                  { b.Attempts = 0 }
func (b Backoff) Exhausted(limit int) bool { return b.Attempts >= limit }
