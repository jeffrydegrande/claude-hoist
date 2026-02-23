package cmd

import (
	"fmt"
	"os"

	"github.com/jeffrydegrande/claude-hoist/hoist"
	"github.com/spf13/cobra"
)

var stepCmd = &cobra.Command{
	Use:   "step",
	Short: "Step through each new permission one by one",
	Run: func(cmd *cobra.Command, args []string) {
		_, user, userPath, newAllow, newDeny, err := hoist.LoadBoth()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

		if len(newAllow) == 0 && len(newDeny) == 0 {
			fmt.Println("nothing new — all project permissions already exist in user config")
			return
		}

		var acceptedAllow, acceptedDeny []string

		if len(newAllow) > 0 {
			fmt.Printf("Allow rules (%d new):\n\n", len(newAllow))
			accepted := stepThrough(newAllow)
			acceptedAllow = append(acceptedAllow, accepted...)
		}

		if len(newDeny) > 0 {
			fmt.Printf("\nDeny rules (%d new):\n\n", len(newDeny))
			accepted := stepThrough(newDeny)
			acceptedDeny = append(acceptedDeny, accepted...)
		}

		if len(acceptedAllow) == 0 && len(acceptedDeny) == 0 {
			fmt.Println("\nnothing selected")
			return
		}

		merged := hoist.Merge(user, acceptedAllow, acceptedDeny)
		if err := hoist.WriteSettings(userPath, merged); err != nil {
			fmt.Fprintf(os.Stderr, "error writing: %v\n", err)
			os.Exit(1)
		}

		total := len(acceptedAllow) + len(acceptedDeny)
		fmt.Printf("\ndone — added %d rule(s) to %s\n", total, userPath)
	},
}

// stepThrough prompts for each rule. Returns accepted ones.
// y = accept, n = skip, q = quit (skip remaining).
func stepThrough(rules []string) []string {
	var accepted []string
	for i, rule := range rules {
		fmt.Printf("  [%d/%d] %s\n", i+1, len(rules), rule)
		fmt.Print("  add? [y/n/q] ")

		var answer string
		fmt.Scanln(&answer)

		switch answer {
		case "y", "Y":
			accepted = append(accepted, rule)
		case "q", "Q":
			fmt.Println("  skipping remaining")
			return accepted
		default:
			// skip
		}
	}
	return accepted
}

func init() {
	rootCmd.AddCommand(stepCmd)
}
