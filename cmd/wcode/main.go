package main

import (
	"fmt"
	"os"
	"strings"
)

const VERSION = "v0.1.0"

const (
	EXIT_OK           = 0
	EXIT_NO_PROJECTS  = 1
	EXIT_BAD_PATH     = 2
	EXIT_NO_SELECTION = 3
	EXIT_TERMINATED   = 9
)

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	if len(args) == 0 || strings.HasPrefix(args[0], "-") {
		return runSelector(args)
	}

	switch args[0] {
	case "init":
		return runInit(args[1:])
	case "selection":
		return runGetSelection(args[1:])
	default:
		fmt.Fprintf(os.Stderr, "wcode: unknown command %q\n", args[0])
		return EXIT_BAD_PATH
	}
}
