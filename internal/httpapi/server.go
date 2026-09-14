// Package httpapi implements the public REST surface defined in api/openapi.json.
package httpapi

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	apicontract "github.com/IWasZ3r0Cool/Resonance/api"
	"github.com/IWasZ3r0Cool/Resonance/internal/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// Server owns HTTP handlers and their local dependencies.
type Server struct {
	config  config.Config
	db      *sql.DB
	version string
}

// New builds the API router.
func New(cfg config.Config, db *sql.DB, version string) http.Handler {
	s := &Server{config: cfg, db: db, version: version}
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(30 * time.Second))
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{cfg.WebOrigin},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodOptions},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Request-ID"},
		ExposedHeaders:   []string{"X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	router.Get("/openapi.json", s.openapi)
	router.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", s.health)
		r.Get("/meta/capabilities", s.capabilities)
	})
	return router
}

func (s *Server) openapi(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(apicontract.OpenAPI)
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second)
	defer cancel()
	if err := s.db.PingContext(ctx); err != nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "database unavailable"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{
		"status":   "ok",
		"database": "ok",
		"version":  s.version,
	})
}

func (s *Server) capabilities(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"agentSdk":          "google-adk-go",
		"database":          "sqlite",
		"llmConfigured":     s.config.LLMConfigured(),
		"llmModel":          s.config.LLM.Model,
		"llmProvider":       s.config.LLM.Provider,
		"localModelSupport": true,
	})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
