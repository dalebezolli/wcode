package main

import (
	_ "embed"
	"fmt"
	"io"
	"os"
)

//go:embed scripts/bash-zsh.sh
var bashZshIntegration string

//go:embed scripts/fish.fish
var fishIntegration string

func runInit(args []string) int {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "usage: wcode init <bash|zsh|fish>")
		return EXIT_BAD_PATH
	}

	var integration string
	switch args[0] {
	case "bash", "zsh":
		integration = bashZshIntegration
	case "fish":
		integration = fishIntegration
	default:
		fmt.Fprintln(os.Stderr, "usage: wcode init <bash|zsh|fish>")
		return EXIT_BAD_PATH
	}

	if _, err := io.WriteString(os.Stdout, integration); err != nil {
		fmt.Fprintf(os.Stderr, "wcode: writing %s integration: %v\n", args[0], err)
		return EXIT_BAD_PATH
	}
	return EXIT_OK
}
