package identity

import (
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/common"
	"time"
)

type Session struct {
	ID, UserID           string
	ExpiresAt, CreatedAt time.Time
	RevokedAt            *time.Time
}

func (s Session) Valid(now time.Time) bool { return s.RevokedAt == nil && now.Before(s.ExpiresAt) }
func NewSession(userID string, ttl time.Duration) Session {
	now := common.Now()
	return Session{ID: common.NewID("ses"), UserID: userID, CreatedAt: now, ExpiresAt: now.Add(ttl)}
}
