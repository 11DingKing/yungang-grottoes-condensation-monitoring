package httpapi

import (
	"context"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/service"
	"net/http"
	"time"
)

type Server struct {
	HTTP    *http.Server
	Service *service.Service
}

func New(addr string, svc *service.Service) *Server {
	h := NewRouter(svc)
	return &Server{HTTP: &http.Server{Addr: addr, Handler: h, ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 15 * time.Second, WriteTimeout: 15 * time.Second, IdleTimeout: 60 * time.Second}, Service: svc}
}
func (s *Server) Start() error                       { return s.HTTP.ListenAndServe() }
func (s *Server) Shutdown(ctx context.Context) error { return s.HTTP.Shutdown(ctx) }
