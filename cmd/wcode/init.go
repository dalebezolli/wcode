package main

import (
	_ "embed"
	"fmt"
	"io"
	"os"
)

//go:embed scripts/bash-zsh.sh
var shellIntegration string

func runInit(args []string) int {
	if len(args) != 1 || (args[0] != "bash" && args[0] != "zsh") {
		fmt.Fprintln(os.Stderr, "usage: wcode init <bash|zsh>")
		return EXIT_BAD_PATH
	}

	if _, err := io.WriteString(os.Stdout, shellIntegration); err != nil {
		fmt.Fprintf(os.Stderr, "wcode: writing %s integration: %v\n", args[0], err)
		return EXIT_BAD_PATH
	}
	return EXIT_OK
}
