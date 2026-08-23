package format

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

func Trim(v string) string { return strings.TrimSpace(v) }
func Title(v string) string {
	parts := strings.Fields(strings.ToLower(v))
	for i, p := range parts {
		r := []rune(p)
		if len(r) > 0 {
			r[0] = unicode.ToUpper(r[0])
		}
		parts[i] = string(r)
	}
	return strings.Join(parts, " ")
}
func SHA(v string) string { h := sha256.Sum256([]byte(v)); return hex.EncodeToString(h[:]) }
func Int(v int) string    { return strconv.Itoa(v) }
func JoinNonEmpty(values ...string) string {
	out := []string{}
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			out = append(out, strings.TrimSpace(v))
		}
	}
	return strings.Join(out, " ")
}
func Label(kind string, n int) string { return fmt.Sprintf("%s-%04d", strings.ToUpper(kind), n) }
