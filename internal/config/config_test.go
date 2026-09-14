package config

import "testing"

func TestLoadDefaults(t *testing.T) {
	for _, key := range []string{
		"APP_ENV", "HTTP_HOST", "HTTP_PORT", "WEB_ORIGIN", "DATABASE_URL",
		"LOG_LEVEL", "MAX_UPLOAD_BYTES", "LLM_PROVIDER", "LLM_MODEL",
		"LLM_BASE_URL", "LLM_API_KEY", "OPENAI_API_KEY", "GOOGLE_API_KEY",
	} {
		t.Setenv(key, "")
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if got, want := cfg.Address(), "127.0.0.1:8080"; got != want {
		t.Fatalf("Address() = %q, want %q", got, want)
	}
	if cfg.LLMConfigured() {
		t.Fatal("default OpenAI provider should not be configured without a key")
	}
}

func TestOllamaIsLocallyConfigurableWithoutASecret(t *testing.T) {
	t.Setenv("LLM_PROVIDER", "ollama")
	t.Setenv("LLM_BASE_URL", "")
	t.Setenv("LLM_MODEL", "qwen3:8b")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if !cfg.LLMConfigured() {
		t.Fatal("Ollama should be configurable without a real secret")
	}
}

func TestRejectsCompatibleProviderWithoutBaseURL(t *testing.T) {
	t.Setenv("LLM_PROVIDER", "openai-compatible")
	t.Setenv("LLM_BASE_URL", "")

	if _, err := Load(); err == nil {
		t.Fatal("Load() error = nil, want validation error")
	}
}
