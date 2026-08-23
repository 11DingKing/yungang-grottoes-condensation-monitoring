package snapshot

import (
	"fmt"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/common"
	"time"
)

type Policy struct {
	Retention       time.Duration
	MaxSize         int
	RequireChecksum bool
}

func (p Policy) Validate() error {
	if p.Retention <= 0 || p.MaxSize < 1 {
		return common.ErrInvalid
	}
	return nil
}
func (p Policy) Accept(a Artifact, now time.Time) error {
	if a.ID == "" || a.Size < 0 {
		return common.ErrInvalid
	}
	if a.Size > p.MaxSize {
		return fmt.Errorf("%w: artifact size", common.ErrConflict)
	}
	if now.Sub(a.CreatedAt) > p.Retention && !a.Pinned {
		return fmt.Errorf("%w: artifact expired", common.ErrConflict)
	}
	return nil
}
func (p Policy) ShouldDelete(a Artifact, now time.Time) bool {
	return !a.Pinned && now.Sub(a.CreatedAt) > p.Retention
}
func (p Policy) Limit(items []Artifact) []Artifact {
	out := Keep(items, 100)
	total := 0
	result := []Artifact{}
	for _, a := range out {
		if total+a.Size > p.MaxSize {
			continue
		}
		result = append(result, a)
		total += a.Size
	}
	return result
}
