package snapshot

import (
	"fmt"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/common"
)

type Limits struct{ MaxBytes, MaxRecords int }

func (l Limits) Validate() error {
	if l.MaxBytes < 1 || l.MaxRecords < 1 {
		return common.ErrInvalid
	}
	return nil
}
func (l Limits) CheckBytes(n int) error {
	if n < 0 || n > l.MaxBytes {
		return fmt.Errorf("%w: bytes", common.ErrConflict)
	}
	return nil
}
func (l Limits) CheckRecords(n int) error {
	if n < 0 || n > l.MaxRecords {
		return fmt.Errorf("%w: records", common.ErrConflict)
	}
	return nil
}
func (l Limits) RemainingBytes(n int) int {
	if n >= l.MaxBytes {
		return 0
	}
	return l.MaxBytes - n
}
func (l Limits) RemainingRecords(n int) int {
	if n >= l.MaxRecords {
		return 0
	}
	return l.MaxRecords - n
}
