package httpapi

import (
	"encoding/json"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/common"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/dispatch"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/identity"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/monitoring"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/service"
	"net/http"
	"strings"
)

type handler struct{ svc *service.Service }

func (h *handler) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, map[string]string{"status": "ok"})
}
func (h *handler) ready(w http.ResponseWriter, r *http.Request) {
	if err := h.svc.DB.SQL.PingContext(r.Context()); err != nil {
		writeErr(w, 503, err)
		return
	}
	writeJSON(w, 200, map[string]string{"status": "ready"})
}
func (h *handler) login(w http.ResponseWriter, r *http.Request) {
	var in struct{ Username, Password string }
	if json.NewDecoder(r.Body).Decode(&in) != nil {
		writeErr(w, 400, common.ErrInvalid)
		return
	}
	s, e := h.svc.Login(r.Context(), in.Username, in.Password)
	if e != nil {
		writeErr(w, 401, e)
		return
	}
	writeJSON(w, 200, map[string]any{"session_id": s.ID, "expires_at": s.ExpiresAt})
}
func (h *handler) dispatch(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/v1/")
	u, _ := currentUser(r)
	switch {
	case path == "logout" && r.Method == "POST":
		e := h.svc.Logout(r.Context(), r.Header.Get("X-Session-ID"))
		if e != nil {
			writeErr(w, 400, e)
			return
		}
		writeJSON(w, 204, nil)
	case path == "dispatches" && r.Method == "POST":
		var in struct{ RequestID, Driver, Vehicle string }
		if json.NewDecoder(r.Body).Decode(&in) != nil || in.RequestID == "" {
			writeErr(w, 400, common.ErrInvalid)
			return
		}
		d := dispatch.Dispatch{ID: common.NewID("dsp"), RequestID: in.RequestID, Driver: in.Driver, Vehicle: in.Vehicle, Status: dispatch.Queued, NextAttemptAt: common.Now(), Version: 1}
		if e := h.svc.CreateDispatch(r.Context(), u, d); e != nil {
			writeErr(w, 403, e)
			return
		}
		writeJSON(w, 201, d)
	case strings.HasPrefix(path, "dispatches/") && r.Method == "POST":
		id := strings.TrimPrefix(path, "dispatches/")
		if e := h.svc.UpdateDispatch(r.Context(), u, id, r.URL.Query().Get("op")); e != nil {
			writeErr(w, 409, e)
			return
		}
		writeJSON(w, 200, map[string]string{"status": "updated"})
	case path == "monitoring/assess" && r.Method == "POST":
		var in monitoring.Sample
		if json.NewDecoder(r.Body).Decode(&in) != nil {
			writeErr(w, 400, common.ErrInvalid)
			return
		}
		decision, e := h.svc.AssessEnvironment(r.Context(), u, in)
		if e != nil {
			writeErr(w, 422, e)
			return
		}
		writeJSON(w, 200, decision)
	default:
		writeErr(w, 404, common.ErrNotFound)
	}
}

var _ identity.Role
