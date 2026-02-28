package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/engineervix/kwelea/internal/config"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the documentation site",
	Long:  `Build generates the complete static documentation site from your docs/ folder.`,
	RunE:  runBuild,
}

func runBuild(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return err
	}

	fmt.Printf("→ building site: %q\n", cfg.Site.Title)
	fmt.Printf("  docs:   %s\n", cfg.Build.DocsDir)
	fmt.Printf("  output: %s\n", cfg.Build.OutputDir)
	fmt.Println("  (build pipeline not yet implemented — Phase 4)")
	return nil
}
