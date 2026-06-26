package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/engineervix/kwelea/internal/builder"
	"github.com/engineervix/kwelea/internal/config"
	"github.com/engineervix/kwelea/internal/nav"
)

var (
	buildBaseURL string
	buildOutput  string
	buildSource  string
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the documentation site",
	Long:  `Build generates the complete static documentation site from your docs/ folder.`,
	RunE:  runBuild,
}

func init() {
	buildCmd.Flags().StringVar(&buildBaseURL, "base-url", "", "override site.base_url from kwelea.toml (e.g. https://kwelea.pages.dev)")
	buildCmd.Flags().StringVar(&buildOutput, "output", "", "override build.output_dir from kwelea.toml")
	buildCmd.Flags().StringVar(&buildSource, "source", "", "override build.docs_dir from kwelea.toml")
}

func runBuild(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return err
	}

	// Apply CLI flag overrides. Changed() is used (not empty-string check)
	// so that legitimate empty values like `--base-url ""` are honoured.
	if cmd.Flags().Changed("base-url") {
		cfg.Site.BaseURL = buildBaseURL
	}
	if cmd.Flags().Changed("output") {
		cfg.Build.OutputDir = buildOutput
	}
	if cmd.Flags().Changed("source") {
		cfg.Build.DocsDir = buildSource
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