package operations

import (
	"fmt"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/common"
	"time"
)

type Approval struct {
	ID, Subject, Approver, Decision string
	RequestedAt, DecidedAt          *time.Time
	Reason                          string
}

func (a *Approval) Request(now time.Time) error {
	if a.ID == "" || a.Subject == "" {
		return common.ErrInvalid
	}
	if a.RequestedAt != nil {
		return common.ErrConflict
	}
	a.RequestedAt = &now
	a.Decision = "pending"
	return nil
}
func (a *Approval) Approve(user, reason string, now time.Time) error {
	if a.Decision != "pending" || user == "" {
		return common.ErrConflict
	}
	a.Approver = user
	a.Decision = "approved"
	a.Reason = reason
	a.DecidedAt = &now
	return nil
}
func (a *Approval) Reject(user, reason string, now time.Time) error {
	if a.Decision != "pending" || user == "" || reason == "" {
		return fmt.Errorf("%w: rejection reason", common.ErrInvalid)
	}
	a.Approver = user
	a.Decision = "rejected"
	a.Reason = reason
	a.DecidedAt = &now
	return nil
}
func (a Approval) Final() bool { return a.Decision == "approved" || a.Decision == "rejected" }
