package aggregate

import (
	"fmt"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/common"
	"sort"
	"time"
)

type Request struct {
	ID, Zone, Destination string
	Items                 map[string]int
	Status                string
	Version               int
	UpdatedAt             time.Time
}

func NewRequest(id, zone, destination string) Request {
	return Request{ID: id, Zone: zone, Destination: destination, Items: map[string]int{}, Status: "submitted", Version: 1, UpdatedAt: common.Now()}
}
func (r *Request) AddItem(sku string, n int) error {
	if sku == "" || n <= 0 {
		return common.ErrInvalid
	}
	if r.Status != "submitted" {
		return common.ErrConflict
	}
	r.Items[sku] += n
	r.Version++
	r.UpdatedAt = common.Now()
	return nil
}
func (r *Request) RemoveItem(sku string, n int) error {
	if n <= 0 || r.Items[sku] < n {
		return common.ErrInvalid
	}
	r.Items[sku] -= n
	if r.Items[sku] == 0 {
		delete(r.Items, sku)
	}
	r.Version++
	r.UpdatedAt = common.Now()
	return nil
}
func (r *Request) Approve() error {
	if len(r.Items) == 0 || r.Status != "submitted" {
		return fmt.Errorf("%w: request not ready", common.ErrConflict)
	}
	r.Status = "approved"
	r.Version++
	return nil
}
func (r *Request) Cancel() error {
	if r.Status == "received" || r.Status == "cancelled" {
		return common.ErrConflict
	}
	r.Status = "cancelled"
	r.Version++
	return nil
}
func (r Request) SKUs() []string {
	out := make([]string, 0, len(r.Items))
	for k := range r.Items {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
func (r Request) Total() int {
	n := 0
	for _, v := range r.Items {
		n += v
	}
	return n
}
