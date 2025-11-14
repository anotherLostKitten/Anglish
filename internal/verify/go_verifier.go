package verify

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
)

// ErrEmptyResponse indicates the agent returned no files or an empty payload.
var ErrEmptyResponse = errors.New("verify: empty agent response")

// ErrNoGoFiles indicates the agent response omitted Go source files.
var ErrNoGoFiles = errors.New("verify: no Go files in agent response")

// AgentFile mirrors the structure emitted by the call space agent system prompt.
type AgentFile struct {
	Filename  string   `json:"filename"`
	FileLines []string `json:"filelines"`
}

// ValidateGoResponse ensures the agent JSON response decodes and contains Go files
// whose contents parse successfully. It validates the line-based format so the
// consumer can materialize files without ambiguity.
func ValidateGoResponse(raw []byte) error {
	payload := bytes.TrimSpace(raw)
	if len(payload) == 0 {
		return ErrEmptyResponse
	}

	var files []AgentFile
	if err := json.Unmarshal(payload, &files); err != nil {
		return fmt.Errorf("verify: decode agent response: %w", err)
	}
	if len(files) == 0 {
		return ErrEmptyResponse
	}

	fset := token.NewFileSet()
	goFileCount := 0

	for _, file := range files {
		filename := strings.TrimSpace(file.Filename)
		if filename == "" {
			return errors.New("verify: file is missing a filename")
		}
		// Enforce one logical line per entry to match the prompt contract.
		for i, line := range file.FileLines {
			if strings.ContainsAny(line, "\r\n") {
				return fmt.Errorf("verify: %s line %d contains unexpected newline characters", filename, i+1)
			}
		}

		if filepath.Ext(filename) != ".go" {
			continue
		}
		goFileCount++

		if len(file.FileLines) == 0 {
			return fmt.Errorf("verify: Go file %s contains no content", filename)
		}

		source := strings.Join(file.FileLines, "\n")
		if _, err := parser.ParseFile(fset, filename, source, parser.AllErrors); err != nil {
			return fmt.Errorf("verify: parse %s: %w", filename, err)
		}
	}

	if goFileCount == 0 {
		return ErrNoGoFiles
	}

	return nil
}
