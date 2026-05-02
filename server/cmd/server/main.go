package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"opspulse/server/internal/api"
)

func main() {
	cfg := api.Config{
		Address:        envOrDefault("OPS_SERVER_ADDR", ":8080"),
		Token:          envOrDefault("OPS_AGENT_TOKEN", "change-me"),
		DatabasePath:   envOrDefault("OPS_DB_PATH", "./opspulse.db"),
		OfflineTimeout: envDurationOrDefault("OPS_OFFLINE_TIMEOUT", 45*time.Second),
	}

	srv, err := api.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer srv.Close()

	log.Printf("server listening on %s", cfg.Address)
	log.Fatal(http.ListenAndServe(cfg.Address, srv.Handler()))
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func envDurationOrDefault(key string, fallback time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		duration, err := time.ParseDuration(value)
		if err == nil {
			return duration
		}
	}
	return fallback
}
