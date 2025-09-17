//go:build vllm

package tests

import (
	"context"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	angllm "github.com/anotherLostKitten/Anglish/internal/llm"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/memory"
	"github.com/tmc/langchaingo/tools"
)

type echoTool struct{}

func (e echoTool) Name() string                                         { return "ECHO" }
func (e echoTool) Description() string                                  { return "Echo back the input" }
func (e echoTool) Call(_ context.Context, input string) (string, error) { return input, nil }

func TestNewAgentExecutor_Errors(t *testing.T) {
	// Empty system prompt should error before any env var checks.
	_, err := angllm.NewAgentExecutor("", []tools.Tool{echoTool{}}, nil)
	require.Error(t, err)

	// Nil tool list should error before any env var checks.
	_, err = angllm.NewAgentExecutor("system", nil, nil)
	require.Error(t, err)
}

func TestNewAgentExecutor_Success_NoMemory(t *testing.T) {
	require.NoError(t, os.Setenv("OPENAI_BASE_URL", "http://localhost:8000/v1"))
	require.NoError(t, os.Setenv("OPENAI_API_KEY", "EMPTY"))
	require.NoError(t, os.Setenv("OPENAI_MODEL", "google/gemma-3-4b-it"))

	exec, err := angllm.NewAgentExecutor("You are helpful assistant.", []tools.Tool{echoTool{}}, nil)
	require.NoError(t, err)

	// Basic sanity checks
	require.NotNil(t, exec.GetMemory())
	require.Len(t, exec.Tools, 1)
	require.Equal(t, "ECHO", exec.Tools[0].Name())

	// Expect input key to include "input"
	keys := exec.GetInputKeys()
	require.Contains(t, keys, "input")
}

func TestNewAgentExecutor_Success_WithMemory(t *testing.T) {
	require.NoError(t, os.Setenv("OPENAI_BASE_URL", "http://localhost:8000/v1"))
	require.NoError(t, os.Setenv("OPENAI_API_KEY", "EMPTY"))
	require.NoError(t, os.Setenv("OPENAI_MODEL", "google/gemma-3-4b-it"))

	mem := memory.NewSimple()
	exec, err := angllm.NewAgentExecutor("System prompt.", []tools.Tool{echoTool{}}, mem)
	require.NoError(t, err)

	gotMem := exec.GetMemory()
	require.NotNil(t, gotMem)
	// Verify the memory type and value
	_, ok := gotMem.(memory.Simple)
	require.True(t, ok)
	require.True(t, reflect.DeepEqual(mem, gotMem))
}

func TestNewAgentExecutor_RunChainResponds(t *testing.T) {
	require.NoError(t, os.Setenv("OPENAI_BASE_URL", "http://localhost:8000/v1"))
	require.NoError(t, os.Setenv("OPENAI_API_KEY", "EMPTY"))
	require.NoError(t, os.Setenv("OPENAI_MODEL", "google/gemma-3-4b-it"))

	exec, err := angllm.NewAgentExecutor("You are a helpful assistant.", []tools.Tool{echoTool{}}, nil)
	if err != nil {
		t.Fatalf("executor creation failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	answer, err := chains.Run(ctx, exec, "Say OK.")
	if err != nil {
		// If the local vLLM server isn't running, skip to keep the suite green by default.
		if os.IsTimeout(err) || (err != nil && (contains(err.Error(), "connection refused") || contains(err.Error(), "dial tcp") || contains(err.Error(), "status code: 404") || contains(err.Error(), "does not exist") || contains(err.Error(), "Not Found"))) {
			t.Skipf("skipping run test: LLM endpoint unavailable: %v", err)
		}
		t.Fatalf("run failed: %v", err)
	}

	if len(answer) == 0 {
		t.Fatalf("expected non-empty answer")
	}
}

func contains(s, substr string) bool { return len(s) >= len(substr) && (indexOf(s, substr) >= 0) }

func indexOf(s, substr string) int {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
