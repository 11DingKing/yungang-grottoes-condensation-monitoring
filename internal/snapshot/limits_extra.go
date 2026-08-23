package snapshot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/common"
	"io"
)

type Buffer struct {
	Max   int
	Value bytes.Buffer
}

func (b *Buffer) Write(p []byte) (int, error) {
	if b.Max < 0 {
		return 0, common.ErrInvalid
	}
	if b.Value.Len()+len(p) > b.Max {
		return 0, fmt.Errorf("%w: buffer full", common.ErrConflict)
	}
	return b.Value.Write(p)
}
func (b *Buffer) ReadAll() []byte { return append([]byte(nil), b.Value.Bytes()...) }
func (b *Buffer) Reset()          { b.Value.Reset() }
func CopyLimited(dst io.Writer, src io.Reader, max int) (int64, error) {
	if max < 0 {
		return 0, common.ErrInvalid
	}
	return io.CopyN(dst, src, int64(max))
}
func EncodeValue(value any) ([]byte, error) { return json.Marshal(value) }
func DecodeValue(data []byte, out any) error {
	if len(data) == 0 {
		return common.ErrInvalid
	}
	return json.Unmarshal(data, out)
}
func Size(data []byte) int { return len(data) }
