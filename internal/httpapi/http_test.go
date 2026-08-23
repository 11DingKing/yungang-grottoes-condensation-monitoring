package httpapi_test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/identity"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/httpapi"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/service"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/storage/sqlite"
	"net/http"
	"net/http/httptest"
	"testing"
)

func server(t *testing.T) http.Handler {
	db, e := sqlite.Open(context.Background(), ":memory:")
	if e != nil {
		t.Fatal(e)
	}
	t.Cleanup(func() { db.Close() })
	s := service.New(db)
	u := identity.NewUser("coordinator", identity.RoleCoordinator)
	u.PasswordHash = "pw"
	if e = s.CreateUser(context.Background(), u); e != nil {
		t.Fatal(e)
	}
	return httpapi.NewRouter(s)
}
func TestHealth(t *testing.T) {
	h := server(t)
	for _, path := range []string{"/healthz", "/readyz"} {
		r := httptest.NewRequest(http.MethodGet, path, nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		if w.Code != 200 {
			t.Fatalf("%s %d", path, w.Code)
		}
	}
}
func TestLoginReturnsSession(t *testing.T) {
	h := server(t)
	body, _ := json.Marshal(map[string]string{"username": "coordinator", "password": "pw"})
	r := httptest.NewRequest(http.MethodPost, "/v1/login", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != 200 {
		t.Fatal(w.Code, w.Body.String())
	}
	var out map[string]any
	if e := json.Unmarshal(w.Body.Bytes(), &out); e != nil {
		t.Fatal(e)
	}
	if out["session_id"] == nil {
		t.Fatal(out)
	}
	if w.Header().Get("X-Request-ID") == "" {
		t.Fatal("missing request id")
	}
}
func TestLoginBadPayload(t *testing.T) {
	h := server(t)
	r := httptest.NewRequest(http.MethodPost, "/v1/login", bytes.NewBufferString("{}"))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != 401 && w.Code != 400 {
		t.Fatal(w.Code)
	}
}
func TestProtectedEndpointNeedsBearer(t *testing.T) {
	h := server(t)
	r := httptest.NewRequest(http.MethodPost, "/v1/dispatches", bytes.NewBufferString(`{"request_id":"r","driver":"d","vehicle":"v"}`))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != 401 {
		t.Fatal(w.Code)
	}
}
func TestUnknownEndpoint(t *testing.T) {
	h := server(t)
	r := httptest.NewRequest(http.MethodGet, "/v1/nope", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != 401 {
		t.Fatal(w.Code)
	}
}
