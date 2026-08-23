package request

import (
	"fmt"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/common"
)

type Status string

const (
	Submitted  Status = "submitted"
	Approved   Status = "approved"
	Allocated  Status = "allocated"
	Dispatched Status = "dispatched"
	Received   Status = "received"
	Cancelled  Status = "cancelled"
)

type Item struct {
	ID, RequestID, SupplyID string
	Quantity, Allocated     int
}
type Request struct {
	ID, ZoneID, RequesterID, IdempotencyKey, Destination string
	Status                                               Status
	Priority, Version                                    int
	Items                                                []Item
}

func (r *Request) Transition(next Status) error {
	valid := map[Status]map[Status]bool{Submitted: {Approved: true, Cancelled: true}, Approved: {Allocated: true, Cancelled: true}, Allocated: {Dispatched: true}, Dispatched: {Received: true}, Received: {}, Cancelled: {}}
	if !valid[r.Status][next] {
		return fmt.Errorf("%w: request transition %s to %s", common.ErrConflict, r.Status, next)
	}
	r.Status = next
	r.Version++
	return nil
}
func (r Request) Complete() bool {
	if r.Status != Received {
		return false
	}
	for _, i := range r.Items {
		if i.Allocated < i.Quantity {
			return false
		}
	}
	return true
}
