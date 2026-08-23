package rules

import (
	"path"
	"strings"
)

type Route struct {
	Method, Pattern string
	Roles           []string
}

func (r Route) Match(method, urlPath string) bool {
	if !strings.EqualFold(r.Method, method) {
		return false
	}
	ok, e := path.Match(r.Pattern, urlPath)
	return e == nil && ok
}
func (r Route) Public() bool { return len(r.Roles) == 0 }
func (r Route) Allows(role string) bool {
	for _, v := range r.Roles {
		if v == role {
			return true
		}
	}
	return false
}
func SortRoutes(routes []Route) []Route {
	out := append([]Route(nil), routes...)
	for i := 0; i < len(out); i++ {
		for j := i + 1; j < len(out); j++ {
			if len(out[j].Pattern) > len(out[i].Pattern) {
				out[i], out[j] = out[j], out[i]
			}
		}
	}
	return out
}
