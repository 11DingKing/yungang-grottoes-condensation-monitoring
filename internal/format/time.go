package format

import (
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/common"
	"time"
)

func RFC3339(t time.Time) string            { return t.In(common.Fujian).Format(time.RFC3339) }
func ParseTime(v string) (time.Time, error) { return time.Parse(time.RFC3339, v) }
func SameDay(a, b time.Time) bool {
	aa := a.In(common.Fujian)
	bb := b.In(common.Fujian)
	return aa.Year() == bb.Year() && aa.YearDay() == bb.YearDay()
}
func StartOfDay(t time.Time) time.Time {
	v := t.In(common.Fujian)
	return time.Date(v.Year(), v.Month(), v.Day(), 0, 0, 0, 0, common.Fujian)
}
func EndOfDay(t time.Time) time.Time { return StartOfDay(t).Add(24 * time.Hour).Add(-time.Nanosecond) }
