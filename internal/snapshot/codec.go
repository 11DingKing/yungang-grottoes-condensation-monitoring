package snapshot

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

type Envelope struct {
	Version int    `json:"version"`
	Kind    string `json:"kind"`
	Payload string `json:"payload"`
}

func Wrap(kind string, version int, value any) (Envelope, error) {
	if kind == "" || version < 1 {
		return Envelope{}, fmt.Errorf("invalid envelope")
	}
	b, e := json.Marshal(value)
	if e != nil {
		return Envelope{}, e
	}
	return Envelope{version, kind, base64.StdEncoding.EncodeToString(b)}, nil
}
func (e Envelope) Unwrap(out any) error {
	b, x := base64.StdEncoding.DecodeString(e.Payload)
	if x != nil {
		return x
	}
	return json.Unmarshal(b, out)
}
func (e Envelope) Valid() bool { return e.Version > 0 && e.Kind != "" && e.Payload != "" }
