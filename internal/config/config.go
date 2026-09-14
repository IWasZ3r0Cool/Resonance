// Package config owns environment parsing and validation.
package config

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const defaultDatabaseURL = "file:./data/resonance.db?_pragma=foreign_keys(1)&_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)"

// Config contains all process configuration. Secrets are deliberately kept in
// a nested server-only value and must never be serialized into API responses.
type Config struct {
	Environment    string
	HTTPHost       string
	HTTPPort       int
	WebOrigin      string
	DatabaseURL    string
	LogLevel       string
	MaxUploadBytes int64
	LLM            LLM
}

// LLM contains model-provider configuration used by the agent factory.
type LLM struct {
	Provider     string
	Model        string
	BaseURL      string
	APIKey       string
	OpenAIAPIKey string
	GoogleAPIKey string
}

// Load reads process environment variables and applies safe local defaults.
func Load() (Config, error) {
	port, err := intEnv("HTTP_PORT", 8080)
	if err != nil {
		return Config{}, err
	}
	maxUpload, err := int64Env("MAX_UPLOAD_BYTES", 10*1024*1024)
	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		Environment:    env("APP_ENV", "development"),
		HTTPHost:       env("HTTP_HOST", "127.0.0.1"),
		HTTPPort:       port,
		WebOrigin:      env("WEB_ORIGIN", "http://localhost:5173"),
		DatabaseURL:    env("DATABASE_URL", defaultDatabaseURL),
		LogLevel:       strings.ToLower(env("LOG_LEVEL", "info")),
		MaxUploadBytes: maxUpload,
		LLM: LLM{
			Provider:     strings.ToLower(env("LLM_PROVIDER", "openai")),
			Model:        env("LLM_MODEL", "gpt-5-mini"),
			BaseURL:      strings.TrimSpace(os.Getenv("LLM_BASE_URL")),
			APIKey:       strings.TrimSpace(os.Getenv("LLM_API_KEY")),
			OpenAIAPIKey: strings.TrimSpace(os.Getenv("OPENAI_API_KEY")),
			GoogleAPIKey: strings.TrimSpace(os.Getenv("GOOGLE_API_KEY")),
		},
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

// Validate rejects ambiguous or unsafe-to-run configuration.
func (c Config) Validate() error {
	if c.HTTPPort < 1 || c.HTTPPort > 65535 {
		return fmt.Errorf("HTTP_PORT must be between 1 and 65535")
	}
	if strings.TrimSpace(c.HTTPHost) == "" {
		return fmt.Errorf("HTTP_HOST must not be empty")
	}
	if strings.TrimSpace(c.DatabaseURL) == "" {
		return fmt.Errorf("DATABASE_URL must not be empty")
	}
	if c.MaxUploadBytes < 1 {
		return fmt.Errorf("MAX_UPLOAD_BYTES must be positive")
	}
	if _, ok := map[string]struct{}{"debug": {}, "info": {}, "warn": {}, "error": {}}[c.LogLevel]; !ok {
		return fmt.Errorf("LOG_LEVEL must be debug, info, warn, or error")
	}
	if _, ok := map[string]struct{}{"openai": {}, "gemini": {}, "ollama": {}, "openai-compatible": {}}[c.LLM.Provider]; !ok {
		return fmt.Errorf("LLM_PROVIDER must be openai, gemini, ollama, or openai-compatible")
	}
	if strings.TrimSpace(c.LLM.Model) == "" {
		return fmt.Errorf("LLM_MODEL must not be empty")
	}
	if c.LLM.Provider == "openai-compatible" && c.LLM.BaseURL == "" {
		return fmt.Errorf("LLM_BASE_URL is required for openai-compatible providers")
	}
	if c.LLM.BaseURL != "" {
		parsed, err := url.Parse(c.LLM.BaseURL)
		if err != nil || parsed.Scheme == "" || parsed.Host == "" {
			return fmt.Errorf("LLM_BASE_URL must be an absolute URL")
		}
	}
	return nil
}

// Address is the host:port value passed to http.Server.
func (c Config) Address() string {
	return net.JoinHostPort(c.HTTPHost, strconv.Itoa(c.HTTPPort))
}

// LLMConfigured reports whether model initialization has enough configuration
// to be attempted. It never makes a remote request.
func (c Config) LLMConfigured() bool {
	switch c.LLM.Provider {
	case "openai":
		return c.LLM.OpenAIAPIKey != ""
	case "gemini":
		return c.LLM.GoogleAPIKey != ""
	case "ollama":
		return true
	case "openai-compatible":
		return c.LLM.BaseURL != ""
	default:
		return false
	}
}

func env(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func intEnv(key string, fallback int) (int, error) {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback, nil
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("%s must be an integer: %w", key, err)
	}
	return parsed, nil
}

func int64Env(key string, fallback int64) (int64, error) {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback, nil
	}
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%s must be an integer: %w", key, err)
	}
	return parsed, nil
}
