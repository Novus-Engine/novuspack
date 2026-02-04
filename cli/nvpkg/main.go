package main

import (
	"os"

	"github.com/novus-engine/novuspack/cli/nvpkg/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
