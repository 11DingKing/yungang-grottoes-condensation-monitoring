package format

import (
	"strings"
	"unicode/utf8"
)

func Mask(v string) string {
	n := utf8.RuneCountInString(v)
	if n <= 2 {
		return strings.Repeat("*", n)
	}
	r := []rune(v)
	return string(r[:1]) + strings.Repeat("*", n-2) + string(r[n-1:])
}
func MaskPhone(v string) string {
	r := []rune(v)
	if len(r) < 7 {
		return Mask(v)
	}
	return string(r[:3]) + "****" + string(r[len(r)-4:])
}
func Redact(values map[string]string) map[string]string {
	out := map[string]string{}
	for k, v := range values {
		switch strings.ToLower(k) {
		case "password", "token", "secret":
			out[k] = "[REDACTED]"
		default:
			out[k] = v
		}
	}
	return out
}
