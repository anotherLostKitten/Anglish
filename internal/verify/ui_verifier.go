package verify

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/robertkrimen/otto/parser"
	"github.com/tdewolff/parse/css"
	"golang.org/x/net/html"
)

var (
	// ErrNoHTMLFiles indicates the agent response omitted HTML source files.
	ErrNoHTMLFiles = errors.New("verify: no HTML files in agent response")
	// ErrNoCSSFiles indicates the agent response omitted CSS source files.
	ErrNoCSSFiles = errors.New("verify: no CSS files in agent response")
	// ErrNoJSFiles indicates the agent response omitted JavaScript source files.
	ErrNoJSFiles = errors.New("verify: no JavaScript files in agent response")
)

// ValidateHTMLResponse ensures any HTML files within the agent response can be parsed.
func ValidateHTMLResponse(raw []byte) error {
	files, err := decodeAgentFiles(raw)
	if err != nil {
		return err
	}

	found := false
	for _, file := range files {
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if ext != ".html" && ext != ".htm" {
			continue
		}
		found = true
		if len(file.FileLines) == 0 {
			return fmt.Errorf("verify: HTML file %s contains no content", file.Filename)
		}

		source := strings.Join(file.FileLines, "\n")
		if _, err := html.Parse(strings.NewReader(source)); err != nil {
			return fmt.Errorf("verify: parse HTML %s: %w", file.Filename, err)
		}
	}

	if !found {
		return ErrNoHTMLFiles
	}

	return nil
}

// ValidateCSSResponse ensures any CSS files within the agent response can be parsed.
func ValidateCSSResponse(raw []byte) error {
	files, err := decodeAgentFiles(raw)
	if err != nil {
		return err
	}

	found := false
	for _, file := range files {
		if strings.ToLower(filepath.Ext(file.Filename)) != ".css" {
			continue
		}
		found = true
		if len(file.FileLines) == 0 {
			return fmt.Errorf("verify: CSS file %s contains no content", file.Filename)
		}

		source := strings.Join(file.FileLines, "\n")
		if err := parseCSS(file.Filename, source); err != nil {
			return err
		}
	}

	if !found {
		return ErrNoCSSFiles
	}

	return nil
}

// ValidateJSResponse ensures any JS files within the agent response can be parsed.
func ValidateJSResponse(raw []byte) error {
	files, err := decodeAgentFiles(raw)
	if err != nil {
		return err
	}

	found := false
	for _, file := range files {
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if ext != ".js" && ext != ".mjs" && ext != ".cjs" {
			continue
		}
		found = true
		if len(file.FileLines) == 0 {
			return fmt.Errorf("verify: JavaScript file %s contains no content", file.Filename)
		}

		source := strings.Join(file.FileLines, "\n")
		if _, err := parser.ParseFile(nil, file.Filename, source, 0); err != nil {
			return fmt.Errorf("verify: parse JavaScript %s: %w", file.Filename, err)
		}
	}

	if !found {
		return ErrNoJSFiles
	}

	return nil
}

// ValidateUIResponse dispatches to the language-specific validators for each file
// within the agent payload and ensures the required UI files are present.
func ValidateUIResponse(raw []byte) error {
	files, err := decodeAgentFiles(raw)
	if err != nil {
		return err
	}

	htmlValidated := false
	cssValidated := false
	jsValidated := false

	for _, file := range files {
		payload, err := json.Marshal([]AgentFile{file})
		if err != nil {
			return fmt.Errorf("verify: marshal agent file %s: %w", file.Filename, err)
		}

		switch ext := strings.ToLower(filepath.Ext(file.Filename)); ext {
		case ".html", ".htm":
			if err := ValidateHTMLResponse(payload); err != nil {
				return err
			}
			htmlValidated = true
		case ".css":
			if err := ValidateCSSResponse(payload); err != nil {
				return err
			}
			cssValidated = true
		case ".js", ".mjs", ".cjs":
			if err := ValidateJSResponse(payload); err != nil {
				return err
			}
			jsValidated = true
		default:
			continue
		}
	}

	if !htmlValidated {
		return ErrNoHTMLFiles
	}
	if !cssValidated {
		return ErrNoCSSFiles
	}
	if !jsValidated {
		return ErrNoJSFiles
	}

	return nil
}

func decodeAgentFiles(raw []byte) ([]AgentFile, error) {
	payload := bytes.TrimSpace(raw)
	if len(payload) == 0 {
		return nil, ErrEmptyResponse
	}

	var files []AgentFile
	if err := json.Unmarshal(payload, &files); err != nil {
		return nil, fmt.Errorf("verify: decode agent response: %w", err)
	}
	if len(files) == 0 {
		return nil, ErrEmptyResponse
	}

	for _, file := range files {
		filename := strings.TrimSpace(file.Filename)
		if filename == "" {
			return nil, errors.New("verify: file is missing a filename")
		}
		for i, line := range file.FileLines {
			if strings.ContainsAny(line, "\r\n") {
				return nil, fmt.Errorf("verify: %s line %d contains unexpected newline characters", filename, i+1)
			}
		}
	}

	return files, nil
}

func parseCSS(filename, source string) error {
	parser := css.NewParser(strings.NewReader(source), false)
	for {
		grammar, _, _ := parser.Next()
		if grammar == css.ErrorGrammar {
			if err := parser.Err(); err != nil {
				if errors.Is(err, io.EOF) {
					return nil
				}
				return fmt.Errorf("verify: parse CSS %s: %w", filename, err)
			}
		}

		if err := parser.Err(); err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return fmt.Errorf("verify: parse CSS %s: %w", filename, err)
		}
	}
}
