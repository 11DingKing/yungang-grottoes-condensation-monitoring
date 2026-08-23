package sqlite

import (
	"context"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/common"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/incident"
)

func CreateIncident(ctx context.Context, q Tx, i incident.Incident) error {
	_, e := q.ExecContext(ctx, "INSERT INTO incidents(id,plan_id,reporter_id,title,severity,status,context_note,created_at,updated_at) VALUES(?,?,?,?,?,?,?,?,?)", i.ID, i.PlanID, i.ReporterID, i.Title, i.Severity, i.Status, i.ContextNote, common.Now().UTC(), common.Now().UTC())
	return e
}
func AddIncidentAction(ctx context.Context, q Tx, id, incidentID, actor, action, note string) error {
	_, e := q.ExecContext(ctx, "INSERT INTO incident_actions(id,incident_id,actor_id,action,note,created_at) VALUES(?,?,?,?,?,?)", id, incidentID, actor, action, note, common.Now().UTC())
	return e
}
func UpdateIncident(ctx context.Context, q Tx, i incident.Incident) error {
	_, e := q.ExecContext(ctx, "UPDATE incidents SET status=?,updated_at=? WHERE id=?", i.Status, common.Now().UTC(), i.ID)
	return e
}
