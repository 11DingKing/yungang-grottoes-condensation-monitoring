package httpapi

import (
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/service"
	"net/http"
)

func NewRouter(svc *service.Service) http.Handler {
	mux := http.NewServeMux()
	h := &handler{svc: svc}
	mux.HandleFunc("GET /healthz", h.health)
	mux.HandleFunc("GET /readyz", h.ready)
	mux.HandleFunc("POST /v1/login", h.login)
	protected := auth(mux, svc)
	_ = protected
	mux.Handle("/v1/", auth(http.HandlerFunc(h.dispatch), svc))
	return requestID(recoverer(mux))
}
