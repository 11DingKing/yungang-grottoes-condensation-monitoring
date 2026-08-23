package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	DatabasePath, HTTPAddr string
	SessionTTL             time.Duration
	DispatchWorkers        int
}

func Load() Config {
	ttl, _ := time.ParseDuration(env("SESSION_TTL", "12h"))
	workers, _ := strconv.Atoi(env("DISPATCH_WORKERS", "2"))
	if workers < 1 {
		workers = 1
	}
	return Config{DatabasePath: env("DATABASE_PATH", "./yungang.db"), HTTPAddr: env("HTTP_ADDR", ":8080"), SessionTTL: ttl, DispatchWorkers: workers}
}
func env(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
