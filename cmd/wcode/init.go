package main

import (
	_ "embed"
	"fmt"
	"io"
	"os"
)

//go:embed scripts/bash.sh
var bashIntegration string

func runInit(args []string) int {
	if len(args) != 1 || args[0] != "bash" {
		fmt.Fprintln(os.Stderr, "usage: wcode init bash")
		return EXIT_BAD_PATH
	}

	if _, err := io.WriteString(os.Stdout, bashIntegration); err != nil {
		fmt.Fprintln(os.Stderr, "wcode: writing Bash integration:", err)
		return EXIT_BAD_PATH
	}
	return EXIT_OK
}
