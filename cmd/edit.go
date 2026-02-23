package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/jeffrydegrande/claude-hoist/hoist"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit [project|user]",
	Short: "Open project or user settings in $EDITOR",
	Long: `Opens the Claude settings.local.json file in your $EDITOR.

  claude-hoist edit project   opens .claude/settings.local.json in the current directory
  claude-hoist edit user      opens ~/.claude/settings.local.json
  claude-hoist edit           defaults to project`,
	Args:      cobra.MaximumNArgs(1),
	ValidArgs: []string{"project", "user"},
	Run: func(cmd *cobra.Command, args []string) {
		target := "project"
		if len(args) > 0 {
			target = args[0]
		}

		var path string
		var err error

		switch target {
		case "project":
			path, err = hoist.FindProjectSettings()
		case "user":
			path, err = hoist.UserSettingsPath()
		default:
			fmt.Fprintf(os.Stderr, "unknown target %q — use 'project' or 'user'\n", target)
			os.Exit(1)
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "vi"
		}

		c := exec.Command(editor, path)
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr

		if err := c.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "editor exited with error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}
