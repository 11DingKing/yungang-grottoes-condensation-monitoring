package format_test

import (
	"errors"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/common"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/format"
	"testing"
	"time"
)

func TestStrings(t *testing.T) {
	if format.Trim(" x ") != "x" || format.Title("coastal response") != "Coastal Response" {
		t.Fatal()
	}
	if len(format.SHA("x")) != 64 {
		t.Fatal()
	}
	if format.Label("plan", 3) != "PLAN-0003" {
		t.Fatal()
	}
	if format.JoinNonEmpty("a", " ", "b") != "a b" {
		t.Fatal()
	}
}
func TestErrors(t *testing.T) {
	for _, c := range []struct {
		e    error
		code string
	}{{common.ErrNotFound, "not_found"}, {common.ErrConflict, "conflict"}, {common.ErrForbidden, "forbidden"}, {common.ErrInvalid, "invalid_request"}, {errors.New("x"), "internal_error"}} {
		if got := format.MapError(c.e); got.Code != c.code || !format.IsPublic(c.e, c.code) {
			t.Fatal(c, got)
		}
	}
}
func TestTime(t *testing.T) {
	a := time.Date(2026, 8, 23, 1, 0, 0, 0, time.UTC)
	b := a.Add(2 * time.Hour)
	if !format.SameDay(a, b) {
		t.Fatal()
	}
	if format.EndOfDay(a).Before(format.StartOfDay(a)) {
		t.Fatal()
	}
	if _, e := format.ParseTime(format.RFC3339(a)); e != nil {
		t.Fatal(e)
	}
}
func TestMask(t *testing.T) {
	if format.Mask("secret") != "s****t" {
		t.Fatal(format.Mask("secret"))
	}
	if format.MaskPhone("13800138000") != "138****8000" {
		t.Fatal(format.MaskPhone("13800138000"))
	}
	m := format.Redact(map[string]string{"password": "p", "name": "n"})
	if m["password"] != "[REDACTED]" || m["name"] != "n" {
		t.Fatal(m)
	}
}
