package snapshot

import (
	"strings"
	"time"
)

type Inspection struct {
	Records, Bytes, Invalid int
	Started, Finished       time.Time
}

func (i Inspection) Duration() time.Duration { return i.Finished.Sub(i.Started) }
func (i Inspection) OK() bool                { return i.Records > 0 && i.Invalid == 0 }
func Inspect(records []Record) Inspection {
	start := time.Now()
	out := Inspection{Started: start}
	for _, r := range records {
		out.Records++
		out.Bytes += len(r.Payload)
		if strings.TrimSpace(r.ID) == "" || r.At.IsZero() {
			out.Invalid++
		}
	}
	out.Finished = time.Now()
	return out
}
func Summarize(records []Record) map[string]int {
	out := map[string]int{}
	for _, r := range records {
		out[r.Kind]++
	}
	return out
}
func Kinds(records []Record) []string {
	set := map[string]bool{}
	for _, r := range records {
		set[r.Kind] = true
	}
	out := []string{}
	for k := range set {
		out = append(out, k)
	}
	return out
}
func Has(records []Record, id string) bool {
	for _, r := range records {
		if r.ID == id {
			return true
		}
	}
	return false
}
func Empty(records []Record) bool { return len(records) == 0 }
func PayloadBytes(records []Record) int {
	n := 0
	for _, r := range records {
		n += len(r.Payload)
	}
	return n
}
func ValidKinds(records []Record) bool {
	for _, r := range records {
		if strings.TrimSpace(r.Kind) == "" {
			return false
		}
	}
	return true
}
func CountValid(records []Record) int {
	n := 0
	for _, r := range records {
		if r.ID != "" && r.Kind != "" {
			n++
		}
	}
	return n
}
