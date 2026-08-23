package plan

import (
	"fmt"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/common"
	"time"
)

type Status string

const (
	Draft   Status = "draft"
	Active  Status = "active"
	Standby Status = "standby"
	Closed  Status = "closed"
)

type Plan struct {
	ID, ZoneID, Name    string
	Level               int
	Status              Status
	Version             int
	StartsAt, UpdatedAt time.Time
	EndsAt              *time.Time
}

func New(zoneID, name string, level int) (Plan, error) {
	if zoneID == "" || name == "" || level < 1 || level > 4 {
		return Plan{}, common.ErrInvalid
	}
	now := common.Now()
	return Plan{ID: common.NewID("plan"), ZoneID: zoneID, Name: name, Level: level, Status: Draft, Version: 1, StartsAt: now, UpdatedAt: now}, nil
}
func (p *Plan) Transition(next Status) error {
	valid := map[Status]map[Status]bool{Draft: {Active: true}, Active: {Standby: true, Closed: true}, Standby: {Closed: true}}
	if !valid[p.Status][next] {
		return fmt.Errorf("%w: plan transition %s to %s", common.ErrConflict, p.Status, next)
	}
	p.Status = next
	p.Version++
	p.UpdatedAt = common.Now()
	if next == Closed {
		now := common.Now()
		p.EndsAt = &now
	}
	return nil
}
