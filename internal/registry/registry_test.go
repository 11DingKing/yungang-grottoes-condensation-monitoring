package registry_test

import (
	"context"
	"errors"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/common"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/registry"
	"testing"
	"time"
)

func TestRegistryLifecycle(t *testing.T) {
	r := registry.New()
	e := registry.Entry{ID: "e1", Name: "shelter", Zone: "z", Metadata: map[string]string{"kind": "hall"}}
	if x := r.Put(e); x != nil {
		t.Fatal(x)
	}
	got, x := r.Get("e1")
	if x != nil || got.Name != "shelter" {
		t.Fatal(got, x)
	}
	got.Metadata["kind"] = "mutated"
	again, _ := r.Get("e1")
	if again.Metadata["kind"] != "hall" {
		t.Fatal("metadata leaked")
	}
	if x = r.Deactivate("e1"); x != nil {
		t.Fatal(x)
	}
	if len(r.List("z")) != 0 {
		t.Fatal(r.List("z"))
	}
}
func TestRegistryDuplicate(t *testing.T) {
	r := registry.New()
	e := registry.Entry{ID: "e", Name: "n", Zone: "z"}
	if x := r.Put(e); x != nil {
		t.Fatal(x)
	}
	if x := r.Put(e); !errors.Is(x, common.ErrConflict) {
		t.Fatal(x)
	}
}
func TestLeases(t *testing.T) {
	l := registry.NewLeases()
	now := time.Now()
	if e := l.Acquire("k", "a", time.Minute, now); e != nil {
		t.Fatal(e)
	}
	if e := l.Acquire("k", "b", time.Minute, now); !errors.Is(e, common.ErrConflict) {
		t.Fatal(e)
	}
	if !l.Active("k", now) {
		t.Fatal()
	}
	if e := l.Renew("k", "a", time.Minute, now); e != nil {
		t.Fatal(e)
	}
	if e := l.Release("k", "b"); !errors.Is(e, common.ErrForbidden) {
		t.Fatal(e)
	}
	if e := l.Release("k", "a"); e != nil {
		t.Fatal(e)
	}
}
func TestCatalog(t *testing.T) {
	c := registry.NewCatalog()
	c.Add(registry.CatalogItem{Code: "water", Name: "Drinking Water", Enabled: true})
	c.Add(registry.CatalogItem{Code: "food", Name: "Food", Enabled: true})
	if len(c.Search("water")) != 1 {
		t.Fatal(c.Search("water"))
	}
	if !c.Disable("water") {
		t.Fatal()
	}
	if len(c.Search("water")) != 0 {
		t.Fatal()
	}
}
func TestHistory(t *testing.T) {
	h := &registry.History{}
	h.Append(registry.Change{Entity: "request", ID: "r", Action: "created"})
	h.Append(registry.Change{Entity: "request", ID: "r", Action: "approved"})
	if got := h.For("request", "r"); len(got) != 2 {
		t.Fatal(got)
	}
	if last, ok := h.Latest("request", "r"); !ok || last.Action != "approved" {
		t.Fatal(last, ok)
	}
}
func TestLocks(t *testing.T) {
	l := registry.NewLocks()
	unlock, e := l.Lock(context.Background(), "k")
	if e != nil {
		t.Fatal(e)
	}
	unlock()
	if _, ok := l.Try("k", time.Millisecond); !ok {
		t.Fatal("lock unavailable")
	}
}
