package llm

import (
	"fmt"

	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/tools"
)

// NewAgentExecutor creates an agent using the provided system prompt and tools,
// optionally attaching the supplied memory, and returns its executor.
//
// - systemPrompt: system instructions for the agent
// - toolList:     tools the agent can call
// - mem:          optional memory (pass nil to use default)
func NewAgentExecutor(systemPrompt string, toolList []tools.Tool, mem schema.Memory) (agents.Executor, error) {
	if systemPrompt == "" {
		return agents.Executor{}, fmt.Errorf("systemPrompt is empty")
	}

	if toolList == nil {
		return agents.Executor{}, fmt.Errorf("toolList is nil")
	}

	llmClient, err := NewOpenAI()
	if err != nil {
		return agents.Executor{}, err
	}

	// Create an OpenAI Functions-style agent configured with the system prompt.
	agent := agents.NewOpenAIFunctionsAgent(
		llmClient,
		toolList,
		agents.NewOpenAIOption().WithSystemMessage(systemPrompt),
	)

	// Build executor options, attaching memory if provided.
	execOpts := make([]agents.CreationOption, 0, 1)
	if mem != nil {
		execOpts = append(execOpts, agents.WithMemory(mem))
	}

	exec := agents.NewExecutor(agent, toolList, execOpts...)
	return exec, nil
}
