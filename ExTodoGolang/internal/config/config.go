package config

import (
	"os"
)

type Config struct {
	Environment string
	HTTPAddr    string
	DBPath      string
	LogLevel    string
	OTLPEndpoint string
}

func Load() *Config {
	return &Config{
		Environment: getEnv("ENV", "development"),
		HTTPAddr:    getEnv("HTTP_ADDR", ":8080"),
		DBPath:      getEnv("DB_PATH", "todo.db"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
		OTLPEndpoint: getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", "http://localhost:4317"),
	}
}

func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}
