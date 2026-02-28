package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/engineervix/kwelea/internal/config"
	"github.com/engineervix/kwelea/internal/nav"
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

	site, err := nav.NewSite(cfg)
	if err != nil {
		return err
	}

	fmt.Printf("→ building site: %q\n", site.Title)
	fmt.Printf("  docs:   %s\n", cfg.Build.DocsDir)
	fmt.Printf("  output: %s\n", cfg.Build.OutputDir)
	fmt.Printf("  pages:  %d\n\n", len(site.Pages))

	// Nav tree dump — Phase 2 verification output.
	// This will be replaced by the actual HTML render pass in Phase 4.
	fmt.Println("nav tree:")
	for _, section := range site.Nav {
		label := section.Label
		if label == "" {
			label = "(root)"
		}
		fmt.Printf("  [%s]\n", label)
		for _, item := range section.Items {
			fmt.Printf("    %-30s %s\n", item.Path, item.Title)
		}
	}

	fmt.Println("\npage order (prev → next):")
	for i, p := range site.Pages {
		prev, next := "—", "—"
		if p.Prev != nil {
			prev = p.Prev.Path
		}
		if p.Next != nil {
			next = p.Next.Path
		}
		fmt.Printf("  %d. %-30s prev:%-25s next:%s\n", i+1, p.Path, prev, next)
	}

	return nil
}
