package http_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gigavard/AIAdoption/ExTodoGolang/internal/config"
	"github.com/gigavard/AIAdoption/ExTodoGolang/internal/http"
	"github.com/gigavard/AIAdoption/ExTodoGolang/pkg/logger"
)

func TestHealthEndpoint(t *testing.T) {
	log := logger.New()
	cfg := &config.Config{HTTPAddr: ":8080"}
	server := http.NewServer(log, cfg)

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	server.Handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", ct)
	}
}
