package main

import (
	"context"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/config"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/httpapi"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/observability"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/service"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/storage/sqlite"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/worker"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := config.Load()
	logger := observability.NewLogger()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db, e := sqlite.Open(ctx, cfg.DatabasePath)
	if e != nil {
		logger.Error("database open failed", "error", e)
		os.Exit(1)
	}
	defer db.Close()
	svc := service.New(db)
	w := worker.New(db, func(context.Context, string) error { return nil }, logger)
	w.Start()
	defer w.Stop(context.Background())
	server := httpapi.New(cfg.HTTPAddr, svc)
	go func() {
		if e := server.Start(); e != nil {
			logger.Error("http stopped", "error", e)
		}
	}()
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	shutdown, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()
	_ = server.Shutdown(shutdown)
}
