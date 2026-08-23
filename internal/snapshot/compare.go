package snapshot

import (
	"reflect"
	"sort"
)

type Difference struct {
	Field         string
	Before, After any
}

func Diff(before, after map[string]any) []Difference {
	keys := map[string]bool{}
	for k := range before {
		keys[k] = true
	}
	for k := range after {
		keys[k] = true
	}
	out := []Difference{}
	for k := range keys {
		if !reflect.DeepEqual(before[k], after[k]) {
			out = append(out, Difference{k, before[k], after[k]})
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Field < out[j].Field })
	return out
}
func EqualMaps(a, b map[string]any) bool { return len(Diff(a, b)) == 0 }
func MergeMaps(base, overlay map[string]any) map[string]any {
	out := map[string]any{}
	for k, v := range base {
		out[k] = v
	}
	for k, v := range overlay {
		out[k] = v
	}
	return out
}
