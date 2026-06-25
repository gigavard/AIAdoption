package http

import (
	"log/slog"
	"net/http"

	"github.com/gigavard/AIAdoption/ExTodoGolang/internal/config"
	"github.com/gigavard/AIAdoption/ExTodoGolang/internal/domain"
)

type Server struct {
	*http.Server
	log *slog.Logger
}

func NewServer(log *slog.Logger, cfg *config.Config, repo domain.Repository) *Server {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", serveIndex)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/public/static"))))

	// Health check
	mux.HandleFunc("/health", healthHandler)

	// Handler routes
	handler := NewHandler(log, repo)
	handler.RegisterRoutes(mux)

	return &Server{
		Server: &http.Server{
			Addr:              cfg.HTTPAddr,
			Handler:           mux,
			ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		},
		log: log,
	}
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.ServeFile(w, r, "web/templates/index.html")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}
