package aggregate

import (
	"fmt"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/common"
	"time"
)

type Plan struct {
	ID, Zone, Name string
	Level          int
	Status         string
	Start, End     *time.Time
	Version        int
}

func (p Plan) Validate() error {
	if p.ID == "" || p.Zone == "" || p.Name == "" || p.Level < 1 || p.Level > 4 {
		return common.ErrInvalid
	}
	return nil
}
func (p *Plan) Activate(now time.Time) error {
	if e := p.Validate(); e != nil {
		return e
	}
	if p.Status != "draft" {
		return common.ErrConflict
	}
	p.Status = "active"
	p.Start = &now
	p.Version++
	return nil
}
func (p *Plan) Standby() error {
	if p.Status != "active" {
		return common.ErrConflict
	}
	p.Status = "standby"
	p.Version++
	return nil
}
func (p *Plan) Close(now time.Time) error {
	if p.Status != "active" && p.Status != "standby" {
		return fmt.Errorf("%w: plan state", common.ErrConflict)
	}
	p.Status = "closed"
	p.End = &now
	p.Version++
	return nil
}
func (p Plan) Running(now time.Time) bool {
	return p.Start != nil && (p.End == nil || now.Before(*p.End)) && p.Status != "closed"
}
