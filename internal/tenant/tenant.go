package tenant

import (
	"context"
	"fmt"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/common"
	"strings"
)

type Key struct{ ZoneID string }

func FromContext(ctx context.Context) (Key, bool) { v, ok := ctx.Value(Key{}).(Key); return v, ok }
func With(ctx context.Context, zone string) context.Context {
	return context.WithValue(ctx, Key{}, Key{ZoneID: zone})
}
func Require(ctx context.Context, zone string) error {
	k, ok := FromContext(ctx)
	if !ok || k.ZoneID == "" || k.ZoneID != zone {
		return fmt.Errorf("%w: zone scope", common.ErrForbidden)
	}
	return nil
}
func Normalize(v string) string { return strings.ToLower(strings.TrimSpace(v)) }
