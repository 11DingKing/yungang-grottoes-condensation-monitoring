package format

import (
	"errors"
	"fmt"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/common"
)

type PublicError struct {
	Code, Message string
	Cause         error
}

func (e PublicError) Error() string {
	if e.Cause == nil {
		return e.Message
	}
	return fmt.Sprintf("%s: %v", e.Message, e.Cause)
}
func (e PublicError) Unwrap() error { return e.Cause }
func MapError(err error) PublicError {
	switch {
	case errors.Is(err, common.ErrNotFound):
		return PublicError{"not_found", "resource not found", err}
	case errors.Is(err, common.ErrConflict):
		return PublicError{"conflict", "state conflict", err}
	case errors.Is(err, common.ErrForbidden):
		return PublicError{"forbidden", "permission denied", err}
	case errors.Is(err, common.ErrInvalid):
		return PublicError{"invalid_request", "invalid request", err}
	default:
		return PublicError{"internal_error", "internal service error", err}
	}
}
func IsPublic(err error, code string) bool { return MapError(err).Code == code }
