package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("cannot determine working directory: %w", err)
	}

	cfgPath := filepath.Join(cwd, "kwelea.toml")
	if _, err := os.Stat(cfgPath); err == nil {
		fmt.Println("kwelea.toml already exists — nothing to do")
		return nil
	}

	projectName := inferProjectName(cwd)

	docsDir := filepath.Join(cwd, "docs")
	if err := os.MkdirAll(docsDir, 0o755); err != nil {
		return fmt.Errorf("cannot create docs/: %w", err)
	}

	tomlContent := fmt.Sprintf(`[site]
title    = %q
version  = "v0.1.0"
base_url = ""
repo     = ""

[build]
docs_dir   = "docs"
output_dir = "site"

[serve]
port         = 4000
open_browser = true

[theme]
light_code_theme = "github"
dark_code_theme  = "github-dark"
`, projectName)

	if err := os.WriteFile(cfgPath, []byte(tomlContent), 0o644); err != nil {
		return fmt.Errorf("cannot write kwelea.toml: %w", err)
	}

	indexPath := filepath.Join(docsDir, "index.md")
	if _, err := os.Stat(indexPath); errors.Is(err, os.ErrNotExist) {
		indexContent := fmt.Sprintf(`---
title: %s
description: Documentation for %s
---

Welcome to the **%s** documentation.

## Getting Started

Add your first section here, then run:

`+"```"+`bash
kwelea serve
`+"```"+`
`, projectName, projectName, projectName)
		if err := os.WriteFile(indexPath, []byte(indexContent), 0o644); err != nil {
			return fmt.Errorf("cannot write docs/index.md: %w", err)
		}
		fmt.Println("  created docs/index.md")
	}

	fmt.Printf("  created kwelea.toml  (title = %q)\n", projectName)
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  edit kwelea.toml     — set base_url and repo")
	fmt.Println("  kwelea serve         — start the dev server")
	fmt.Println("  kwelea build         — build to site/")
	return nil
}

// inferProjectName reads the module path from go.mod and returns the last
// path segment as a human-readable project name. Falls back to the directory
// name if go.mod is absent or unparseable.
func inferProjectName(dir string) string {
	data, err := os.ReadFile(filepath.Join(dir, "go.mod"))
	if err != nil {
		return filepath.Base(dir)
	}
	for line := range strings.SplitSeq(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			module := strings.TrimSpace(strings.TrimPrefix(line, "module "))
			if i := strings.LastIndex(module, "/"); i >= 0 {
				return module[i+1:]
			}
			return module
		}
	}
	return filepath.Base(dir)
}
