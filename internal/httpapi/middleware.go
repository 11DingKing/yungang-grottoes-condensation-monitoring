package httpapi

import (
	"context"
	"encoding/json"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/common"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/identity"
	"net/http"
	"runtime/debug"
)

type ctxKey string

const userKey ctxKey = "user"

func withUser(ctx context.Context, u identity.User) context.Context {
	return context.WithValue(ctx, userKey, u)
}
func currentUser(r *http.Request) (identity.User, bool) {
	u, ok := r.Context().Value(userKey).(identity.User)
	return u, ok
}
func auth(next http.Handler, svc Authenticator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if len(token) < 8 || token[:7] != "Bearer " {
			writeErr(w, http.StatusUnauthorized, common.ErrForbidden)
			return
		}
		u, e := svc.Authenticate(r.Context(), token[7:])
		if e != nil {
			writeErr(w, http.StatusUnauthorized, e)
			return
		}
		next.ServeHTTP(w, r.WithContext(withUser(r.Context(), u)))
	})
}

type Authenticator interface {
	Authenticate(context.Context, string) (identity.User, error)
}

func recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if v := recover(); v != nil {
				_ = debug.Stack()
				writeErr(w, http.StatusInternalServerError, common.ErrUnavailable)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
func requestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if id == "" {
			id = common.NewID("http")
		}
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r)
	})
}
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
func writeErr(w http.ResponseWriter, status int, e error) {
	code := "internal_error"
	if e == common.ErrForbidden {
		code = "forbidden"
	}
	if e == common.ErrConflict {
		code = "conflict"
	}
	if e == common.ErrInvalid {
		code = "invalid_request"
	}
	writeJSON(w, status, map[string]any{"error": map[string]string{"code": code, "message": e.Error()}})
}
