package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "claude-hoist",
	Short: "Hoist Claude project permissions to your user config",
	Long: `claude-hoist reads .claude/settings.local.json from the current project
and merges its permission rules into your user-level ~/.claude/settings.local.json.

This lets you promote project-specific permission decisions to apply globally,
so you don't have to re-approve the same tools across projects.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
