package agentengine

import (
	"context"
	"testing"

	"github.com/IWasZ3r0Cool/Resonance/internal/config"
)

func TestNewResumeAgentRequiresOpenAIKey(t *testing.T) {
	_, err := NewResumeAgent(context.Background(), config.LLM{
		Provider: "openai",
		Model:    "gpt-5-mini",
	})
	if err == nil {
		t.Fatal("NewResumeAgent() error = nil, want missing key error")
	}
}

func TestNewResumeAgentSupportsLocalResponsesEndpoint(t *testing.T) {
	result, err := NewResumeAgent(context.Background(), config.LLM{
		Provider: "ollama",
		Model:    "qwen3:8b",
	})
	if err != nil {
		t.Fatalf("NewResumeAgent() error = %v", err)
	}
	if result == nil {
		t.Fatal("NewResumeAgent() returned nil agent")
	}
}
