package main

import (
	"embed"
	"fmt"
	"os"

	"github.com/engineervix/kwelea/cmd"
)

//go:embed assets templates
var embeddedFS embed.FS

// Set by -ldflags at build time.
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	if err := cmd.Execute(embeddedFS, version, commit, date); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
