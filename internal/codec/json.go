package codec

import (
	"encoding/json"
	"io"
	"net/http"
)

func Decode[T any](r io.Reader) (T, error) {
	var out T
	e := json.NewDecoder(io.LimitReader(r, 1<<20)).Decode(&out)
	return out, e
}
func Encode(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}
func Must(v any) []byte { b, _ := json.Marshal(v); return b }
