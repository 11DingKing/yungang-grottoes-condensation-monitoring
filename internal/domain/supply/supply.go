package supply

import (
	"fmt"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/common"
	"time"
)

type LotStatus string

const (
	Available LotStatus = "available"
	Exhausted LotStatus = "exhausted"
	Expired   LotStatus = "expired"
)

type Lot struct {
	ID, SupplyID, LotCode       string
	Quantity, Reserved, Version int
	ExpiresAt                   *time.Time
	Status                      LotStatus
}

func (l Lot) Free(now time.Time) int {
	if l.ExpiresAt != nil && !now.Before(*l.ExpiresAt) {
		return 0
	}
	return l.Quantity - l.Reserved
}
func (l *Lot) Reserve(n int, now time.Time) error {
	if n <= 0 {
		return common.ErrInvalid
	}
	if l.ExpiresAt != nil && !now.Before(*l.ExpiresAt) {
		l.Status = Expired
		return fmt.Errorf("%w: lot expired", common.ErrConflict)
	}
	if l.Free(now) < n {
		return fmt.Errorf("%w: insufficient lot", common.ErrConflict)
	}
	l.Reserved += n
	return nil
}
func (l *Lot) Consume(n int) error {
	if n <= 0 || n > l.Reserved {
		return common.ErrInvalid
	}
	l.Quantity -= n
	l.Reserved -= n
	if l.Quantity == 0 {
		l.Status = Exhausted
	}
	l.Version++
	return nil
}
