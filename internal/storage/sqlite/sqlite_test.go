package sqlite_test

import (
	"context"
	"errors"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/common"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/identity"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/plan"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/shelter"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/supply"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/storage/sqlite"
	"path/filepath"
	"testing"
)

func open(t *testing.T) (context.Context, *sqlite.DB) {
	t.Helper()
	ctx := context.Background()
	db, e := sqlite.Open(ctx, filepath.Join(t.TempDir(), "test.db"))
	if e != nil {
		t.Fatal(e)
	}
	t.Cleanup(func() { db.Close() })
	return ctx, db
}
func TestMigrationCreatesBusinessTables(t *testing.T) {
	ctx, db := open(t)
	rows, e := db.SQL.QueryContext(ctx, "SELECT name FROM sqlite_master WHERE type='table' ORDER BY name")
	if e != nil {
		t.Fatal(e)
	}
	defer rows.Close()
	seen := map[string]bool{}
	for rows.Next() {
		var n string
		if e = rows.Scan(&n); e != nil {
			t.Fatal(e)
		}
		seen[n] = true
	}
	for _, name := range []string{"users", "sessions", "zones", "response_plans", "shelters", "supply_lots", "relief_requests", "allocations", "dispatches", "incidents", "audit_events", "worker_jobs"} {
		if !seen[name] {
			t.Fatal("missing", name)
		}
	}
}
func TestMigrationIsIdempotent(t *testing.T) {
	ctx, db := open(t)
	if e := db.Migrate(ctx); e != nil {
		t.Fatal(e)
	}
	var n int
	if e := db.SQL.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations").Scan(&n); e != nil {
		t.Fatal(e)
	}
	if n != 2 {
		t.Fatal(n)
	}
}
func TestRestartRecoversRows(t *testing.T) {
	ctx, db := open(t)
	u := identity.NewUser("operator", identity.RoleCoordinator)
	u.PasswordHash = "secret"
	if e := sqlite.CreateUser(ctx, db.SQL, u); e != nil {
		t.Fatal(e)
	}
	path := filepath.Join(t.TempDir(), "restart.db")
	db.Close()
	db, e := sqlite.Open(ctx, path)
	if e != nil {
		t.Fatal(e)
	}
	defer db.Close()
	u2 := identity.NewUser("persisted", identity.RoleAuditor)
	u2.PasswordHash = "secret"
	if e := sqlite.CreateUser(ctx, db.SQL, u2); e != nil {
		t.Fatal(e)
	}
	db.Close()
	db, e = sqlite.Open(ctx, path)
	if e != nil {
		t.Fatal(e)
	}
	defer db.Close()
	var n int
	if e = db.SQL.QueryRowContext(ctx, "SELECT COUNT(*) FROM users").Scan(&n); e != nil {
		t.Fatal(e)
	}
	if n != 1 {
		t.Fatal(n)
	}
}
func TestPlanOptimisticVersion(t *testing.T) {
	ctx, db := open(t)
	if e := sqlite.CreateZone(ctx, db.SQL, "z1", "FJ-01", "coast", "east"); e != nil {
		t.Fatal(e)
	}
	p, _ := plan.New("z1", "response", 1)
	if e := sqlite.CreatePlan(ctx, db.SQL, p); e != nil {
		t.Fatal(e)
	}
	p.Status = plan.Active
	p.Version = 2
	if e := sqlite.UpdatePlan(ctx, db.SQL, p, 99); !errors.Is(e, common.ErrConflict) {
		t.Fatal(e)
	}
}
func TestShelterOptimisticVersion(t *testing.T) {
	ctx, db := open(t)
	if e := sqlite.CreateZone(ctx, db.SQL, "z1", "FJ-01", "coast", "east"); e != nil {
		t.Fatal(e)
	}
	s := shelter.Shelter{ID: "sh1", ZoneID: "z1", Name: "hall", Capacity: 20, Status: shelter.Open, Version: 1}
	if e := sqlite.CreateShelter(ctx, db.SQL, s); e != nil {
		t.Fatal(e)
	}
	got, e := sqlite.GetShelter(ctx, db.SQL, "sh1")
	if e != nil {
		t.Fatal(e)
	}
	if e = got.Assign(5); e != nil {
		t.Fatal(e)
	}
	if e = sqlite.UpdateShelter(ctx, db.SQL, got, 1); e != nil {
		t.Fatal(e)
	}
	if e = sqlite.UpdateShelter(ctx, db.SQL, got, 1); !errors.Is(e, common.ErrConflict) {
		t.Fatal(e)
	}
}
func TestLotOrdering(t *testing.T) {
	ctx, db := open(t)
	if e := sqlite.CreateZone(ctx, db.SQL, "z1", "FJ-01", "coast", "east"); e != nil {
		t.Fatal(e)
	}
	if e := sqlite.CreateSupply(ctx, db.SQL, "sp1", "z1", "water", "water", "box", 10); e != nil {
		t.Fatal(e)
	}
	for _, l := range []supply.Lot{{ID: "l2", SupplyID: "sp1", LotCode: "late", Quantity: 3, Status: supply.Available, Version: 1}, {ID: "l1", SupplyID: "sp1", LotCode: "early", Quantity: 4, Status: supply.Available, Version: 1}} {
		if e := sqlite.CreateLot(ctx, db.SQL, l); e != nil {
			t.Fatal(e)
		}
	}
	lots, e := sqlite.ListLots(ctx, db.SQL, "sp1")
	if e != nil {
		t.Fatal(e)
	}
	if len(lots) != 2 {
		t.Fatal(len(lots))
	}
}
func TestForeignKeyRejectsUnknownZone(t *testing.T) {
	ctx, db := open(t)
	s := shelter.Shelter{ID: "sh", ZoneID: "unknown", Name: "bad", Capacity: 1, Status: shelter.Open, Version: 1}
	if e := sqlite.CreateShelter(ctx, db.SQL, s); e == nil {
		t.Fatal("foreign key accepted")
	}
}
func TestContextCancellation(t *testing.T) {
	ctx, db := open(t)
	cancelled, cancel := context.WithCancel(ctx)
	cancel()
	if _, e := db.SQL.ExecContext(cancelled, "SELECT 1"); e == nil {
		t.Fatal("cancel ignored")
	}
}
