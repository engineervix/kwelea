package cmd

import (
	"github.com/spf13/cobra"

	"github.com/engineervix/kwelea/internal/config"
)

// applyFlagOverrides applies CLI flag overrides to cfg. It is called from
// every subcommand that supports flag overrides (build, serve), so the
// override logic stays in one place — no risk of drift between commands.
//
// Flag values are read directly from cmd.Flags() so this function does not
// depend on package-level variables. This makes the override behaviour
// reproducible across multiple Execute() calls in the same process.
func applyFlagOverrides(cmd *cobra.Command, cfg *config.Config) error {
	if cmd.Flags().Changed("source") {
		src, _ := cmd.Flags().GetString("source")
		cfg.Build.DocsDir = src
	}

	if cmd.Flags().Changed("output") {
		out, _ := cmd.Flags().GetString("output")
		cfg.Build.OutputDir = out
	}

	if cmd.Flags().Changed("base-url") {
		url, _ := cmd.Flags().GetString("base-url")
		cfg.Site.BaseURL = url
	}

	return nil
}
