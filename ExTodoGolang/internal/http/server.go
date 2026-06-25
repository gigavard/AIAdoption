package http

import (
	"log/slog"
	"net/http"

	"github.com/gigavard/AIAdoption/ExTodoGolang/internal/config"
)

type Server struct {
	*http.Server
	log *slog.Logger
}

func NewServer(log *slog.Logger, cfg *config.Config) *Server {
	mux := http.NewServeMux()

	// Routes will be added in SPEC-003
	mux.HandleFunc("/health", healthHandler)

	return &Server{
		Server: &http.Server{
			Addr:    cfg.HTTPAddr,
			Handler: mux,
		},
		log: log,
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}
