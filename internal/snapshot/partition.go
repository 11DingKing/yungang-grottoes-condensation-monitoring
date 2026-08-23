package snapshot

import (
	"sort"
	"strings"
)

type Partition struct {
	Key     string
	Records []Record
}

func PartitionBy(records []Record) []Partition {
	groups := map[string][]Record{}
	for _, r := range records {
		groups[r.Kind] = append(groups[r.Kind], r)
	}
	out := []Partition{}
	for k, v := range groups {
		out = append(out, Partition{k, v})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Key < out[j].Key })
	return out
}
func FilterRecords(records []Record, kind, prefix string) []Record {
	out := []Record{}
	for _, r := range records {
		if kind != "" && r.Kind != kind {
			continue
		}
		if prefix != "" && !strings.HasPrefix(r.ID, prefix) {
			continue
		}
		out = append(out, r)
	}
	return out
}
func RecordIDs(records []Record) []string {
	out := []string{}
	for _, r := range records {
		out = append(out, r.ID)
	}
	sort.Strings(out)
	return out
}
func Deduplicate(records []Record) []Record {
	seen := map[string]bool{}
	out := []Record{}
	for _, r := range records {
		if !seen[r.ID] {
			seen[r.ID] = true
			out = append(out, r)
		}
	}
	return out
}
