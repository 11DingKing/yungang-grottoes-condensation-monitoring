package aggregate

import (
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/common"
	"strings"
)

func Required(values ...string) error {
	for _, v := range values {
		if strings.TrimSpace(v) == "" {
			return common.ErrInvalid
		}
	}
	return nil
}
func ValidID(v string) bool { return len(strings.TrimSpace(v)) >= 3 }
func ValidStatus(v string, allowed ...string) bool {
	for _, a := range allowed {
		if v == a {
			return true
		}
	}
	return false
}
func CloneItems(in map[string]int) map[string]int {
	out := map[string]int{}
	for k, v := range in {
		out[k] = v
	}
	return out
}
func EqualItems(a, b map[string]int) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}
func Difference(a, b map[string]int) map[string]int {
	out := CloneItems(a)
	for k, v := range b {
		out[k] -= v
		if out[k] == 0 {
			delete(out, k)
		}
	}
	return out
}
func Sum(values map[string]int) int {
	n := 0
	for _, v := range values {
		n += v
	}
	return n
}
