package snapshot

import (
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/common"
	"net/url"
	"strings"
)

func ValidateKind(v string) error {
	if strings.TrimSpace(v) == "" || len(v) > 64 {
		return common.ErrInvalid
	}
	return nil
}
func ValidateURL(v string) error {
	u, e := url.Parse(v)
	if e != nil || u.Scheme == "" || u.Host == "" {
		return common.ErrInvalid
	}
	return nil
}
func ValidatePage(page, size int) error {
	if page < 1 || size < 1 || size > 100 {
		return common.ErrInvalid
	}
	return nil
}
func ValidateChecksum(v string) bool {
	if len(v) != 64 {
		return false
	}
	for _, r := range v {
		if !(r >= '0' && r <= '9' || r >= 'a' && r <= 'f') {
			return false
		}
	}
	return true
}
func ValidateRecord(r Record) error {
	if r.Kind == "" || r.ID == "" || r.At.IsZero() {
		return common.ErrInvalid
	}
	return nil
}
