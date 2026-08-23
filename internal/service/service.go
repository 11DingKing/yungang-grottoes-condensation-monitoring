package service

import (
	"context"
	"errors"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/audit"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/common"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/dispatch"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/identity"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/incident"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/plan"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/request"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/shelter"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/supply"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/idempotency"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/monitoring"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/storage/sqlite"
	"time"
)

type Service struct {
	DB         *sqlite.DB
	SessionTTL time.Duration
	Monitor    *monitoring.Ledger
}

func New(db *sqlite.DB) *Service {
	return &Service{DB: db, SessionTTL: 12 * time.Hour, Monitor: monitoring.NewLedger()}
}

func (s *Service) AssessEnvironment(ctx context.Context, actor identity.User, sample monitoring.Sample) (monitoring.Decision, error) {
	if !actor.Can("plan:write") {
		return monitoring.Decision{}, common.ErrForbidden
	}
	return s.Monitor.Record(ctx, sample, common.Now(), 15*time.Minute)
}
func (s *Service) tx(ctx context.Context, fn func(sqlite.Tx) error) error {
	tx, e := s.DB.SQL.BeginTx(ctx, nil)
	if e != nil {
		return e
	}
	if e = fn(tx); e != nil {
		_ = tx.Rollback()
		return e
	}
	return tx.Commit()
}
func (s *Service) CreateUser(ctx context.Context, u identity.User) error {
	return s.tx(ctx, func(q sqlite.Tx) error { return sqlite.CreateUser(ctx, q, u) })
}
func (s *Service) Login(ctx context.Context, username, password string) (identity.Session, error) {
	var out identity.Session
	e := s.tx(ctx, func(q sqlite.Tx) error {
		u, e := sqlite.FindUser(ctx, q, username)
		if e != nil {
			return e
		}
		if !u.Active || u.PasswordHash != password {
			return common.ErrForbidden
		}
		out = identity.NewSession(u.ID, s.SessionTTL)
		return sqlite.SaveSession(ctx, q, out)
	})
	return out, e
}
func (s *Service) Logout(ctx context.Context, id string) error {
	return s.tx(ctx, func(q sqlite.Tx) error { return sqlite.RevokeSession(ctx, q, id) })
}
func (s *Service) Authenticate(ctx context.Context, id string) (identity.User, error) {
	var u identity.User
	e := s.tx(ctx, func(q sqlite.Tx) error {
		ses, e := sqlite.FindSession(ctx, q, id)
		if e != nil {
			return e
		}
		if !ses.Valid(common.Now()) {
			return common.ErrForbidden
		}
		u, e = sqlite.FindUser(ctx, q, "coordinator")
		if e != nil && !errors.Is(e, common.ErrNotFound) {
			return e
		}
		if e == nil && u.ID == ses.UserID {
			return nil
		}
		rows, e := q.QueryContext(ctx, "SELECT username FROM users WHERE id=?", ses.UserID)
		if e != nil {
			return e
		}
		defer rows.Close()
		if !rows.Next() {
			return common.ErrNotFound
		}
		var name string
		_ = rows.Scan(&name)
		u, e = sqlite.FindUser(ctx, q, name)
		return e
	})
	return u, e
}
func (s *Service) CreatePlan(ctx context.Context, actor identity.User, p plan.Plan) error {
	if !actor.Can("plan:write") {
		return common.ErrForbidden
	}
	return s.tx(ctx, func(q sqlite.Tx) error {
		if e := sqlite.CreatePlan(ctx, q, p); e != nil {
			return e
		}
		return audit.Record(ctx, q, audit.Event{ActorID: actor.ID, ObjectType: "plan", ObjectID: p.ID, Action: "create", Result: "ok", RequestID: common.NewID("req"), Detail: p.Name})
	})
}
func (s *Service) TransitionPlan(ctx context.Context, actor identity.User, id string, next plan.Status) error {
	if !actor.Can("plan:write") {
		return common.ErrForbidden
	}
	return s.tx(ctx, func(q sqlite.Tx) error {
		p, e := sqlite.GetPlan(ctx, q, id)
		if e != nil {
			return e
		}
		old := p.Version
		if e = p.Transition(next); e != nil {
			return e
		}
		if e = sqlite.UpdatePlan(ctx, q, p, old); e != nil {
			return e
		}
		return audit.Record(ctx, q, audit.Event{ActorID: actor.ID, ObjectType: "plan", ObjectID: id, Action: "transition", Result: "ok", RequestID: common.NewID("req"), Detail: string(next)})
	})
}
func (s *Service) CreateZone(ctx context.Context, id, code, name, coast string) error {
	return s.tx(ctx, func(q sqlite.Tx) error { return sqlite.CreateZone(ctx, q, id, code, name, coast) })
}
func (s *Service) CreateShelter(ctx context.Context, sh shelter.Shelter) error {
	return s.tx(ctx, func(q sqlite.Tx) error { return sqlite.CreateShelter(ctx, q, sh) })
}
func (s *Service) AssignShelter(ctx context.Context, actor identity.User, shelterID, household string, people int) error {
	if !actor.Can("plan:write") {
		return common.ErrForbidden
	}
	return s.tx(ctx, func(q sqlite.Tx) error {
		sh, e := sqlite.GetShelter(ctx, q, shelterID)
		if e != nil {
			return e
		}
		old := sh.Version
		if e = sh.Assign(people); e != nil {
			return e
		}
		if e = sqlite.UpdateShelter(ctx, q, sh, old); e != nil {
			return e
		}
		_, e = q.ExecContext(ctx, "INSERT INTO shelter_assignments(id,shelter_id,household_ref,people,status,created_at) VALUES(?,?,?,?,?,?)", common.NewID("asn"), shelterID, household, people, "active", common.Now().UTC())
		return e
	})
}
func (s *Service) CreateSupply(ctx context.Context, id, zone, sku, name, unit string, reorder int) error {
	return s.tx(ctx, func(q sqlite.Tx) error { return sqlite.CreateSupply(ctx, q, id, zone, sku, name, unit, reorder) })
}
func (s *Service) CreateLot(ctx context.Context, l supply.Lot) error {
	return s.tx(ctx, func(q sqlite.Tx) error { return sqlite.CreateLot(ctx, q, l) })
}
func (s *Service) CreateRequest(ctx context.Context, r request.Request) error {
	return s.tx(ctx, func(q sqlite.Tx) error {
		if old, e := idempotency.Find(ctx, q, r.IdempotencyKey); e == nil && old.ExpiresAt.After(common.Now()) {
			return nil
		}
		if e := sqlite.CreateRequest(ctx, q, r); e != nil {
			return e
		}
		return idempotency.Save(ctx, q, r.IdempotencyKey, r.IdempotencyKey, 201, r)
	})
}
func (s *Service) ApproveRequest(ctx context.Context, actor identity.User, id, requestID string) error {
	if !actor.Can("request:approve") {
		return common.ErrForbidden
	}
	return s.tx(ctx, func(q sqlite.Tx) error {
		r, e := sqlite.GetRequest(ctx, q, id)
		if e != nil {
			return e
		}
		if r.Status != request.Submitted {
			return common.ErrConflict
		}
		if e = r.Transition(request.Approved); e != nil {
			return e
		}
		if e = sqlite.UpdateRequest(ctx, q, r, r.Version-1); e != nil {
			return e
		}
		for _, item := range r.Items {
			lots, e := sqlite.ListLots(ctx, q, item.SupplyID)
			if e != nil {
				return e
			}
			need := item.Quantity
			for _, lot := range lots {
				if need == 0 {
					break
				}
				take := lot.Free(common.Now())
				if take > need {
					take = need
				}
				if take == 0 {
					continue
				}
				old := lot.Version
				if e = lot.Reserve(take, common.Now()); e != nil {
					return e
				}
				if e = sqlite.UpdateLot(ctx, q, lot, old); e != nil {
					return e
				}
				if _, e = q.ExecContext(ctx, "INSERT INTO allocations(id,request_id,lot_id,quantity,status,created_at) VALUES(?,?,?,?,?,?)", common.NewID("alloc"), r.ID, lot.ID, take, "reserved", common.Now().UTC()); e != nil {
					return e
				}
				if e = sqlite.SetAllocated(ctx, q, r.ID, item.SupplyID, take); e != nil {
					return e
				}
				need -= take
			}
			if need > 0 {
				return common.ErrConflict
			}
		}
		r.Status = request.Allocated
		r.Version++
		if e = sqlite.UpdateRequest(ctx, q, r, r.Version-1); e != nil {
			return e
		}
		return audit.Record(ctx, q, audit.Event{ActorID: actor.ID, ObjectType: "request", ObjectID: id, Action: "approve", Result: "ok", RequestID: requestID, Detail: "stock reserved"})
	})
}
func (s *Service) CreateDispatch(ctx context.Context, actor identity.User, d dispatch.Dispatch) error {
	if !actor.Can("dispatch:write") {
		return common.ErrForbidden
	}
	return s.tx(ctx, func(q sqlite.Tx) error { return sqlite.CreateDispatch(ctx, q, d) })
}
func (s *Service) UpdateDispatch(ctx context.Context, actor identity.User, id string, op string) error {
	if !actor.Can("dispatch:write") {
		return common.ErrForbidden
	}
	return s.tx(ctx, func(q sqlite.Tx) error {
		d, e := sqlite.GetDispatch(ctx, q, id)
		if e != nil {
			return e
		}
		old := d.Version
		switch op {
		case "start":
			e = d.Start()
		case "deliver":
			e = d.Deliver()
		case "retry":
			e = d.Retry("operator request")
		default:
			e = common.ErrInvalid
		}
		if e != nil {
			return e
		}
		if e = sqlite.UpdateDispatch(ctx, q, d, old); e != nil {
			return e
		}
		return audit.Record(ctx, q, audit.Event{ActorID: actor.ID, ObjectType: "dispatch", ObjectID: id, Action: op, Result: "ok", RequestID: common.NewID("req"), Detail: string(d.Status)})
	})
}
func (s *Service) CreateIncident(ctx context.Context, actor identity.User, i incident.Incident) error {
	if !actor.Can("incident:write") {
		return common.ErrForbidden
	}
	return s.tx(ctx, func(q sqlite.Tx) error { return sqlite.CreateIncident(ctx, q, i) })
}
func (s *Service) TransitionIncident(ctx context.Context, actor identity.User, id string, next incident.Status) error {
	if !actor.Can("incident:write") {
		return common.ErrForbidden
	}
	return s.tx(ctx, func(q sqlite.Tx) error {
		var i incident.Incident
		var st string
		e := q.QueryRowContext(ctx, "SELECT id,plan_id,reporter_id,title,severity,status,context_note FROM incidents WHERE id=?", id).Scan(&i.ID, &i.PlanID, &i.ReporterID, &i.Title, &i.Severity, &st, &i.ContextNote)
		if e != nil {
			return e
		}
		i.Status = incident.Status(st)
		if e = i.Transition(next); e != nil {
			return e
		}
		if e = sqlite.UpdateIncident(ctx, q, i); e != nil {
			return e
		}
		return sqlite.AddIncidentAction(ctx, q, common.NewID("act"), id, actor.ID, string(next), "status transition")
	})
}
