package common

import "errors"

var (
	ErrNotFound    = errors.New("not found")
	ErrConflict    = errors.New("conflict")
	ErrForbidden   = errors.New("forbidden")
	ErrInvalid     = errors.New("invalid request")
	ErrUnavailable = errors.New("dependency unavailable")
)
