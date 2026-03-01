package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/engineervix/kwelea/internal/builder"
	"github.com/engineervix/kwelea/internal/config"
	"github.com/engineervix/kwelea/internal/nav"
)

var buildBaseURL string

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the documentation site",
	Long:  `Build generates the complete static documentation site from your docs/ folder.`,
	RunE:  runBuild,
}

func init() {
	buildCmd.Flags().StringVar(&buildBaseURL, "base-url", "", "override base_url from kwelea.toml (e.g. https://kwelea.pages.dev)")
}

func runBuild(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return err
	}

	if buildBaseURL != "" {
		cfg.Site.BaseURL = buildBaseURL
	}

	site, err := nav.NewSite(cfg)
	if err != nil {
		return err
	}

	fmt.Printf("→ building %q  (%d pages)\n", site.Title, len(site.Pages))

	if err := builder.Build(site, assets, false); err != nil {
		return err
	}

	fmt.Printf("✓ site written to %q\n", cfg.Build.OutputDir)
	return nil
}
