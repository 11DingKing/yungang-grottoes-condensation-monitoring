package rules

import (
	"fmt"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/common"
)

type Check struct {
	Name    string
	Test    func() bool
	Message string
}
type Report struct {
	Passed []string
	Failed []string
}

func RunChecks(checks []Check) Report {
	r := Report{}
	for _, c := range checks {
		if c.Test != nil && c.Test() {
			r.Passed = append(r.Passed, c.Name)
		} else {
			r.Failed = append(r.Failed, c.Name)
		}
	}
	return r
}
func (r Report) OK() bool { return len(r.Failed) == 0 }
func (r Report) Error() error {
	if r.OK() {
		return nil
	}
	return fmt.Errorf("%w: %v", common.ErrConflict, r.Failed)
}
func RequireAll(checks ...Check) error { return RunChecks(checks).Error() }
