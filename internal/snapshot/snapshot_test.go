package snapshot_test

import (
	"bytes"
	"context"
	"errors"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/common"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/snapshot"
	"testing"
	"time"
)

func TestStoreRoundTrip(t *testing.T) {
	s := snapshot.New()
	if e := s.Put("request", "r", map[string]int{"n": 2}); e != nil {
		t.Fatal(e)
	}
	h, b, e := s.Get("r")
	if e != nil || h.Version != 1 {
		t.Fatal(h, e)
	}
	var got map[string]int
	got, e = snapshot.Decode[map[string]int](b)
	if e != nil || got["n"] != 2 {
		t.Fatal(got, e)
	}
	if e = s.Delete("r"); e != nil || s.Count() != 0 {
		t.Fatal(e, s.Count())
	}
}
func TestVersioned(t *testing.T) {
	v := snapshot.NewVersions[string]()
	if e := v.Put("r", "one", 0); e != nil {
		t.Fatal(e)
	}
	if e := v.Put("r", "bad", 0); !errors.Is(e, common.ErrConflict) {
		t.Fatal(e)
	}
	if e := v.Put("r", "two", 1); e != nil {
		t.Fatal(e)
	}
	x, e := v.Get("r")
	if e != nil || x.Value != "two" {
		t.Fatal(x, e)
	}
}
func TestStream(t *testing.T) {
	now := time.Now()
	in := []snapshot.Record{{Kind: "a", ID: "1", At: now}, {Kind: "b", ID: "2", At: now}}
	out, e := snapshot.RoundTrip(in)
	if e != nil || len(out) != 2 {
		t.Fatal(out, e)
	}
	var b bytes.Buffer
	if e = snapshot.Encode(&b, in); e != nil {
		t.Fatal(e)
	}
}
func TestChecksum(t *testing.T) {
	if snapshot.Hash([]byte("x")) != snapshot.HashParts([]byte("x")) {
		t.Fatal()
	}
	if !snapshot.Equal([]byte("a"), []byte("a")) {
		t.Fatal()
	}
	if snapshot.Equal([]byte("a"), []byte("b")) {
		t.Fatal()
	}
	if len(snapshot.StableMap(map[string]string{"b": "2", "a": "1"})) != 4 {
		t.Fatal()
	}
}
func TestRecovery(t *testing.T) {
	m := &snapshot.Manager{}
	r := &fakeRecover{}
	m.Add(r)
	if e := m.Run(context.Background()); len(e) != 0 || !r.called {
		t.Fatal(e, r.called)
	}
	marker := snapshot.Marker{}
	marker.Mark()
	if !marker.Valid() {
		t.Fatal(marker)
	}
}
func TestManifest(t *testing.T) {
	m := snapshot.Manifest{ID: "m", Service: "api", Environment: "test", Checksum: "1234567890123456789012345678901234567890123456789012345678901234", CreatedAt: time.Now(), Files: []string{"a"}}
	if e := m.AddFile("b"); e != nil || !m.HasFile("b") {
		t.Fatal(e)
	}
	if e := m.Validate(); e != nil {
		t.Fatal(e)
	}
	if m.Summary() == "" {
		t.Fatal()
	}
}
func TestRetention(t *testing.T) {
	now := time.Now()
	items := []snapshot.Artifact{{ID: "old", CreatedAt: now.Add(-time.Hour), Size: 2}, {ID: "new", CreatedAt: now, Size: 3}, {ID: "pinned", CreatedAt: now.Add(-2 * time.Hour), Pinned: true}}
	if len(snapshot.Keep(items, 2)) != 2 || len(snapshot.Expired(items, now.Add(-time.Minute))) != 1 || snapshot.TotalSize(items) != 5 {
		t.Fatal()
	}
}

type fakeRecover struct{ called bool }

func (f *fakeRecover) Recover(context.Context) error { f.called = true; return nil }
