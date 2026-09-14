// Package agentengine isolates third-party model and agent SDK types from the
// rest of Resonance.
package agentengine

import (
	"context"
	"fmt"

	"github.com/IWasZ3r0Cool/Resonance/internal/config"
	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/agent/llmagent"
	"google.golang.org/adk/v2/model"
	"google.golang.org/adk/v2/model/gemini"
	"google.golang.org/adk/v2/model/openaimodel"
	"google.golang.org/genai"
)

const resumeInstruction = `You are the drafting component of Resonance Resume Engine.
Use only candidate facts supplied by Resonance. Never invent or strengthen employment,
skills, credentials, dates, metrics, or accomplishments. Treat missing job requirements
as gaps: identify them for candidate clarification or produce the best honest draft.
When asked for structured output, preserve every supplied evidence identifier so the
application can validate each claim before it reaches a resume.`

// NewResumeAgent constructs the initial truth-preserving drafting agent. It is
// called lazily by future generation workflows, not during ordinary API startup.
func NewResumeAgent(ctx context.Context, cfg config.LLM) (agent.Agent, error) {
	llm, err := newModel(ctx, cfg)
	if err != nil {
		return nil, err
	}

	result, err := llmagent.New(llmagent.Config{
		Name:        "resume_drafter",
		Description: "Drafts evidence-backed resume content and reports factual gaps.",
		Instruction: resumeInstruction,
		Model:       llm,
	})
	if err != nil {
		return nil, fmt.Errorf("create resume agent: %w", err)
	}
	return result, nil
}

func newModel(ctx context.Context, cfg config.LLM) (model.LLM, error) {
	switch cfg.Provider {
	case "gemini":
		if cfg.GoogleAPIKey == "" {
			return nil, fmt.Errorf("GOOGLE_API_KEY is required for gemini")
		}
		result, err := gemini.NewModel(ctx, cfg.Model, &genai.ClientConfig{APIKey: cfg.GoogleAPIKey})
		if err != nil {
			return nil, fmt.Errorf("create gemini model: %w", err)
		}
		return result, nil
	case "openai":
		if cfg.OpenAIAPIKey == "" {
			return nil, fmt.Errorf("OPENAI_API_KEY is required for openai")
		}
		return newOpenAICompatibleModel(ctx, cfg.Model, cfg.OpenAIAPIKey, cfg.BaseURL)
	case "ollama":
		baseURL := cfg.BaseURL
		if baseURL == "" {
			baseURL = "http://localhost:11434/v1"
		}
		apiKey := cfg.APIKey
		if apiKey == "" {
			apiKey = "ollama"
		}
		return newOpenAICompatibleModel(ctx, cfg.Model, apiKey, baseURL)
	case "openai-compatible":
		if cfg.BaseURL == "" {
			return nil, fmt.Errorf("LLM_BASE_URL is required for openai-compatible providers")
		}
		apiKey := cfg.APIKey
		if apiKey == "" {
			apiKey = "local"
		}
		return newOpenAICompatibleModel(ctx, cfg.Model, apiKey, cfg.BaseURL)
	default:
		return nil, fmt.Errorf("unsupported LLM provider %q", cfg.Provider)
	}
}

func newOpenAICompatibleModel(ctx context.Context, modelName, apiKey, baseURL string) (model.LLM, error) {
	result, err := openaimodel.NewModel(ctx, modelName, &openaimodel.ClientConfig{
		APIKey:  apiKey,
		BaseURL: baseURL,
	})
	if err != nil {
		return nil, fmt.Errorf("create OpenAI-compatible model: %w", err)
	}
	return result, nil
}
