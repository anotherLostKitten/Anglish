package llm

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

const (
	FuncSpaceAgentSystem = "func_space_agent_system"
	UiSpaceAgentSystem   = "ui_space_agent_system"
)

// ReadPrompt reads a specific prompt file from the prompts directory
// and returns its contents as a string.
// promptName should be the name of the file without the .md extension,
// e.g., "func_space_agent_system" or "ui_space_agent_system"
func ReadPrompt(promptName string) (string, error) {
	// Get the directory where this file is located
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("failed to get current file path")
	}
	baseDir := filepath.Dir(filename)
	promptsDir := filepath.Join(baseDir, "prompts")

	// Construct the full path to the prompt file
	promptPath := filepath.Join(promptsDir, promptName+".md")

	// Read the file
	content, err := os.ReadFile(promptPath)
	if err != nil {
		return "", fmt.Errorf("failed to read prompt file %s: %w", promptName, err)
	}

	return string(content), nil
}
