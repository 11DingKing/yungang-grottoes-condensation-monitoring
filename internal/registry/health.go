package registry

import (
	"context"
	"time"
)

type Check func(context.Context) error
type Health struct {
	Name    string
	Healthy bool
	Latency time.Duration
	Error   string
}

func Run(ctx context.Context, name string, check Check) Health {
	start := time.Now()
	e := check(ctx)
	h := Health{Name: name, Healthy: e == nil, Latency: time.Since(start)}
	if e != nil {
		h.Error = e.Error()
	}
	return h
}
func All(ctx context.Context, checks map[string]Check) []Health {
	out := []Health{}
	for name, check := range checks {
		out = append(out, Run(ctx, name, check))
	}
	return out
}
