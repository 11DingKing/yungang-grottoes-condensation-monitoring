package snapshot

import (
	"fmt"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/common"
	"strings"
	"time"
)

type Manifest struct {
	ID, Service, Environment, Checksum string
	CreatedAt                          time.Time
	Files                              []string
}

func (m Manifest) Validate() error {
	if m.ID == "" || m.Service == "" || m.Environment == "" || len(m.Files) == 0 || len(m.Checksum) != 64 {
		return common.ErrInvalid
	}
	return nil
}
func (m Manifest) HasFile(name string) bool {
	for _, f := range m.Files {
		if f == name {
			return true
		}
	}
	return false
}
func (m *Manifest) AddFile(name string) error {
	if strings.TrimSpace(name) == "" || m.HasFile(name) {
		return common.ErrInvalid
	}
	m.Files = append(m.Files, name)
	return nil
}
func (m Manifest) Summary() string {
	if e := m.Validate(); e != nil {
		return fmt.Sprintf("invalid:%v", e)
	}
	return fmt.Sprintf("%s/%s files=%d", m.Service, m.Environment, len(m.Files))
}
