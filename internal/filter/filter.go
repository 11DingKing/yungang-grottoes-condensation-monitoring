package filter

import (
	"sort"
	"strings"
)

type Predicate[T any] func(T) bool

func Contains(value, needle string) bool {
	return needle == "" || strings.Contains(strings.ToLower(value), strings.ToLower(needle))
}
func Apply[T any](items []T, p Predicate[T]) []T {
	out := make([]T, 0, len(items))
	for _, v := range items {
		if p == nil || p(v) {
			out = append(out, v)
		}
	}
	return out
}
func SortBy[T any](items []T, less func(T, T) bool) []T {
	out := append([]T(nil), items...)
	sort.SliceStable(out, func(i, j int) bool { return less(out[i], out[j]) })
	return out
}
func Page[T any](items []T, page, size int) []T {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 20
	}
	start := (page - 1) * size
	if start >= len(items) {
		return []T{}
	}
	end := start + size
	if end > len(items) {
		end = len(items)
	}
	return append([]T(nil), items[start:end]...)
}
