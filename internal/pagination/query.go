package pagination

import (
	"net/url"
	"strconv"
	"strings"
)

type Query struct {
	Page, Size              int
	Sort, Direction, Filter string
}

func Parse(v url.Values) Query {
	p, _ := strconv.Atoi(v.Get("page"))
	s, _ := strconv.Atoi(v.Get("size"))
	if p < 1 {
		p = 1
	}
	if s < 1 {
		s = 20
	}
	if s > 100 {
		s = 100
	}
	dir := strings.ToLower(v.Get("direction"))
	if dir != "desc" {
		dir = "asc"
	}
	sort := v.Get("sort")
	if sort == "" {
		sort = "created_at"
	}
	return Query{p, s, sort, dir, v.Get("filter")}
}
func (q Query) Offset() int                            { return (q.Page - 1) * q.Size }
func (q Query) Limit() int                             { return q.Size }
func (q Query) ValidSort(allowed map[string]bool) bool { return allowed[q.Sort] }

type Result[T any] struct {
	Items []T `json:"items"`
	Page  int `json:"page"`
	Size  int `json:"size"`
	Total int `json:"total"`
}
