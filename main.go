package main

import (
	"embed"
	"fmt"
	"os"

	"github.com/engineervix/kwelea/cmd"
)

//go:embed assets templates
var embeddedFS embed.FS

func main() {
	if err := cmd.Execute(embeddedFS); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
