package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/dalebezolli/wcode/internal/config"
	"github.com/dalebezolli/wcode/internal/detailers"
	"github.com/dalebezolli/wcode/internal/matchers"
	"github.com/dalebezolli/wcode/internal/projects"
	"github.com/dalebezolli/wcode/internal/selection"
	"github.com/dalebezolli/wcode/internal/tui"
)

const (
	EXIT_OK           = 0
	EXIT_NO_PROJECTS  = 1
	EXIT_BAD_PATH     = 2
	EXIT_NO_SELECTION = 3
	EXIT_TERMINATED   = 9
)

const usage = `                     wcode - Unf*ck project navigation

 Usage: wcode [-m | --matcher <matcher_type>] [-d | --detailer <detailer_type>]
              [-n | --navigate-only] [-h | --help] [-v | --version]
        wcode init bash


 Navigate through your ocean of projects in a simple and effective way.

 Looks for projects under the directories defined in $WCODE_PATHS, a semicolon
 separated list of directories, and displays them in a list.

 If tmux is installed once a project is selected.
   - It checks if a session already exists and navigates to it,
   - otherwise it spins up a new session which runs command from the local
     .tmux.conf if it exists

 Controls:

  Arrow Up or CTRL-p      - Navigate up the list
  Arrow Down or CTRL-n    - Navigate down the list
  Enter                   - Open the selected project
  Any other key           - Type in the search

 Arguments:

   -h, --help              Displays the help message

   -m, --matcher           Specifies the matcher to be used. Default: regex
                           Available Options: linear, regex

   -d, --detailer          Specifies the detailer to be used. Default: git
                           Available Options: git, common

   -n, --navigate-only     Disables tmux navigation.

   -v, --version           Display the version of wcode.

 Commands:

   init bash               Print Bash integration for eval "$(wcode init bash)".`

func runSelector(args []string) int {
	flags := flag.NewFlagSet("wcode", flag.ContinueOnError)
	flags.Usage = func() { fmt.Fprintln(os.Stdout, usage) }

	var showVersion bool
	flags.BoolVar(&showVersion, "v", false, "Display the version of wcode.")
	flags.BoolVar(&showVersion, "version", false, "Display the version of wcode.")

	var matcherName string
	flags.StringVar(&matcherName, "m", "regex", "Specifies the matcher to be used.")
	flags.StringVar(&matcherName, "matcher", "regex", "Specifies the matcher to be used.")

	var detailerName string
	flags.StringVar(&detailerName, "d", "git", "Specifies the detailer to be used.")
	flags.StringVar(&detailerName, "detailer", "git", "Specifies the detailer to be used.")

	var navigateOnly bool
	flags.BoolVar(&navigateOnly, "n", false, "Disables tmux navigation")
	flags.BoolVar(&navigateOnly, "navigate-only", false, "Disables tmux navigation")

	if err := flags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return EXIT_TERMINATED
		}
		return EXIT_BAD_PATH
	}
	if flags.NArg() != 0 {
		fmt.Fprintf(os.Stderr, "wcode: unexpected argument %q\n", flags.Arg(0))
		return EXIT_BAD_PATH
	}
	if showVersion {
		fmt.Println(VERSION)
		return EXIT_TERMINATED
	}

	if err := config.EnsureDir(); err != nil {
		fmt.Fprintln(os.Stderr, "An unexpected error occured while initializing the config directory:", err)
		return EXIT_BAD_PATH
	}

	projectRoots := config.GetProjectRoots()
	directories, err := projects.DiscoverProjects(projectRoots)
	if err != nil {
		fmt.Fprintln(os.Stderr, "There was a problem while collecting the projects:", err)
		return EXIT_BAD_PATH
	}
	if len(directories) == 0 {
		fmt.Fprintln(os.Stderr, "There don't exist any projects in the directories:", projectRoots)
		return EXIT_NO_PROJECTS
	}

	model := tui.NewWcodeProject(
		matchers.NewMatcher(matchers.MatcherType(matcherName)),
		detailers.NewDetailer(detailers.DetailerType(detailerName)),
		directories,
	)

	t := tui.NewTUI(model)
	defer t.Close()
	t.Run()
	t.Clear()
	t.Flush()

	selectionPath := model.SelectedPath()
	if err := selection.Save(selectionPath); err != nil {
		fmt.Fprintln(os.Stderr, "An unexpected error occured while saving the selection:", err)
		return EXIT_BAD_PATH
	}
	if selectionPath == "" {
		return EXIT_NO_SELECTION
	}
	return EXIT_OK
}
