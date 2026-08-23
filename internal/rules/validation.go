package rules

import (
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/common"
	"regexp"
	"strings"
)

var codePattern = regexp.MustCompile(`^[A-Z]{2}-[0-9]{2,6}$`)

func ValidateZoneCode(v string) error {
	if !codePattern.MatchString(strings.TrimSpace(v)) {
		return common.ErrInvalid
	}
	return nil
}
func ValidatePositive(values ...int) error {
	for _, v := range values {
		if v <= 0 {
			return common.ErrInvalid
		}
	}
	return nil
}
func ValidateNonNegative(values ...int) error {
	for _, v := range values {
		if v < 0 {
			return common.ErrInvalid
		}
	}
	return nil
}
func ValidateUnique(values []string) bool {
	seen := map[string]bool{}
	for _, v := range values {
		if v == "" || seen[v] {
			return false
		}
		seen[v] = true
	}
	return true
}
func NormalizeTags(values []string) []string {
	out := make([]string, 0, len(values))
	seen := map[string]bool{}
	for _, v := range values {
		v = strings.ToLower(strings.TrimSpace(v))
		if v != "" && !seen[v] {
			seen[v] = true
			out = append(out, v)
		}
	}
	return out
}
