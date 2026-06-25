package observability

import (
	"log/slog"
	"net/http"

	"github.com/gigavard/AIAdoption/ExTodoGolang/pkg/logger"
)

type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
	written    int64
}

func (w *ResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *ResponseWriter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.written += int64(n)
	return n, err
}

func LoggingMiddleware(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rw := &ResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
			next.ServeHTTP(rw, r)

			log.Info("HTTP request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", rw.statusCode,
				"bytes_written", rw.written,
			)
		})
	}
}

// MetricsCollector tracks request metrics for Prometheus
type MetricsCollector struct {
	requests map[string]int64
}

func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		requests: make(map[string]int64),
	}
}

func (m *MetricsCollector) RecordRequest(method, path string) {
	key := method + " " + path
	m.requests[key]++
}

func (m *MetricsCollector) GetMetrics() map[string]int64 {
	return m.requests
}
