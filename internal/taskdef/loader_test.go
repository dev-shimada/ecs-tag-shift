package taskdef

import (
	"strings"
	"testing"
)

func TestRemoveJSONComments(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		shouldContain string
		shouldNotContain string
	}{
		{
			name:          "No comments",
			input:         `{"key": "value"}`,
			shouldContain: "key",
			shouldNotContain: "",
		},
		{
			name:          "Single line comment",
			input:         `{"key": "value"} // This is a comment`,
			shouldContain: "key",
			shouldNotContain: "This is a",
		},
		{
			name:          "Multi-line comment",
			input:         `{"key": "value"} /* This is a multi-line comment */ {"key2": "value2"}`,
			shouldContain: "key",
			shouldNotContain: "This is",
		},
		{
			name:          "Comment inside string should not be removed",
			input:         `{"key": "value // not a comment"}`,
			shouldContain: "value // not a comment",
			shouldNotContain: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := string(removeJSONComments([]byte(tt.input)))
			if !strings.Contains(result, tt.shouldContain) {
				t.Errorf("removeJSONComments() should contain %q, got %q", tt.shouldContain, result)
			}
			if tt.shouldNotContain != "" && strings.Contains(result, tt.shouldNotContain) {
				t.Errorf("removeJSONComments() should not contain %q, got %q", tt.shouldNotContain, result)
			}
		})
	}
}

func TestLoadTaskDefinition(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		check   func(*TaskDefinition) bool
	}{
		{
			name:    "Valid task definition",
			input:   `{"family": "my-app", "containerDefinitions": []}`,
			wantErr: false,
			check: func(td *TaskDefinition) bool {
				return td.Family == "my-app" && len(td.ContainerDefinitions) == 0
			},
		},
		{
			name:    "Task definition with containers",
			input:   `{"family": "app", "containerDefinitions": [{"name": "web", "image": "nginx:latest"}]}`,
			wantErr: false,
			check: func(td *TaskDefinition) bool {
				return td.Family == "app" && len(td.ContainerDefinitions) == 1 && td.ContainerDefinitions[0].Name == "web"
			},
		},
		{
			name:    "Task definition with JSONC comments",
			input:   "{\"family\": \"app\", // my app\n\"containerDefinitions\": []}",
			wantErr: false,
			check: func(td *TaskDefinition) bool {
				return td.Family == "app"
			},
		},
		{
			name:    "Invalid JSON",
			input:   `{invalid json}`,
			wantErr: true,
			check:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := LoadTaskDefinition(strings.NewReader(tt.input))
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadTaskDefinition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.check != nil {
				if !tt.check(result) {
					t.Errorf("LoadTaskDefinition() result does not match expected")
				}
			}
		})
	}
}

func TestLoadContainerDefinitions(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		check   func([]ContainerDefinition) bool
	}{
		{
			name:    "Valid container definitions array",
			input:   `[{"name": "web", "image": "nginx:latest"}]`,
			wantErr: false,
			check: func(cd []ContainerDefinition) bool {
				return len(cd) == 1 && cd[0].Name == "web"
			},
		},
		{
			name:    "Empty container definitions array",
			input:   `[]`,
			wantErr: false,
			check: func(cd []ContainerDefinition) bool {
				return len(cd) == 0
			},
		},
		{
			name:    "Single object should fail",
			input:   `{"name": "web", "image": "nginx:latest"}`,
			wantErr: true,
			check:   nil,
		},
		{
			name:    "Container definitions with JSONC comments",
			input:   "[\n  // web container\n  {\"name\": \"web\", \"image\": \"nginx:latest\"}\n]",
			wantErr: false,
			check: func(cd []ContainerDefinition) bool {
				return len(cd) == 1
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := LoadContainerDefinitions(strings.NewReader(tt.input))
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadContainerDefinitions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.check != nil {
				if !tt.check(result) {
					t.Errorf("LoadContainerDefinitions() result does not match expected")
				}
			}
		})
	}
}
