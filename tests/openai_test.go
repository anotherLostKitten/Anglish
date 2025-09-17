//go:build vllm

package tests

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	angllm "github.com/anotherLostKitten/Anglish/internal/llm"
)

func TestNewOpenAI(t *testing.T) {

	require.NoError(t, os.Setenv("OPENAI_BASE_URL", "http://localhost:8000/v1"))
	require.NoError(t, os.Setenv("OPENAI_API_KEY", "EMPTY"))
	require.NoError(t, os.Setenv("OPENAI_MODEL", "google/gemma3-4b-it"))

	model, err := angllm.NewOpenAI()
	require.NoError(t, err)
	require.NotNil(t, model)
}
