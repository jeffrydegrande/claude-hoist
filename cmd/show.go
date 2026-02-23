package cmd

import (
	"fmt"
	"os"

	"github.com/jeffrydegrande/claude-hoist/hoist"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show project permissions that aren't in your user config yet",
	Run: func(cmd *cobra.Command, args []string) {
		_, _, _, newAllow, newDeny, err := hoist.LoadBoth()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

		if len(newAllow) == 0 && len(newDeny) == 0 {
			fmt.Println("nothing new — all project permissions already exist in user config")
			return
		}

		if len(newAllow) > 0 {
			fmt.Printf("New allow rules (%d):\n", len(newAllow))
			for _, rule := range newAllow {
				fmt.Printf("  + %s\n", rule)
			}
		}

		if len(newDeny) > 0 {
			if len(newAllow) > 0 {
				fmt.Println()
			}
			fmt.Printf("New deny rules (%d):\n", len(newDeny))
			for _, rule := range newDeny {
				fmt.Printf("  + %s\n", rule)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}
