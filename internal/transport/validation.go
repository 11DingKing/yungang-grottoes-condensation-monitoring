package transport

import (
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/common"
	"net/url"
	"strings"
	"unicode/utf8"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (l LoginRequest) Validate() error {
	if utf8.RuneCountInString(strings.TrimSpace(l.Username)) < 3 {
		return common.ErrInvalid
	}
	if len(l.Password) < 3 {
		return common.ErrInvalid
	}
	return nil
}

type PlanRequest struct {
	Name   string `json:"name"`
	ZoneID string `json:"zone_id"`
	Level  int    `json:"level"`
}

func (p PlanRequest) Validate() error {
	if strings.TrimSpace(p.Name) == "" || p.ZoneID == "" || p.Level < 1 || p.Level > 4 {
		return common.ErrInvalid
	}
	return nil
}

type RequestItem struct {
	SupplyID string `json:"supply_id"`
	Quantity int    `json:"quantity"`
}
type ReliefRequest struct {
	ZoneID      string        `json:"zone_id"`
	Destination string        `json:"destination"`
	Priority    int           `json:"priority"`
	Items       []RequestItem `json:"items"`
}

func (r ReliefRequest) Validate() error {
	if r.ZoneID == "" || strings.TrimSpace(r.Destination) == "" || len(r.Items) == 0 || r.Priority < 1 || r.Priority > 5 {
		return common.ErrInvalid
	}
	seen := map[string]bool{}
	for _, i := range r.Items {
		if i.SupplyID == "" || i.Quantity <= 0 || seen[i.SupplyID] {
			return common.ErrInvalid
		}
		seen[i.SupplyID] = true
	}
	return nil
}
func SafeNext(v url.Values) string {
	next := v.Get("next")
	if strings.HasPrefix(next, "/") {
		return next
	}
	return "/"
}
