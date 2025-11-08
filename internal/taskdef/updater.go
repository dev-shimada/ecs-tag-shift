package taskdef

import (
	"fmt"
	"strings"
)

// UpdateOptions represents options for updating container image tags
type UpdateOptions struct {
	Tag           string
	ContainerName string
	ImageName     string
}

// parseImage splits an image string into repository and tag
// e.g., "nginx:latest" -> ("nginx", "latest")
// e.g., "123456789.dkr.ecr.us-east-1.amazonaws.com/my-app:v1.2.3" -> ("123456789.dkr.ecr.us-east-1.amazonaws.com/my-app", "v1.2.3")
func parseImage(image string) (repository string, tag string) {
	parts := strings.Split(image, ":")
	if len(parts) == 1 {
		return image, ""
	}

	// Handle case where there are multiple colons (e.g., URL with port)
	// We assume the last colon separates the tag
	repository = strings.Join(parts[:len(parts)-1], ":")
	tag = parts[len(parts)-1]

	return repository, tag
}

// matchesFilter checks if a container matches the filter criteria
func matchesFilter(container *ContainerDefinition, opts UpdateOptions) bool {
	// Filter by container name
	if opts.ContainerName != "" && container.Name != opts.ContainerName {
		return false
	}

	// Filter by image repository name
	if opts.ImageName != "" {
		repository, _ := parseImage(container.Image)
		// Extract just the repository name (without registry URL)
		repoName := repository
		if strings.Contains(repository, "/") {
			parts := strings.Split(repository, "/")
			repoName = parts[len(parts)-1]
		}

		if repoName != opts.ImageName && repository != opts.ImageName {
			return false
		}
	}

	return true
}

// updateContainerImage updates the image tag for a single container
func updateContainerImage(container *ContainerDefinition, newTag string) {
	repository, _ := parseImage(container.Image)
	container.Image = fmt.Sprintf("%s:%s", repository, newTag)
}

// UpdateTaskDefinition updates the container image tags in a task definition
func UpdateTaskDefinition(taskDef *TaskDefinition, opts UpdateOptions) error {
	updated := false

	for i := range taskDef.ContainerDefinitions {
		container := &taskDef.ContainerDefinitions[i]
		if matchesFilter(container, opts) {
			updateContainerImage(container, opts.Tag)
			updated = true
		}
	}

	if !updated && (opts.ContainerName != "" || opts.ImageName != "") {
		// User specified a filter but no containers matched
		if opts.ContainerName != "" {
			return fmt.Errorf("container '%s' not found in definitions", opts.ContainerName)
		}
		if opts.ImageName != "" {
			return fmt.Errorf("image '%s' not found in definitions", opts.ImageName)
		}
	}

	return nil
}

// UpdateContainerDefinitions updates the container image tags in a list of container definitions
func UpdateContainerDefinitions(containers []ContainerDefinition, opts UpdateOptions) ([]ContainerDefinition, error) {
	updated := false

	for i := range containers {
		container := &containers[i]
		if matchesFilter(container, opts) {
			updateContainerImage(container, opts.Tag)
			updated = true
		}
	}

	if !updated && (opts.ContainerName != "" || opts.ImageName != "") {
		// User specified a filter but no containers matched
		if opts.ContainerName != "" {
			return nil, fmt.Errorf("container '%s' not found in definitions", opts.ContainerName)
		}
		if opts.ImageName != "" {
			return nil, fmt.Errorf("image '%s' not found in definitions", opts.ImageName)
		}
	}

	return containers, nil
}
