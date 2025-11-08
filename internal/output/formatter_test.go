package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/dev-shimada/ecs-tag-shift/internal/taskdef"
	"gopkg.in/yaml.v3"
)

func TestFormatTaskDefinitionJSON(t *testing.T) {
	td := &taskdef.TaskDefinition{
		Family:   "my-app",
		Revision: 15,
		ContainerDefinitions: []taskdef.ContainerDefinition{
			{Name: "web", Image: "nginx:latest"},
		},
	}

	buf := &bytes.Buffer{}
	err := formatTaskDefinitionJSON(buf, td, false)
	if err != nil {
		t.Errorf("formatTaskDefinitionJSON() error = %v", err)
		return
	}

	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Errorf("Result is not valid JSON: %v", err)
		return
	}

	if result["family"] != "my-app" {
		t.Errorf("formatTaskDefinitionJSON() family = %v, expected my-app", result["family"])
	}
}

func TestFormatTaskDefinitionYAML(t *testing.T) {
	td := &taskdef.TaskDefinition{
		Family:   "my-app",
		Revision: 15,
		ContainerDefinitions: []taskdef.ContainerDefinition{
			{Name: "web", Image: "nginx:latest"},
		},
	}

	buf := &bytes.Buffer{}
	err := formatTaskDefinitionYAML(buf, td, false)
	if err != nil {
		t.Errorf("formatTaskDefinitionYAML() error = %v", err)
		return
	}

	var result map[string]interface{}
	if err := yaml.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Errorf("Result is not valid YAML: %v", err)
		return
	}

	if result["family"] != "my-app" {
		t.Errorf("formatTaskDefinitionYAML() family = %v, expected my-app", result["family"])
	}
}

func TestFormatTaskDefinitionText(t *testing.T) {
	td := &taskdef.TaskDefinition{
		Family:   "my-app",
		Revision: 15,
		ContainerDefinitions: []taskdef.ContainerDefinition{
			{Name: "web", Image: "nginx:latest"},
		},
	}

	buf := &bytes.Buffer{}
	err := formatTaskDefinitionText(buf, td, false)
	if err != nil {
		t.Errorf("formatTaskDefinitionText() error = %v", err)
		return
	}

	output := buf.String()
	if !strings.Contains(output, "Family: my-app") {
		t.Errorf("formatTaskDefinitionText() output does not contain 'Family: my-app'")
	}
	if !strings.Contains(output, "web: nginx:latest") {
		t.Errorf("formatTaskDefinitionText() output does not contain 'web: nginx:latest'")
	}
}

func TestFormatContainerDefinitionsJSON(t *testing.T) {
	containers := []taskdef.ContainerDefinition{
		{Name: "web", Image: "nginx:latest"},
		{Name: "api", Image: "api:v1.0"},
	}

	buf := &bytes.Buffer{}
	err := formatContainerDefinitionsJSON(buf, containers, false)
	if err != nil {
		t.Errorf("formatContainerDefinitionsJSON() error = %v", err)
		return
	}

	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Errorf("Result is not valid JSON: %v", err)
		return
	}

	containerMap, ok := result["containers"].(map[string]interface{})
	if !ok {
		t.Errorf("formatContainerDefinitionsJSON() containers field is not a map")
		return
	}

	if containerMap["web"] != "nginx:latest" {
		t.Errorf("formatContainerDefinitionsJSON() web image = %v, expected nginx:latest", containerMap["web"])
	}
}

func TestFormatTaskDefinitionFull(t *testing.T) {
	td := &taskdef.TaskDefinition{
		Family:   "my-app",
		Revision: 15,
		ContainerDefinitions: []taskdef.ContainerDefinition{
			{
				Name:  "web",
				Image: "nginx:latest",
				CPU:   256,
			},
		},
	}

	tests := []struct {
		name   string
		format OutputFormat
	}{
		{
			name:   "JSON format",
			format: FormatJSON,
		},
		{
			name:   "YAML format",
			format: FormatYAML,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			err := FormatTaskDefinitionFull(buf, td, tt.format)
			if err != nil {
				t.Errorf("FormatTaskDefinitionFull() error = %v", err)
				return
			}

			if buf.Len() == 0 {
				t.Errorf("FormatTaskDefinitionFull() produced no output")
			}
		})
	}
}

func TestFormatContainerDefinitionsFull(t *testing.T) {
	containers := []taskdef.ContainerDefinition{
		{Name: "web", Image: "nginx:latest", CPU: 256},
		{Name: "api", Image: "api:v1.0", CPU: 512},
	}

	tests := []struct {
		name   string
		format OutputFormat
	}{
		{
			name:   "JSON format",
			format: FormatJSON,
		},
		{
			name:   "YAML format",
			format: FormatYAML,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			err := FormatContainerDefinitionsFull(buf, containers, tt.format)
			if err != nil {
				t.Errorf("FormatContainerDefinitionsFull() error = %v", err)
				return
			}

			if buf.Len() == 0 {
				t.Errorf("FormatContainerDefinitionsFull() produced no output")
			}
		})
	}
}
