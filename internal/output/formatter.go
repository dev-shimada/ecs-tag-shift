package output

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/dev-shimada/ecs-tag-shift/internal/taskdef"
	"gopkg.in/yaml.v3"
)

// OutputFormat represents the output format
type OutputFormat string

const (
	FormatJSON OutputFormat = "json"
	FormatYAML OutputFormat = "yaml"
	FormatText OutputFormat = "text"
)

// FormatTaskDefinition formats a task definition for output
func FormatTaskDefinition(w io.Writer, taskDef *taskdef.TaskDefinition, format OutputFormat, showAll bool) error {
	switch format {
	case FormatJSON:
		return formatTaskDefinitionJSON(w, taskDef, showAll)
	case FormatYAML:
		return formatTaskDefinitionYAML(w, taskDef, showAll)
	case FormatText:
		return formatTaskDefinitionText(w, taskDef, showAll)
	default:
		return fmt.Errorf("unsupported output format: %s", format)
	}
}

// FormatContainerDefinitions formats container definitions for output
func FormatContainerDefinitions(w io.Writer, containers []taskdef.ContainerDefinition, format OutputFormat, showAll bool) error {
	switch format {
	case FormatJSON:
		return formatContainerDefinitionsJSON(w, containers, showAll)
	case FormatYAML:
		return formatContainerDefinitionsYAML(w, containers, showAll)
	case FormatText:
		return formatContainerDefinitionsText(w, containers, showAll)
	default:
		return fmt.Errorf("unsupported output format: %s", format)
	}
}

// formatTaskDefinitionJSON formats task definition as JSON
func formatTaskDefinitionJSON(w io.Writer, taskDef *taskdef.TaskDefinition, showAll bool) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")

	if showAll {
		return encoder.Encode(taskDef)
	}

	// Simplified output
	simplified := map[string]interface{}{
		"family":   taskDef.Family,
		"revision": taskDef.Revision,
		"containers": func() map[string]string {
			containers := make(map[string]string)
			for _, c := range taskDef.ContainerDefinitions {
				containers[c.Name] = c.Image
			}
			return containers
		}(),
	}
	return encoder.Encode(simplified)
}

// formatTaskDefinitionYAML formats task definition as YAML
func formatTaskDefinitionYAML(w io.Writer, taskDef *taskdef.TaskDefinition, showAll bool) error {
	encoder := yaml.NewEncoder(w)
	defer encoder.Close()

	if showAll {
		return encoder.Encode(taskDef)
	}

	// Simplified output
	simplified := map[string]interface{}{
		"family":   taskDef.Family,
		"revision": taskDef.Revision,
		"containers": func() map[string]string {
			containers := make(map[string]string)
			for _, c := range taskDef.ContainerDefinitions {
				containers[c.Name] = c.Image
			}
			return containers
		}(),
	}
	return encoder.Encode(simplified)
}

// formatTaskDefinitionText formats task definition as text
func formatTaskDefinitionText(w io.Writer, taskDef *taskdef.TaskDefinition, showAll bool) error {
	fmt.Fprintf(w, "Family: %s\n", taskDef.Family)
	if taskDef.Revision > 0 {
		fmt.Fprintf(w, "Revision: %d\n", taskDef.Revision)
	}
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Containers:")
	for _, c := range taskDef.ContainerDefinitions {
		fmt.Fprintf(w, "  - %s: %s\n", c.Name, c.Image)
	}
	return nil
}

// formatContainerDefinitionsJSON formats container definitions as JSON
func formatContainerDefinitionsJSON(w io.Writer, containers []taskdef.ContainerDefinition, showAll bool) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")

	if showAll {
		return encoder.Encode(containers)
	}

	// Simplified output
	simplified := map[string]interface{}{
		"containers": func() map[string]string {
			result := make(map[string]string)
			for _, c := range containers {
				result[c.Name] = c.Image
			}
			return result
		}(),
	}
	return encoder.Encode(simplified)
}

// formatContainerDefinitionsYAML formats container definitions as YAML
func formatContainerDefinitionsYAML(w io.Writer, containers []taskdef.ContainerDefinition, showAll bool) error {
	encoder := yaml.NewEncoder(w)
	defer encoder.Close()

	if showAll {
		return encoder.Encode(containers)
	}

	// Simplified output
	simplified := map[string]interface{}{
		"containers": func() map[string]string {
			result := make(map[string]string)
			for _, c := range containers {
				result[c.Name] = c.Image
			}
			return result
		}(),
	}
	return encoder.Encode(simplified)
}

// formatContainerDefinitionsText formats container definitions as text
func formatContainerDefinitionsText(w io.Writer, containers []taskdef.ContainerDefinition, showAll bool) error {
	fmt.Fprintln(w, "Containers:")
	for _, c := range containers {
		fmt.Fprintf(w, "  - %s: %s\n", c.Name, c.Image)
	}
	return nil
}

// FormatTaskDefinitionFull formats a full task definition (for shift command output)
func FormatTaskDefinitionFull(w io.Writer, taskDef *taskdef.TaskDefinition, format OutputFormat) error {
	switch format {
	case FormatJSON:
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		return encoder.Encode(taskDef)
	case FormatYAML:
		encoder := yaml.NewEncoder(w)
		defer encoder.Close()
		return encoder.Encode(taskDef)
	default:
		return fmt.Errorf("unsupported output format: %s", format)
	}
}

// FormatContainerDefinitionsFull formats full container definitions (for shift command output)
func FormatContainerDefinitionsFull(w io.Writer, containers []taskdef.ContainerDefinition, format OutputFormat) error {
	switch format {
	case FormatJSON:
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		return encoder.Encode(containers)
	case FormatYAML:
		encoder := yaml.NewEncoder(w)
		defer encoder.Close()
		return encoder.Encode(containers)
	default:
		return fmt.Errorf("unsupported output format: %s", format)
	}
}
