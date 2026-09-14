package apicontract

import (
	"encoding/json"
	"testing"
)

func TestOpenAPIIsValidJSONWithExpectedVersion(t *testing.T) {
	var document struct {
		OpenAPI string         `json:"openapi"`
		Paths   map[string]any `json:"paths"`
	}
	if err := json.Unmarshal(OpenAPI, &document); err != nil {
		t.Fatalf("OpenAPI document is not valid JSON: %v", err)
	}
	if document.OpenAPI != "3.1.0" {
		t.Fatalf("OpenAPI version = %q, want 3.1.0", document.OpenAPI)
	}
	for _, path := range []string{"/api/v1/health", "/api/v1/meta/capabilities", "/openapi.json"} {
		if _, ok := document.Paths[path]; !ok {
			t.Errorf("OpenAPI document is missing %s", path)
		}
	}
}
