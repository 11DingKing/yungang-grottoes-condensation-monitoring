package aggregate

import (
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/common"
	"strings"
	"time"
)

type Incident struct {
	ID, PlanID, Title, Status string
	Severity                  int
	Notes                     []string
	UpdatedAt                 time.Time
}

func (i *Incident) Validate() error {
	if i.ID == "" || i.PlanID == "" || strings.TrimSpace(i.Title) == "" || i.Severity < 1 || i.Severity > 5 {
		return common.ErrInvalid
	}
	return nil
}
func (i *Incident) Assign(note string) error {
	if i.Status != "open" || note == "" {
		return common.ErrConflict
	}
	i.Status = "assigned"
	i.Notes = append(i.Notes, note)
	i.UpdatedAt = common.Now()
	return nil
}
func (i *Incident) Resolve(note string) error {
	if i.Status != "assigned" || note == "" {
		return common.ErrConflict
	}
	i.Status = "resolved"
	i.Notes = append(i.Notes, note)
	i.UpdatedAt = common.Now()
	return nil
}
func (i *Incident) Close(note string) error {
	if i.Status != "resolved" {
		return common.ErrConflict
	}
	i.Status = "closed"
	i.Notes = append(i.Notes, note)
	i.UpdatedAt = common.Now()
	return nil
}
func (i Incident) Open() bool { return i.Status == "open" || i.Status == "assigned" }
