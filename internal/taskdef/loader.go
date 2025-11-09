package taskdef

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

// LoadMode represents the input mode (task or container)
type LoadMode string

const (
	ModeTask      LoadMode = "task"
	ModeContainer LoadMode = "container"
)

// removeJSONComments removes single-line (//) and multi-line (/* */) comments from JSONC
func removeJSONComments(data []byte) []byte {
	var result bytes.Buffer
	scanner := bufio.NewScanner(bytes.NewReader(data))

	inMultiLineComment := false

	for scanner.Scan() {
		line := scanner.Text()

		// Handle multi-line comments
		if inMultiLineComment {
			if idx := strings.Index(line, "*/"); idx != -1 {
				line = line[idx+2:]
				inMultiLineComment = false
			} else {
				continue
			}
		}

		// Check for multi-line comment start
		if idx := strings.Index(line, "/*"); idx != -1 {
			endIdx := strings.Index(line[idx:], "*/")
			if endIdx != -1 {
				// Comment starts and ends on same line
				line = line[:idx] + line[idx+endIdx+2:]
			} else {
				// Comment continues to next line
				line = line[:idx]
				inMultiLineComment = true
			}
		}

		// Remove single-line comments
		if idx := strings.Index(line, "//"); idx != -1 {
			// Make sure it's not inside a string
			inString := false
			escaped := false
			for i := 0; i < idx; i++ {
				if line[i] == '\\' && !escaped {
					escaped = true
					continue
				}
				if line[i] == '"' && !escaped {
					inString = !inString
				}
				escaped = false
			}
			if !inString {
				line = line[:idx]
			}
		}

		result.WriteString(line)
		result.WriteString("\n")
	}

	return result.Bytes()
}

// LoadTaskDefinition loads a task definition from a reader
func LoadTaskDefinition(r io.Reader) (*TaskDefinition, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read input: %w", err)
	}

	// Remove JSONC comments
	cleanData := removeJSONComments(data)

	var taskDef TaskDefinition
	if err := json.Unmarshal(cleanData, &taskDef); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return &taskDef, nil
}

// LoadContainerDefinitions loads container definitions from a reader
func LoadContainerDefinitions(r io.Reader) ([]ContainerDefinition, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read input: %w", err)
	}

	// Remove JSONC comments
	cleanData := removeJSONComments(data)

	// First check if it's an array
	var containers []ContainerDefinition
	if err := json.Unmarshal(cleanData, &containers); err != nil {
		// Check if it's a single object (which should be an error)
		var singleContainer ContainerDefinition
		if err2 := json.Unmarshal(cleanData, &singleContainer); err2 == nil {
			return nil, fmt.Errorf("input must be an array of container definitions")
		}
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return containers, nil
}

// LoadFromFile loads data from a file based on mode
func LoadFromFile(filename string, mode LoadMode) (interface{}, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to close file: %v\n", err)
		}
	}()

	return Load(file, mode)
}

// Load loads data from a reader based on mode
func Load(r io.Reader, mode LoadMode) (interface{}, error) {
	switch mode {
	case ModeTask:
		return LoadTaskDefinition(r)
	case ModeContainer:
		return LoadContainerDefinitions(r)
	default:
		return nil, fmt.Errorf("invalid mode: %s", mode)
	}
}
