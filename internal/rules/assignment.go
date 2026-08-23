package rules

import (
	"fmt"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/common"
	"strings"
)

type Assignment struct {
	ID, Resource, Owner string
	Amount              int
	Active              bool
}

func (a Assignment) Validate() error {
	if a.ID == "" || strings.TrimSpace(a.Resource) == "" || strings.TrimSpace(a.Owner) == "" || a.Amount <= 0 {
		return common.ErrInvalid
	}
	return nil
}
func (a *Assignment) Activate() error {
	if e := a.Validate(); e != nil {
		return e
	}
	if a.Active {
		return fmt.Errorf("%w: already active", common.ErrConflict)
	}
	a.Active = true
	return nil
}
func (a *Assignment) Release() error {
	if !a.Active {
		return fmt.Errorf("%w: not active", common.ErrConflict)
	}
	a.Active = false
	return nil
}
func (a Assignment) SameResource(other Assignment) bool {
	return a.Resource == other.Resource && a.Active && other.Active
}
func Merge(items []Assignment) map[string]Assignment {
	out := map[string]Assignment{}
	for _, a := range items {
		if old, ok := out[a.ID]; ok && old.Active != a.Active {
			continue
		}
		out[a.ID] = a
	}
	return out
}
