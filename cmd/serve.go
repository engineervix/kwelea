package cmd

import (
	"github.com/spf13/cobra"

	"github.com/engineervix/kwelea/internal/config"
	"github.com/engineervix/kwelea/internal/server"
)

var serveSource string

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the live-reloading development server",
	Long:  `Serve watches your docs/ folder and rebuilds on every save, with live browser reload.`,
	RunE:  runServe,
}

func init() {
	serveCmd.Flags().StringVar(&serveSource, "source", "", "override build.docs_dir from kwelea.toml")
}

func runServe(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return err
	}

	if cmd.Flags().Changed("source") {
		cfg.Build.DocsDir = serveSource
	}

	return server.Start(cfg, assets, cfgFile)
}