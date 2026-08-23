package snapshot

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
)

func Hash(data []byte) string { h := sha256.Sum256(data); return hex.EncodeToString(h[:]) }
func HashParts(parts ...[]byte) string {
	all := []byte{}
	for _, p := range parts {
		all = append(all, p...)
	}
	return Hash(all)
}
func StableMap(values map[string]string) []string {
	keys := make([]string, 0, len(values))
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	out := []string{}
	for _, k := range keys {
		out = append(out, k, values[k])
	}
	return out
}
func Equal(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
