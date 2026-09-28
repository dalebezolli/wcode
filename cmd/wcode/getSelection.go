package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/dalebezolli/wcode/internal/selection"
)

func runGetSelection(args []string) int {
	if len(args) != 0 {
		fmt.Fprintln(os.Stderr, "usage: wcode selection")
		return EXIT_BAD_PATH
	}

	selected, err := selection.Load()
	if errors.Is(err, os.ErrNotExist) || (err == nil && selected == "") {
		fmt.Fprintln(os.Stderr, "wcode: no project selected")
		return EXIT_NO_SELECTION
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "wcode: reading selection:", err)
		return EXIT_BAD_PATH
	}
	if _, err := fmt.Fprintln(os.Stdout, selected); err != nil {
		fmt.Fprintln(os.Stderr, "wcode: writing selection:", err)
		return EXIT_BAD_PATH
	}
	return EXIT_OK
}
