package main

import (
	"fmt"
	"os"
	"strings"
)

const VERSION = "v0.1.0"

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
	default:
		fmt.Fprintf(os.Stderr, "wcode: unknown command %q\n", args[0])
		return EXIT_BAD_PATH
	}
}
