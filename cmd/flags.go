package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/engineervix/kwelea/internal/config"
)

// applyFlagOverrides applies CLI flag overrides to cfg. It is called from
// every subcommand that supports flag overrides (build, serve), so the
// override logic stays in one place — no risk of drift between commands.
//
// Empty values are rejected with a clear error pointing at the offending
// flag. Without this guard:
//   - `kwelea build --base-url ""` silently breaks every URL-emitting
//     feature (og:url, sitemap, og:image, etc.)
//   - `kwelea build --output ""` produces a cryptic OS error from
//     os.MkdirAll("", ...)
//   - `kwelea build --source ""` produces a cryptic OS error from
//     filepath.WalkDir("", ...)
//
// Flag values are read directly from cmd.Flags() so this function does not
// depend on package-level variables. This makes the override behaviour
// reproducible across multiple Execute() calls in the same process.
func applyFlagOverrides(cmd *cobra.Command, cfg *config.Config) error {
	if cmd.Flags().Changed("source") {
		src, _ := cmd.Flags().GetString("source")
		if strings.TrimSpace(src) == "" {
			return fmt.Errorf("--source: value must not be empty")
		}
		cfg.Build.DocsDir = src
	}

	if cmd.Flags().Changed("output") {
		out, _ := cmd.Flags().GetString("output")
		if strings.TrimSpace(out) == "" {
			return fmt.Errorf("--output: value must not be empty")
		}
		cfg.Build.OutputDir = out
	}

	if cmd.Flags().Changed("base-url") {
		url, _ := cmd.Flags().GetString("base-url")
		if strings.TrimSpace(url) == "" {
			return fmt.Errorf("--base-url: value must not be empty")
		}
		cfg.Site.BaseURL = url
	}

	return nil
}
