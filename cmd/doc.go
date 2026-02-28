// Package cmd implements the kwelea CLI subcommands using Cobra.
//
// Three subcommands are registered on the root command:
//
//   - build   reads kwelea.toml, walks the docs directory, parses all
//     Markdown files, and writes a self-contained static HTML site
//     to the configured output directory.
//
//   - serve   runs the build pipeline then starts a local HTTP dev server
//     with WebSocket live reload. The browser is opened automatically
//     unless serve.open_browser = false in kwelea.toml.
//
//   - init    scaffolds a minimal kwelea.toml and a docs/index.md starter
//     file in the current directory. It is a no-op if kwelea.toml
//     already exists.
//
// The package exports a single function, [Execute], which is called by
// main.go to hand off control to Cobra.
package cmd
