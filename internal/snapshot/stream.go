package snapshot

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"time"
)

type Record struct {
	Kind, ID string
	Payload  json.RawMessage
	At       time.Time
}

func Encode(w io.Writer, records []Record) error {
	e := json.NewEncoder(w)
	for _, r := range records {
		if x := e.Encode(r); x != nil {
			return x
		}
	}
	return nil
}
func DecodeStream(r io.Reader) ([]Record, error) {
	scan := bufio.NewScanner(r)
	out := []Record{}
	for scan.Scan() {
		var v Record
		if e := json.Unmarshal(scan.Bytes(), &v); e != nil {
			return nil, e
		}
		out = append(out, v)
	}
	return out, scan.Err()
}
func RoundTrip(records []Record) ([]Record, error) {
	b := bytes.Buffer{}
	if e := Encode(&b, records); e != nil {
		return nil, e
	}
	return DecodeStream(&b)
}
