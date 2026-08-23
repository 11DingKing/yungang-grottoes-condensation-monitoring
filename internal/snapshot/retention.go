package snapshot

import (
	"sort"
	"time"
)

type Artifact struct {
	ID        string
	CreatedAt time.Time
	Size      int
	Pinned    bool
}

func Keep(items []Artifact, n int) []Artifact {
	out := append([]Artifact(nil), items...)
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt.After(out[j].CreatedAt) })
	if n < 0 {
		n = 0
	}
	if len(out) > n {
		out = out[:n]
	}
	return out
}
func Expired(items []Artifact, before time.Time) []Artifact {
	out := []Artifact{}
	for _, a := range items {
		if !a.Pinned && a.CreatedAt.Before(before) {
			out = append(out, a)
		}
	}
	return out
}
func TotalSize(items []Artifact) int {
	n := 0
	for _, a := range items {
		n += a.Size
	}
	return n
}
