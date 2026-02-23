package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/jeffrydegrande/claude-hoist/hoist"
	"github.com/spf13/cobra"
)

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Show a unified diff of what would change in your user config",
	Run: func(cmd *cobra.Command, args []string) {
		_, user, userPath, newAllow, newDeny, err := hoist.LoadBoth()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

		if len(newAllow) == 0 && len(newDeny) == 0 {
			fmt.Println("nothing to do — user config already has all project permissions")
			return
		}

		merged := hoist.Merge(user, newAllow, newDeny)

		before, _ := json.MarshalIndent(user, "", "  ")
		after, _ := json.MarshalIndent(merged, "", "  ")

		d := hoist.UnifiedDiff(userPath, userPath+" (merged)", string(before)+"\n", string(after)+"\n")
		fmt.Print(d)
	},
}

func init() {
	rootCmd.AddCommand(diffCmd)
}
