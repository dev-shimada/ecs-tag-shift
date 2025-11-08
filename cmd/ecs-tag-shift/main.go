package main

import (
	"fmt"
	"os"

	"github.com/dev-shimada/ecs-tag-shift/internal/command"
	"github.com/dev-shimada/ecs-tag-shift/internal/taskdef"
	"github.com/spf13/cobra"
)

var version = "dev"

func main() {
	if err := newRootCommand().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func newRootCommand() *cobra.Command {
	var mode string
	globalMode := taskdef.ModeTask

	rootCmd := &cobra.Command{
		Use:   "ecs-tag-shift",
		Short: "A CLI tool to update ECS task definition and container definition image tags",
		Long: `ecs-tag-shift is a CLI tool for updating container image tags in
ECS task definitions and container definitions. It supports JSONC input
and can output in JSON, YAML, or text formats.`,
		Version: version,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Validate and set global mode
			switch mode {
			case "task":
				globalMode = taskdef.ModeTask
			case "container":
				globalMode = taskdef.ModeContainer
			default:
				return fmt.Errorf("invalid mode: %s (must be 'task' or 'container')", mode)
			}
			return nil
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// Global flags
	rootCmd.PersistentFlags().StringVarP(&mode, "mode", "m", "task", "Input mode (task or container)")

	// Add subcommands
	rootCmd.AddCommand(command.NewShowCommand(&globalMode))
	rootCmd.AddCommand(command.NewShiftCommand(&globalMode))

	return rootCmd
}
