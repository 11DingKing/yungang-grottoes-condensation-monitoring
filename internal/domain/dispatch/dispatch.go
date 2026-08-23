package dispatch

import (
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/common"
	"time"
)

type Status string

const (
	Queued    Status = "queued"
	InTransit Status = "in_transit"
	Delivered Status = "delivered"
	Failed    Status = "failed"
)

type Dispatch struct {
	ID, RequestID, Driver, Vehicle string
	Status                         Status
	Attempt                        int
	NextAttemptAt                  time.Time
	DeliveredAt                    *time.Time
	Version                        int
}

func (d *Dispatch) Start() error {
	if d.Status != Queued {
		return common.ErrConflict
	}
	d.Status = InTransit
	d.Version++
	return nil
}
func (d *Dispatch) Deliver() error {
	if d.Status != InTransit {
		return common.ErrConflict
	}
	now := common.Now()
	d.Status = Delivered
	d.DeliveredAt = &now
	d.Version++
	return nil
}
func (d *Dispatch) Retry(reason string) error {
	if d.Status == Delivered {
		return common.ErrConflict
	}
	d.Attempt++
	if d.Attempt >= 5 {
		d.Status = Failed
		return nil
	}
	d.Status = Queued
	d.NextAttemptAt = common.Now().Add(time.Duration(d.Attempt*d.Attempt) * time.Minute)
	return nil
}
