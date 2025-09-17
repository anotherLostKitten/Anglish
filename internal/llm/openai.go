package llm

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

// NewOpenAI initializes and returns an OpenAI-compatible model client using
// configuration from environment variables. It attempts to load a local .env
// file (if present) using godotenv before reading variables.
//
// Required variables (see env.example):
// - OPENAI_API_KEY: API key for the OpenAI-compatible server
// - OPENAI_MODEL: Default model name to use
// Optional variables:
// - OPENAI_BASE_URL: Base URL for an OpenAI-compatible endpoint (e.g., vLLM)
func NewOpenAI() (llms.Model, error) {
	_ = godotenv.Load()

	apiKey := os.Getenv("OPENAI_API_KEY")
	model := os.Getenv("OPENAI_MODEL")
	baseURL := os.Getenv("OPENAI_BASE_URL")

	if apiKey == "" {
		return nil, fmt.Errorf("missing OPENAI_API_KEY environment variable")
	}
	if model == "" {
		return nil, fmt.Errorf("missing OPENAI_MODEL environment variable")
	}

	opts := []openai.Option{
		openai.WithModel(model),
	}
	if baseURL != "" {
		opts = append(opts, openai.WithBaseURL(baseURL))
	}

	// Note: The OpenAI client reads the API key from the OPENAI_API_KEY env var.
	client, err := openai.New(opts...)
	if err != nil {
		return nil, err
	}

	return client, nil
}
