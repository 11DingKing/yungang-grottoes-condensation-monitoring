package snapshot

import (
	"bytes"
	"encoding/json"
	"sort"
	"strings"
	"time"
)

type Cursor struct {
	Offset int
	Token  string
}

func (c Cursor) Valid() bool { return c.Offset >= 0 && c.Token != "" }
func (c Cursor) Advance(n int) Cursor {
	if n < 0 {
		n = 0
	}
	c.Offset += n
	return c
}
func (c Cursor) Encode() string { return c.Token + ":" + string(rune(c.Offset)) }

type Page[T any] struct {
	Items  []T
	Cursor Cursor
	Total  int
}

func NewPage[T any](items []T, total int, c Cursor) Page[T] {
	return Page[T]{append([]T(nil), items...), c, total}
}
func (p Page[T]) Empty() bool   { return len(p.Items) == 0 }
func (p Page[T]) HasNext() bool { return p.Cursor.Offset+len(p.Items) < p.Total }
func (p Page[T]) Copy() Page[T] { return NewPage(p.Items, p.Total, p.Cursor) }

type EnvelopeMeta struct {
	TraceID, RequestID, Source string
	At                         time.Time
}

func (m EnvelopeMeta) Valid() bool {
	return m.TraceID != "" && m.RequestID != "" && m.Source != "" && !m.At.IsZero()
}
func MarshalMeta(m EnvelopeMeta) []byte { b, _ := json.Marshal(m); return b }
func UnmarshalMeta(b []byte) (EnvelopeMeta, error) {
	var m EnvelopeMeta
	e := json.Unmarshal(b, &m)
	return m, e
}
func MergeRecords(a, b []Record) []Record {
	out := append([]Record(nil), a...)
	out = append(out, b...)
	return Deduplicate(out)
}
func SortByID(records []Record) []Record {
	out := append([]Record(nil), records...)
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}
func EncodeCompact(records []Record) []byte {
	var b bytes.Buffer
	for _, r := range records {
		b.WriteString(strings.Join([]string{r.Kind, r.ID, r.At.Format(time.RFC3339), string(r.Payload)}, "|"))
		b.WriteByte('\n')
	}
	return b.Bytes()
}
func DecodeCompact(data []byte) []Record {
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	out := []Record{}
	for _, line := range lines {
		parts := strings.SplitN(line, "|", 4)
		if len(parts) != 4 {
			continue
		}
		at, _ := time.Parse(time.RFC3339, parts[2])
		out = append(out, Record{Kind: parts[0], ID: parts[1], At: at, Payload: json.RawMessage(parts[3])})
	}
	return out
}
func CountKind(records []Record, kind string) int {
	n := 0
	for _, r := range records {
		if r.Kind == kind {
			n++
		}
	}
	return n
}
func LatestByKind(records []Record, kind string) (Record, bool) {
	var out Record
	ok := false
	for _, r := range records {
		if r.Kind == kind && (!ok || r.At.After(out.At)) {
			out = r
			ok = true
		}
	}
	return out, ok
}
func ValidateOrder(records []Record) bool {
	for i := 1; i < len(records); i++ {
		if records[i].At.Before(records[i-1].At) {
			return false
		}
	}
	return true
}
