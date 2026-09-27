package main

import (
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

const VERSION = "0.1.0"

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

   -v, --version           Display the version of wcode.`

func main() {
	flag.Usage = func() {
		fmt.Println(usage)
		os.Exit(EXIT_TERMINATED)
	}

	var flagShowVersion bool
	flag.BoolVar(&flagShowVersion, "v", false, "Display the version of wcode.")
	flag.BoolVar(&flagShowVersion, "version", false, "Display the version of wcode.")

	var flagMatcher string
	flag.StringVar(&flagMatcher, "m", "regex", "Specifies the matcher to be used.")
	flag.StringVar(&flagMatcher, "matcher", "regex", "Specifies the matcher to be used.")

	var flagDetailer string
	flag.StringVar(&flagDetailer, "d", "git", "Specifies the detailer to be used.")
	flag.StringVar(&flagDetailer, "detailer", "git", "Specifies the detailer to be used.")

	var flagNavigateOnly bool
	flag.BoolVar(&flagNavigateOnly, "n", false, "Disables tmux navigation")
	flag.BoolVar(&flagNavigateOnly, "navigate-only", false, "Disables tmux navigation")

	flag.Parse()

	if flagShowVersion {
		fmt.Println(VERSION)
		os.Exit(EXIT_TERMINATED)
	}

	err := config.EnsureDir()
	if err != nil {
		fmt.Println("An unexpected error occured while initializing the config directory:", err.Error())
		os.Exit(EXIT_BAD_PATH)
	}

	projectRoots := config.GetProjectRoots()

	directories, err := projects.DiscoverProjects(projectRoots)
	if err != nil {
		fmt.Printf("There was a problem while collecting the projects: %v\n", err.Error())
		os.Exit(EXIT_BAD_PATH)
	}

	if len(directories) == 0 {
		fmt.Println("There don't exist any projects in the directories: ", projectRoots)
		os.Exit(EXIT_NO_PROJECTS)
	}

	model := tui.NewWcodeProject(
		matchers.NewMatcher(matchers.MatcherType(flagMatcher)),
		detailers.NewDetailer(detailers.DetailerType(flagDetailer)),
		directories,
	)

	t := tui.NewTUI(model)
	defer t.Close()

	t.Run()
	t.Clear()
	t.Flush()

	selectionPath := model.SelectedPath()

	err = selection.Save(selectionPath)
	if err != nil {
		fmt.Println("An unexpected error occured while saving the selection:", err.Error())
		t.Close()
		os.Exit(EXIT_BAD_PATH)
	}

	if len(selectionPath) == 0 {
		t.Close()
		os.Exit(EXIT_NO_SELECTION)
	}
}
