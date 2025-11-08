package command

import (
	"fmt"
	"os"

	"github.com/dev-shimada/ecs-tag-shift/internal/output"
	"github.com/dev-shimada/ecs-tag-shift/internal/taskdef"
	"github.com/spf13/cobra"
)

// ShiftOptions represents options for the shift command
type ShiftOptions struct {
	Mode          taskdef.LoadMode
	Tag           string
	ContainerName string
	ImageName     string
	OutputFormat  string
	Format        output.OutputFormat
	Overwrite     bool
}

// NewShiftCommand creates a new shift command
func NewShiftCommand(globalMode *taskdef.LoadMode) *cobra.Command {
	opts := &ShiftOptions{}

	cmd := &cobra.Command{
		Use:   "shift [file]",
		Short: "Update container image tags",
		Long:  `Update the image tags for containers in a task definition or container definitions file.`,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Mode = *globalMode
			return runShift(args, opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Tag, "tag", "t", "", "New image tag (required)")
	cmd.Flags().StringVarP(&opts.ContainerName, "container", "c", "", "Filter by container name")
	cmd.Flags().StringVarP(&opts.ImageName, "image", "i", "", "Filter by image repository name")
	cmd.Flags().StringVarP(&opts.OutputFormat, "output", "o", "json", "Output format (json, yaml)")
	cmd.Flags().BoolVarP(&opts.Overwrite, "overwrite", "w", false, "Overwrite input file (only with file input)")

	cmd.MarkFlagRequired("tag")

	return cmd
}

func runShift(args []string, opts *ShiftOptions) error {
	// Validate tag
	if opts.Tag == "" {
		return fmt.Errorf("tag is required")
	}

	// Parse output format
	opts.Format = output.OutputFormat(opts.OutputFormat)
	if opts.Format != output.FormatJSON && opts.Format != output.FormatYAML {
		return fmt.Errorf("invalid output format: %s (must be json or yaml)", opts.OutputFormat)
	}

	// Validate overwrite option
	if opts.Overwrite && len(args) == 0 {
		// Overwrite without file input - just ignore and output to stdout
		opts.Overwrite = false
	}

	// Load input
	var data interface{}
	var err error
	var inputFile string

	if len(args) > 0 {
		// Load from file
		inputFile = args[0]
		data, err = taskdef.LoadFromFile(inputFile, opts.Mode)
	} else {
		// Load from stdin
		data, err = taskdef.Load(os.Stdin, opts.Mode)
	}

	if err != nil {
		return err
	}

	// Create update options
	updateOpts := taskdef.UpdateOptions{
		Tag:           opts.Tag,
		ContainerName: opts.ContainerName,
		ImageName:     opts.ImageName,
	}

	// Update and output
	switch opts.Mode {
	case taskdef.ModeTask:
		taskDef := data.(*taskdef.TaskDefinition)
		if err := taskdef.UpdateTaskDefinition(taskDef, updateOpts); err != nil {
			return err
		}

		// Determine output destination
		if opts.Overwrite && inputFile != "" {
			return writeToFile(inputFile, taskDef, opts.Format, false)
		}
		return output.FormatTaskDefinitionFull(os.Stdout, taskDef, opts.Format)

	case taskdef.ModeContainer:
		containers := data.([]taskdef.ContainerDefinition)
		updatedContainers, err := taskdef.UpdateContainerDefinitions(containers, updateOpts)
		if err != nil {
			return err
		}

		// Determine output destination
		if opts.Overwrite && inputFile != "" {
			return writeToFile(inputFile, updatedContainers, opts.Format, true)
		}
		return output.FormatContainerDefinitionsFull(os.Stdout, updatedContainers, opts.Format)

	default:
		return fmt.Errorf("invalid mode: %s", opts.Mode)
	}
}

// writeToFile writes the result to a file
func writeToFile(filename string, data interface{}, format output.OutputFormat, isContainers bool) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to open file for writing: %w", err)
	}
	defer file.Close()

	if isContainers {
		containers := data.([]taskdef.ContainerDefinition)
		return output.FormatContainerDefinitionsFull(file, containers, format)
	}
	taskDef := data.(*taskdef.TaskDefinition)
	return output.FormatTaskDefinitionFull(file, taskDef, format)
}
