package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/engineervix/kwelea/internal/config"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the live-reloading development server",
	Long:  `Serve watches your docs/ folder and rebuilds on every save, with live browser reload.`,
	RunE:  runServe,
}

func runServe(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return err
	}

	fmt.Printf("→ starting dev server on http://localhost:%d\n", cfg.Serve.Port)
	fmt.Println("  (dev server not yet implemented — Phase 6)")
	return nil
}
