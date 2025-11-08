package command

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/dev-shimada/ecs-tag-shift/internal/output"
	"github.com/dev-shimada/ecs-tag-shift/internal/taskdef"
)

func TestOverwriteOption(t *testing.T) {
	// Create a temporary task definition file
	tmpFile := filepath.Join(t.TempDir(), "task-def.json")
	originalContent := `{
  "family": "my-app",
  "containerDefinitions": [
    {"name": "web", "image": "nginx:latest"}
  ]
}`

	if err := os.WriteFile(tmpFile, []byte(originalContent), 0644); err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	tests := []struct {
		name          string
		inputFile     string
		overwrite     bool
		tag           string
		shouldUpdateFile bool
		format        string
	}{
		{
			name:          "Overwrite enabled should update file",
			inputFile:     tmpFile,
			overwrite:     true,
			tag:           "v1.0.0",
			shouldUpdateFile: true,
			format:        "json",
		},
		{
			name:          "Overwrite disabled should not update file",
			inputFile:     tmpFile,
			overwrite:     false,
			tag:           "v2.0.0",
			shouldUpdateFile: false,
			format:        "json",
		},
		{
			name:          "Overwrite with YAML format",
			inputFile:     tmpFile,
			overwrite:     true,
			tag:           "v3.0.0",
			shouldUpdateFile: true,
			format:        "yaml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset file content
			if err := os.WriteFile(tt.inputFile, []byte(originalContent), 0644); err != nil {
				t.Fatalf("Failed to reset temp file: %v", err)
			}

			// Save original stat
			originalStat, err := os.Stat(tt.inputFile)
			if err != nil {
				t.Fatalf("Failed to stat file: %v", err)
			}

			// Create ShiftOptions
			opts := &ShiftOptions{
				Mode:         taskdef.ModeTask,
				Tag:          tt.tag,
				OutputFormat: tt.format,
				Format:       output.OutputFormat(tt.format),
			}

			// Simulate shift operation
			data, err := taskdef.LoadFromFile(tt.inputFile, opts.Mode)
			if err != nil {
				t.Fatalf("Failed to load file: %v", err)
			}

			updateOpts := taskdef.UpdateOptions{
				Tag: tt.tag,
			}

			taskDef := data.(*taskdef.TaskDefinition)
			if err := taskdef.UpdateTaskDefinition(taskDef, updateOpts); err != nil {
				t.Fatalf("Failed to update: %v", err)
			}

			// Write result
			if tt.overwrite {
				buf := &bytes.Buffer{}
				if err := output.FormatTaskDefinitionFull(buf, taskDef, opts.Format); err != nil {
					t.Fatalf("Failed to format: %v", err)
				}
				if err := os.WriteFile(tt.inputFile, buf.Bytes(), 0644); err != nil {
					t.Fatalf("Failed to write file: %v", err)
				}
			}

			// Check if file was modified
			newStat, err := os.Stat(tt.inputFile)
			if err != nil {
				t.Fatalf("Failed to stat file: %v", err)
			}

			if tt.shouldUpdateFile {
				if originalStat.ModTime() == newStat.ModTime() && originalStat.Size() == newStat.Size() {
					t.Errorf("File should have been modified but wasn't")
				}
				// Verify content was updated
				content, err := os.ReadFile(tt.inputFile)
				if err != nil {
					t.Fatalf("Failed to read file: %v", err)
				}
				contentStr := string(content)
				if !bytes.Contains([]byte(contentStr), []byte(tt.tag)) {
					t.Errorf("File content should contain new tag %q", tt.tag)
				}
			} else {
				if originalStat.Size() != newStat.Size() {
					t.Errorf("File should not have been modified but was")
				}
			}
		})
	}
}

func TestOverwriteWithContainerMode(t *testing.T) {
	// Create a temporary container definitions file
	tmpFile := filepath.Join(t.TempDir(), "containers.json")
	originalContent := `[
  {"name": "web", "image": "nginx:latest"}
]`

	if err := os.WriteFile(tmpFile, []byte(originalContent), 0644); err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	opts := &ShiftOptions{
		Mode:         taskdef.ModeContainer,
		Tag:          "v1.0.0",
		OutputFormat: "json",
		Format:       output.FormatJSON,
	}

	// Load
	data, err := taskdef.LoadFromFile(tmpFile, opts.Mode)
	if err != nil {
		t.Fatalf("Failed to load file: %v", err)
	}

	containers := data.([]taskdef.ContainerDefinition)

	// Update
	updateOpts := taskdef.UpdateOptions{
		Tag: opts.Tag,
	}

	updated, err := taskdef.UpdateContainerDefinitions(containers, updateOpts)
	if err != nil {
		t.Fatalf("Failed to update: %v", err)
	}

	// Write back to file (overwrite)
	buf := &bytes.Buffer{}
	if err := output.FormatContainerDefinitionsFull(buf, updated, opts.Format); err != nil {
		t.Fatalf("Failed to format: %v", err)
	}
	if err := os.WriteFile(tmpFile, buf.Bytes(), 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	// Verify
	content, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}
	contentStr := string(content)
	if !bytes.Contains([]byte(contentStr), []byte("v1.0.0")) {
		t.Errorf("File should contain new tag v1.0.0")
	}
}
