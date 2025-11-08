package taskdef

// ContainerDefinition represents an ECS container definition
type ContainerDefinition struct {
	Name         string                 `json:"name" yaml:"name"`
	Image        string                 `json:"image" yaml:"image"`
	CPU          int                    `json:"cpu,omitempty" yaml:"cpu,omitempty"`
	Memory       int                    `json:"memory,omitempty" yaml:"memory,omitempty"`
	Essential    *bool                  `json:"essential,omitempty" yaml:"essential,omitempty"`
	PortMappings []PortMapping          `json:"portMappings,omitempty" yaml:"portMappings,omitempty"`
	Environment  []EnvironmentVariable  `json:"environment,omitempty" yaml:"environment,omitempty"`
	// Store all other fields as-is
	Extra map[string]interface{} `json:"-" yaml:"-"`
}

// PortMapping represents a port mapping configuration
type PortMapping struct {
	ContainerPort int    `json:"containerPort,omitempty" yaml:"containerPort,omitempty"`
	HostPort      int    `json:"hostPort,omitempty" yaml:"hostPort,omitempty"`
	Protocol      string `json:"protocol,omitempty" yaml:"protocol,omitempty"`
}

// EnvironmentVariable represents an environment variable
type EnvironmentVariable struct {
	Name  string `json:"name" yaml:"name"`
	Value string `json:"value" yaml:"value"`
}

// TaskDefinition represents an ECS task definition
type TaskDefinition struct {
	Family                  string                 `json:"family,omitempty" yaml:"family,omitempty"`
	TaskRoleArn             string                 `json:"taskRoleArn,omitempty" yaml:"taskRoleArn,omitempty"`
	ExecutionRoleArn        string                 `json:"executionRoleArn,omitempty" yaml:"executionRoleArn,omitempty"`
	NetworkMode             string                 `json:"networkMode,omitempty" yaml:"networkMode,omitempty"`
	ContainerDefinitions    []ContainerDefinition  `json:"containerDefinitions" yaml:"containerDefinitions"`
	RequiresCompatibilities []string               `json:"requiresCompatibilities,omitempty" yaml:"requiresCompatibilities,omitempty"`
	CPU                     string                 `json:"cpu,omitempty" yaml:"cpu,omitempty"`
	Memory                  string                 `json:"memory,omitempty" yaml:"memory,omitempty"`
	Revision                int                    `json:"revision,omitempty" yaml:"revision,omitempty"`
	// Store all other fields as-is
	Extra map[string]interface{} `json:"-" yaml:"-"`
}
