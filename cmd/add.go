package cmd

import (
	"fmt"
	"os"

	"github.com/jeffrydegrande/claude-hoist/hoist"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add all project permissions to your user config",
	Run: func(cmd *cobra.Command, args []string) {
		_, user, userPath, newAllow, newDeny, err := hoist.LoadBoth()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

		if len(newAllow) == 0 && len(newDeny) == 0 {
			fmt.Println("nothing to do — all project permissions already exist in user config")
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

		yes, _ := cmd.Flags().GetBool("yes")
		if !yes {
			fmt.Printf("\nMerge into %s? [y/N] ", userPath)
			var answer string
			fmt.Scanln(&answer)
			if answer != "y" && answer != "Y" {
				fmt.Println("aborted")
				return
			}
		}

		merged := hoist.Merge(user, newAllow, newDeny)
		if err := hoist.WriteSettings(userPath, merged); err != nil {
			fmt.Fprintf(os.Stderr, "error writing: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("done — wrote %s\n", userPath)
	},
}

func init() {
	addCmd.Flags().BoolP("yes", "y", false, "skip confirmation prompt")
	rootCmd.AddCommand(addCmd)
}
