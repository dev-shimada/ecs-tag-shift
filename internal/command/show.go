package command

import (
	"fmt"
	"os"

	"github.com/dev-shimada/ecs-tag-shift/internal/output"
	"github.com/dev-shimada/ecs-tag-shift/internal/taskdef"
	"github.com/spf13/cobra"
)

// ShowOptions represents options for the show command
type ShowOptions struct {
	Mode       taskdef.LoadMode
	OutputFile string
	Format     output.OutputFormat
	ShowAll    bool
}

// NewShowCommand creates a new show command
func NewShowCommand(globalMode *taskdef.LoadMode) *cobra.Command {
	opts := &ShowOptions{}

	cmd := &cobra.Command{
		Use:   "show [file]",
		Short: "Display task definition or container definitions",
		Long:  `Display the contents of a task definition or container definitions file.`,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Mode = *globalMode
			return runShow(args, opts)
		},
	}

	cmd.Flags().StringVarP(&opts.OutputFile, "output", "o", "json", "Output format (json, yaml, text)")
	cmd.Flags().BoolVar(&opts.ShowAll, "all", false, "Show all fields")

	return cmd
}

func runShow(args []string, opts *ShowOptions) error {
	// Parse output format
	opts.Format = output.OutputFormat(opts.OutputFile)
	if opts.Format != output.FormatJSON && opts.Format != output.FormatYAML && opts.Format != output.FormatText {
		return fmt.Errorf("invalid output format: %s (must be json, yaml, or text)", opts.OutputFile)
	}

	// Load input
	var data interface{}
	var err error

	if len(args) > 0 {
		// Load from file
		data, err = taskdef.LoadFromFile(args[0], opts.Mode)
	} else {
		// Load from stdin
		data, err = taskdef.Load(os.Stdin, opts.Mode)
	}

	if err != nil {
		return err
	}

	// Format and output
	switch opts.Mode {
	case taskdef.ModeTask:
		taskDef := data.(*taskdef.TaskDefinition)
		return output.FormatTaskDefinition(os.Stdout, taskDef, opts.Format, opts.ShowAll)
	case taskdef.ModeContainer:
		containers := data.([]taskdef.ContainerDefinition)
		return output.FormatContainerDefinitions(os.Stdout, containers, opts.Format, opts.ShowAll)
	default:
		return fmt.Errorf("invalid mode: %s", opts.Mode)
	}
}
