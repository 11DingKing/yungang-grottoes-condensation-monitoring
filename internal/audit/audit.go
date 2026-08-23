package audit

import (
	"context"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/common"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/storage/sqlite"
)

type Event struct{ ID, ActorID, ObjectType, ObjectID, Action, Result, RequestID, Detail string }

func Record(ctx context.Context, q sqlite.Tx, e Event) error {
	if e.ID == "" {
		e.ID = common.NewID("aud")
	}
	_, err := q.ExecContext(ctx, "INSERT INTO audit_events(id,actor_id,object_type,object_id,action,result,request_id,detail,created_at) VALUES(?,?,?,?,?,?,?,?,?)", e.ID, e.ActorID, e.ObjectType, e.ObjectID, e.Action, e.Result, e.RequestID, e.Detail, common.Now().UTC())
	return err
}
func List(ctx context.Context, q sqlite.Tx, objectType, objectID string) ([]Event, error) {
	rows, e := q.QueryContext(ctx, "SELECT id,COALESCE(actor_id,''),object_type,object_id,action,result,request_id,detail FROM audit_events WHERE object_type=? AND object_id=? ORDER BY created_at", objectType, objectID)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []Event{}
	for rows.Next() {
		var x Event
		if e = rows.Scan(&x.ID, &x.ActorID, &x.ObjectType, &x.ObjectID, &x.Action, &x.Result, &x.RequestID, &x.Detail); e != nil {
			return nil, e
		}
		out = append(out, x)
	}
	return out, rows.Err()
}
