package retry

import (
	"context"
	"fmt"
	"math"
	"time"
)

type Policy struct {
	Max  int
	Base time.Duration
	Cap  time.Duration
}

func Default() Policy { return Policy{Max: 5, Base: time.Second, Cap: 5 * time.Minute} }
func (p Policy) Delay(attempt int) time.Duration {
	if attempt < 1 {
		return 0
	}
	d := time.Duration(float64(p.Base) * math.Pow(2, float64(attempt-1)))
	if d > p.Cap {
		return p.Cap
	}
	return d
}
func (p Policy) Run(ctx context.Context, fn func(context.Context, int) error) error {
	var last error
	for i := 1; i <= p.Max; i++ {
		if e := ctx.Err(); e != nil {
			return e
		}
		e := fn(ctx, i)
		if e == nil {
			return nil
		}
		last = e
		if i < p.Max {
			t := time.NewTimer(p.Delay(i))
			select {
			case <-ctx.Done():
				t.Stop()
				return ctx.Err()
			case <-t.C:
			}
		}
	}
	return fmt.Errorf("retry exhausted: %w", last)
}
