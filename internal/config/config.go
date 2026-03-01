// Package config loads and validates kwelea.toml configuration files.
package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

// Config is the top-level configuration loaded from kwelea.toml.
type Config struct {
	Site  SiteConfig  `toml:"site"`
	Build BuildConfig `toml:"build"`
	Serve ServeConfig `toml:"serve"`
	Theme ThemeConfig `toml:"theme"`
	Nav   []NavEntry  `toml:"nav"`
}

// SiteConfig holds site-level metadata shown in every page.
type SiteConfig struct {
	Title       string `toml:"title"`
	Version     string `toml:"version"`
	BaseURL     string `toml:"base_url"`
	Repo        string `toml:"repo"`
	ExtraHead   string `toml:"extra_head"`   // injected verbatim into <head> before </head>
	ExtraFooter string `toml:"extra_footer"` // injected verbatim into <footer> after attribution
}

// BuildConfig controls where kwelea reads and writes files.
type BuildConfig struct {
	DocsDir   string `toml:"docs_dir"`
	OutputDir string `toml:"output_dir"`
}

// ServeConfig controls the development server.
type ServeConfig struct {
	Port        int  `toml:"port"`
	OpenBrowser bool `toml:"open_browser"`
}

// ThemeConfig controls syntax highlighting colour schemes.
type ThemeConfig struct {
	LightCodeTheme string `toml:"light_code_theme"`
	DarkCodeTheme  string `toml:"dark_code_theme"`
}

// NavEntry is a manually-ordered navigation section defined in kwelea.toml.
// When no [[nav]] sections exist the walker builds the nav automatically from
// the docs/ directory structure.
type NavEntry struct {
	Section string   `toml:"section"`
	Pages   []string `toml:"pages"`
}

// Load reads kwelea.toml from path, applies sensible defaults for any fields
// not present in the file, and returns the validated config.
func Load(path string) (*Config, error) {
	// Pre-populate defaults; toml.Decode only writes fields present in the file,
	// so missing keys naturally inherit these values.
	cfg := defaults()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf(
				"config file not found: %s\n\nRun `kwelea init` to create one.", path,
			)
		}
		return nil, fmt.Errorf("reading config %s: %w", path, err)
	}

	if _, err := toml.Decode(string(data), cfg); err != nil {
		return nil, fmt.Errorf("parsing config %s: %w", path, err)
	}

	return cfg, nil
}

// defaults returns a Config pre-populated with sensible default values.
func defaults() *Config {
	return &Config{
		Build: BuildConfig{
			DocsDir:   "docs",
			OutputDir: "site",
		},
		Serve: ServeConfig{
			Port:        4000,
			OpenBrowser: true,
		},
		Theme: ThemeConfig{
			LightCodeTheme: "github",
			DarkCodeTheme:  "github-dark",
		},
	}
}
