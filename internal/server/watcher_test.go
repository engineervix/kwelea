package server

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/engineervix/kwelea/internal/config"
)

// --- isConfigFile ---

func TestIsConfigFile_SamePath(t *testing.T) {
	tmp := t.TempDir()
	cfgPath := filepath.Join(tmp, "kwelea.toml")
	if !isConfigFile(cfgPath, cfgPath) {
		t.Fatal("same path should match")
	}
}

func TestIsConfigFile_RelativeVsAbsolute(t *testing.T) {
	abs, _ := filepath.Abs("kwelea.toml")
	if !isConfigFile("kwelea.toml", abs) {
		t.Fatal("relative vs absolute same file should match")
	}
}

func TestIsConfigFile_DifferentFile(t *testing.T) {
	tmp := t.TempDir()
	a := filepath.Join(tmp, "kwelea.toml")
	b := filepath.Join(tmp, "other.toml")
	if isConfigFile(a, b) {
		t.Fatal("different files should not match")
	}
}

// --- reloadConfig ---

func writeTOML(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func TestReloadConfig_UpdatesDocsDir(t *testing.T) {
	tmp := t.TempDir()
	cfgPath := filepath.Join(tmp, "kwelea.toml")
	writeTOML(t, cfgPath, `
[build]
docs_dir = "new/docs"
`)

	got, err := reloadConfig(cfgPath, func(c *config.Config) error { return nil })
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Build.DocsDir != "new/docs" {
		t.Errorf("docs_dir: got %q, want %q", got.Build.DocsDir, "new/docs")
	}
}

func TestReloadConfig_OverrideWins(t *testing.T) {
	tmp := t.TempDir()
	cfgPath := filepath.Join(tmp, "kwelea.toml")
	writeTOML(t, cfgPath, `
[build]
docs_dir = "toml/docs"
`)

	got, err := reloadConfig(cfgPath, func(c *config.Config) error {
		c.Build.DocsDir = "cli/docs"
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Build.DocsDir != "cli/docs" {
		t.Errorf("override wins: got %q, want %q", got.Build.DocsDir, "cli/docs")
	}
}

func TestReloadConfig_BadTOML(t *testing.T) {
	tmp := t.TempDir()
	cfgPath := filepath.Join(tmp, "kwelea.toml")
	writeTOML(t, cfgPath, `not valid toml ][`)

	_, err := reloadConfig(cfgPath, func(c *config.Config) error { return nil })
	if err == nil {
		t.Fatal("expected error for bad TOML, got nil")
	}
}

func TestReloadConfig_MissingFile(t *testing.T) {
	_, err := reloadConfig("/nonexistent/path/kwelea.toml", func(c *config.Config) error { return nil })
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestReloadConfig_OverrideError(t *testing.T) {
	tmp := t.TempDir()
	cfgPath := filepath.Join(tmp, "kwelea.toml")
	writeTOML(t, cfgPath, ``)

	want := errors.New("override failed")
	_, err := reloadConfig(cfgPath, func(c *config.Config) error { return want })
	if !errors.Is(err, want) {
		t.Fatalf("expected override error, got %v", err)
	}
}
