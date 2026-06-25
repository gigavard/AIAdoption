package config

import (
	"os"
	"time"
)

type Config struct {
	Environment       string
	HTTPAddr          string
	DBPath            string
	LogLevel          string
	OTLPEndpoint      string
	ShutdownTimeout   time.Duration
	ReadHeaderTimeout time.Duration
}

func Load() *Config {
	return &Config{
		Environment:       getEnv("ENV", "development"),
		HTTPAddr:          getEnv("HTTP_ADDR", ":8080"),
		DBPath:            getEnv("DB_PATH", "todo.db"),
		LogLevel:          getEnv("LOG_LEVEL", "info"),
		OTLPEndpoint:      getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", "http://localhost:4317"),
		ShutdownTimeout:   30 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}
}

func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}
