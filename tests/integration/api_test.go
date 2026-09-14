//go:build integration

package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/IWasZ3r0Cool/Resonance/internal/config"
	"github.com/IWasZ3r0Cool/Resonance/internal/database"
	"github.com/IWasZ3r0Cool/Resonance/internal/httpapi"
)

func TestFreshDatabaseAndHTTPStack(t *testing.T) {
	dsn := "file:" + t.TempDir() + "/integration.db?_pragma=foreign_keys(1)"
	db, err := database.Open(context.Background(), dsn)
	if err != nil {
		t.Fatalf("open migrated database: %v", err)
	}
	defer db.Close()

	cfg := config.Config{
		WebOrigin: "http://localhost:5173",
		LLM: config.LLM{
			Provider: "ollama",
			Model:    "qwen3:8b",
		},
	}
	handler := httpapi.New(cfg, db, "integration")

	request := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("health status = %d, want 200", response.Code)
	}

	var health struct {
		Status  string `json:"status"`
		Version string `json:"version"`
	}
	if err := json.NewDecoder(response.Body).Decode(&health); err != nil {
		t.Fatalf("decode health response: %v", err)
	}
	if health.Status != "ok" || health.Version != "integration" {
		t.Fatalf("unexpected health response: %#v", health)
	}

	request = httptest.NewRequest(http.MethodGet, "/api/v1/meta/capabilities", nil)
	response = httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	var capabilities struct {
		LLMConfigured bool   `json:"llmConfigured"`
		LLMProvider   string `json:"llmProvider"`
	}
	if err := json.NewDecoder(response.Body).Decode(&capabilities); err != nil {
		t.Fatalf("decode capabilities response: %v", err)
	}
	if !capabilities.LLMConfigured || capabilities.LLMProvider != "ollama" {
		t.Fatalf("unexpected capabilities response: %#v", capabilities)
	}

	for _, table := range []string{"resume_bullets", "bullet_skills", "bullet_tags", "resume_versions"} {
		var tableCount int
		if err := db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type = 'table' AND name = ?", table).Scan(&tableCount); err != nil {
			t.Fatalf("query migrated schema for %s: %v", table, err)
		}
		if tableCount != 1 {
			t.Fatalf("%s table count = %d, want 1", table, tableCount)
		}
	}
}
