package cmd

import (
	"io/fs"

	"github.com/spf13/cobra"
)

var (
	cfgFile string
	assets  fs.FS
)

var rootCmd = &cobra.Command{
	Use:   "kwelea",
	Short: "A fast documentation generator for Go projects",
	Long: `Kwelea weaves Markdown, templates, and assets into beautiful,
fast documentation sites for Go projects.

Install once globally, use across all your Go projects.
Zero runtime dependencies in consuming projects.`,
	SilenceErrors: true, // main.go prints the error; avoid double-printing
	SilenceUsage:  true, // don't print usage on runtime errors (config not found, etc.)
}

// Execute initialises the CLI with the embedded assets FS and runs the root command.
func Execute(embeddedFS fs.FS) error {
	assets = embeddedFS
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(
		&cfgFile, "config", "kwelea.toml",
		"path to kwelea.toml config file",
	)
	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(initCmd)
}
