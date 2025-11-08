package taskdef

import (
	"testing"
)

func TestParseImage(t *testing.T) {
	tests := []struct {
		name           string
		image          string
		expectedRepo   string
		expectedTag    string
	}{
		{
			name:           "Simple image with tag",
			image:          "nginx:latest",
			expectedRepo:   "nginx",
			expectedTag:    "latest",
		},
		{
			name:           "ECR image with tag",
			image:          "123456789.dkr.ecr.us-east-1.amazonaws.com/my-app:v1.2.3",
			expectedRepo:   "123456789.dkr.ecr.us-east-1.amazonaws.com/my-app",
			expectedTag:    "v1.2.3",
		},
		{
			name:           "Image without tag",
			image:          "nginx",
			expectedRepo:   "nginx",
			expectedTag:    "",
		},
		{
			name:           "Image with registry and no tag",
			image:          "123456789.dkr.ecr.us-east-1.amazonaws.com/my-app",
			expectedRepo:   "123456789.dkr.ecr.us-east-1.amazonaws.com/my-app",
			expectedTag:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, tag := parseImage(tt.image)
			if repo != tt.expectedRepo {
				t.Errorf("parseImage() repo = %q, expected %q", repo, tt.expectedRepo)
			}
			if tag != tt.expectedTag {
				t.Errorf("parseImage() tag = %q, expected %q", tag, tt.expectedTag)
			}
		})
	}
}

func TestUpdateTaskDefinition(t *testing.T) {
	tests := []struct {
		name    string
		taskDef *TaskDefinition
		opts    UpdateOptions
		wantErr bool
		check   func(*TaskDefinition) bool
	}{
		{
			name: "Update all containers",
			taskDef: &TaskDefinition{
				Family: "app",
				ContainerDefinitions: []ContainerDefinition{
					{Name: "web", Image: "nginx:latest"},
					{Name: "api", Image: "api:v1.0"},
				},
			},
			opts: UpdateOptions{Tag: "v2.0"},
			check: func(td *TaskDefinition) bool {
				return td.ContainerDefinitions[0].Image == "nginx:v2.0" &&
					td.ContainerDefinitions[1].Image == "api:v2.0"
			},
		},
		{
			name: "Update specific container",
			taskDef: &TaskDefinition{
				Family: "app",
				ContainerDefinitions: []ContainerDefinition{
					{Name: "web", Image: "nginx:latest"},
					{Name: "api", Image: "api:v1.0"},
				},
			},
			opts: UpdateOptions{Tag: "v3.0", ContainerName: "web"},
			check: func(td *TaskDefinition) bool {
				return td.ContainerDefinitions[0].Image == "nginx:v3.0" &&
					td.ContainerDefinitions[1].Image == "api:v1.0"
			},
		},
		{
			name: "Update by image filter",
			taskDef: &TaskDefinition{
				Family: "app",
				ContainerDefinitions: []ContainerDefinition{
					{Name: "web", Image: "nginx:latest"},
					{Name: "app", Image: "my-app:v1.0"},
				},
			},
			opts: UpdateOptions{Tag: "stable", ImageName: "my-app"},
			check: func(td *TaskDefinition) bool {
				return td.ContainerDefinitions[0].Image == "nginx:latest" &&
					td.ContainerDefinitions[1].Image == "my-app:stable"
			},
		},
		{
			name: "Error: container not found",
			taskDef: &TaskDefinition{
				Family: "app",
				ContainerDefinitions: []ContainerDefinition{
					{Name: "web", Image: "nginx:latest"},
				},
			},
			opts:    UpdateOptions{Tag: "v2.0", ContainerName: "notfound"},
			wantErr: true,
		},
		{
			name: "Error: image not found",
			taskDef: &TaskDefinition{
				Family: "app",
				ContainerDefinitions: []ContainerDefinition{
					{Name: "web", Image: "nginx:latest"},
				},
			},
			opts:    UpdateOptions{Tag: "v2.0", ImageName: "notfound"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := UpdateTaskDefinition(tt.taskDef, tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateTaskDefinition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.check != nil {
				if !tt.check(tt.taskDef) {
					t.Errorf("UpdateTaskDefinition() result does not match expected")
				}
			}
		})
	}
}

func TestUpdateContainerDefinitions(t *testing.T) {
	tests := []struct {
		name       string
		containers []ContainerDefinition
		opts       UpdateOptions
		wantErr    bool
		check      func([]ContainerDefinition) bool
	}{
		{
			name: "Update all containers",
			containers: []ContainerDefinition{
				{Name: "web", Image: "nginx:latest"},
				{Name: "api", Image: "api:v1.0"},
			},
			opts: UpdateOptions{Tag: "v2.0"},
			check: func(cd []ContainerDefinition) bool {
				return cd[0].Image == "nginx:v2.0" &&
					cd[1].Image == "api:v2.0"
			},
		},
		{
			name: "Update specific container",
			containers: []ContainerDefinition{
				{Name: "web", Image: "nginx:latest"},
				{Name: "api", Image: "api:v1.0"},
			},
			opts: UpdateOptions{Tag: "v3.0", ContainerName: "api"},
			check: func(cd []ContainerDefinition) bool {
				return cd[0].Image == "nginx:latest" &&
					cd[1].Image == "api:v3.0"
			},
		},
		{
			name: "Error: container not found",
			containers: []ContainerDefinition{
				{Name: "web", Image: "nginx:latest"},
			},
			opts:    UpdateOptions{Tag: "v2.0", ContainerName: "notfound"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := UpdateContainerDefinitions(tt.containers, tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateContainerDefinitions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.check != nil {
				if !tt.check(result) {
					t.Errorf("UpdateContainerDefinitions() result does not match expected")
				}
			}
		})
	}
}
