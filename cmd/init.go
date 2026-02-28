package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Scaffold a kwelea.toml and docs/ folder in the current project",
	Long: `Init creates a kwelea.toml and a docs/ directory with a starter index.md
in your Go project. Run it once at the root of the project you want to document.`,
	RunE: runInit,
}

func runInit(cmd *cobra.Command, args []string) error {
	fmt.Println("→ kwelea init (not yet implemented — Phase 7)")
	return nil
}
