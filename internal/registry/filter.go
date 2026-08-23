package registry

import (
	"sort"
	"strings"
)

func Filter(entries []Entry, zone, term string) []Entry {
	out := []Entry{}
	term = strings.ToLower(term)
	for _, e := range entries {
		if zone != "" && e.Zone != zone {
			continue
		}
		if term != "" && !strings.Contains(strings.ToLower(e.Name), term) {
			continue
		}
		out = append(out, clone(e))
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}
func Active(entries []Entry) []Entry {
	out := []Entry{}
	for _, e := range entries {
		if e.Active {
			out = append(out, clone(e))
		}
	}
	return out
}
func Zones(entries []Entry) []string {
	set := map[string]bool{}
	for _, e := range entries {
		set[e.Zone] = true
	}
	out := []string{}
	for z := range set {
		out = append(out, z)
	}
	sort.Strings(out)
	return out
}
