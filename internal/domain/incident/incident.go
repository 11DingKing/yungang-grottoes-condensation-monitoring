package incident

import (
	"fmt"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/common"
)

type Status string

const (
	Open     Status = "open"
	Assigned Status = "assigned"
	Resolved Status = "resolved"
	Closed   Status = "closed"
)

type Incident struct {
	ID, PlanID, ReporterID, Title, ContextNote string
	Severity                                   int
	Status                                     Status
}

func (i *Incident) Transition(next Status) error {
	valid := map[Status]map[Status]bool{Open: {Assigned: true}, Assigned: {Resolved: true}, Resolved: {Closed: true}}
	if !valid[i.Status][next] {
		return fmt.Errorf("%w: incident transition", common.ErrConflict)
	}
	i.Status = next
	return nil
}
