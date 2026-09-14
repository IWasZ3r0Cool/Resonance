package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/IWasZ3r0Cool/Resonance/internal/config"
	"github.com/IWasZ3r0Cool/Resonance/internal/database"
)

func TestHealth(t *testing.T) {
	db, err := database.Open(context.Background(), "file::memory:?cache=shared&_pragma=foreign_keys(1)")
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	defer db.Close()

	cfg := config.Config{WebOrigin: "http://localhost:5173"}
	request := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	response := httptest.NewRecorder()
	New(cfg, db, "test").ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", response.Code, http.StatusOK, response.Body.String())
	}
	var body map[string]string
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body["version"] != "test" || body["database"] != "ok" {
		t.Fatalf("unexpected response: %#v", body)
	}
}

func TestCapabilitiesDoesNotExposeSecrets(t *testing.T) {
	db, err := database.Open(context.Background(), "file::memory:?cache=shared&_pragma=foreign_keys(1)")
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	defer db.Close()
	cfg := config.Config{
		WebOrigin: "http://localhost:5173",
		LLM: config.LLM{
			Provider:     "openai",
			Model:        "test-model",
			OpenAIAPIKey: "super-secret",
		},
	}

	request := httptest.NewRequest(http.MethodGet, "/api/v1/meta/capabilities", nil)
	response := httptest.NewRecorder()
	New(cfg, db, "test").ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}
	if got := response.Body.String(); contains(got, "super-secret") {
		t.Fatalf("response exposed API key: %s", got)
	}
}

func contains(value, substring string) bool {
	for i := 0; i+len(substring) <= len(value); i++ {
		if value[i:i+len(substring)] == substring {
			return true
		}
	}
	return false
}
