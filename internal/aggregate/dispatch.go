package aggregate

import (
	"fmt"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/common"
	"time"
)

type Dispatch struct {
	ID, RequestID, Vehicle, Driver, Status string
	Attempts                               int
	CreatedAt, UpdatedAt                   time.Time
	Events                                 []string
}

func (d *Dispatch) Record(event string) {
	d.Events = append(d.Events, event)
	d.UpdatedAt = common.Now()
}
func (d *Dispatch) Assign(driver, vehicle string) error {
	if driver == "" || vehicle == "" || d.Status != "queued" {
		return common.ErrInvalid
	}
	d.Driver = driver
	d.Vehicle = vehicle
	d.Status = "assigned"
	d.Record("assigned")
	return nil
}
func (d *Dispatch) Start() error {
	if d.Status != "assigned" {
		return common.ErrConflict
	}
	d.Status = "in_transit"
	d.Record("started")
	return nil
}
func (d *Dispatch) Deliver() error {
	if d.Status != "in_transit" {
		return common.ErrConflict
	}
	d.Status = "delivered"
	d.Record("delivered")
	return nil
}
func (d *Dispatch) Fail(reason string) error {
	if d.Status == "delivered" || reason == "" {
		return common.ErrConflict
	}
	d.Status = "failed"
	d.Attempts++
	d.Record("failed:" + reason)
	return nil
}
func (d Dispatch) Validate() error {
	if d.ID == "" || d.RequestID == "" {
		return common.ErrInvalid
	}
	switch d.Status {
	case "queued", "assigned", "in_transit", "delivered", "failed":
		return nil
	default:
		return fmt.Errorf("%w: status", common.ErrInvalid)
	}
}
