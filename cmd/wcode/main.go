package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/dalebezolli/wcode/internal/detailers"
	"github.com/dalebezolli/wcode/internal/matchers"
	"github.com/dalebezolli/wcode/internal/tui"
)

const VERSION = "0.1.0"

const DIR_CONFIG = "$HOME/.config/wcode"
const FILENAME_SELECTION = "selection"
const ENV_PROJECT_PATHS = "$WCODE_PATHS"

const (
	EXIT_OK           = 0
	EXIT_NO_PROJECTS  = 1
	EXIT_BAD_PATH     = 2
	EXIT_NO_SELECTION = 3
	EXIT_TERMINATED   = 9
)

const SELECTED_LINE_INDICATOR = "•"
const PATH_LABEL = "Path: "
const INFO_LABEL = "Info"
const INFO_NO_DATA_LABEL = "No info available"

type model struct {
	matcher  matchers.Matcher
	detailer detailers.Detailer

	selection          int
	queryInput         []byte
	directories        []string
	queriedDirectories []string

	projectDetails map[string]detailers.Details

	// Used to reduce rendering
	prevQueriedDirectories []string
	previousSelection      int

	list  tui.Box
	input tui.Box
	info  tui.Box
}

func (m *model) Start(t *tui.TUI) {
	m.input = tui.Box{
		Title:  "What project are you working on today?",
		Width:  t.Width / 2,
		Height: 4,
	}

	m.list = tui.Box{
		Title:  "Projects",
		Width:  t.Width / 2,
		Height: t.Height - m.input.Height,
	}

	m.info = tui.Box{
		Title:  "Info",
		Width:  t.Width/2 - 1,
		Height: t.Height,
	}

	// Prefetch first directory
	m.projectDetails[m.directories[0]] = m.detailer.GetDetails(m.directories[0])
	go func() {
		for _, dir := range m.directories {
			details := m.detailer.GetDetails(dir)
			m.projectDetails[dir] = details
		}
	}()

	t.Clear()

	t.Add(tui.ANSI_CLEAR_MODIFIER + "\x1b[38;5;45m")
	m.list.Render(t)

	t.MoveAt(t.Width/2+1, 0)
	t.Add("\x1b[38;5;45m")
	m.info.Render(t)

	t.MoveAt(0, t.Height-3)
	t.Add("\x1b[38;5;45m")
	m.input.Render(t)
}

func (m *model) View(t *tui.TUI) {
	var listBuilder strings.Builder
	for i, dir := range m.queriedDirectories {
		if len(dir) == 0 {
			continue
		}

		splitPath := strings.Split(dir, "/")
		if len(splitPath) < 2 {
			continue
		}

		path := "[" + splitPath[len(splitPath)-2] + "]"
		project := splitPath[len(splitPath)-1]

		selectedMod := ""
		listBuilder.WriteString(tui.AnsiMoveTo(3, 3+i))
		if m.selection == i {
			selectedMod = ";1"
			listBuilder.WriteString("\x1b[2;1m")
			listBuilder.WriteString(SELECTED_LINE_INDICATOR)
			listBuilder.WriteString(tui.ANSI_CLEAR_MODIFIER)
		} else {
			listBuilder.WriteString(" ")
		}

		listBuilder.WriteString(tui.AnsiMoveTo(5, 3+i))
		listBuilder.WriteString(fmt.Sprintf("\x1b[2%vm", selectedMod))
		listBuilder.WriteString(path)
		listBuilder.WriteString(tui.ANSI_CLEAR_MODIFIER)

		listBuilder.WriteString(fmt.Sprintf("\x1b[%vm ", selectedMod))
		listBuilder.WriteString(project)
		listBuilder.WriteString(tui.ANSI_CLEAR_MODIFIER)
		listBuilder.WriteString(strings.Repeat(" ", max(0, (t.Width/2-6)-(len(project)+len(path)))))
	}

	for i := 0; i < len(m.prevQueriedDirectories)-len(m.queriedDirectories); i++ {
		listBuilder.WriteString(tui.AnsiMoveTo(3, 3+i+len(m.queriedDirectories)))
		listBuilder.WriteString(strings.Repeat(" ", max(0, (t.Width/2-6))))
	}

	t.Add(tui.ANSI_CLEAR_MODIFIER)
	t.Add(tui.AnsiMoveDown(1))
	t.Add(tui.AnsiMoveRight(1))
	t.Add(listBuilder.String())

	t.MoveAt(t.Width/2+1, 0)
	t.Add(tui.AnsiMoveDown(1))
	t.Add(tui.AnsiMoveRight(1))
	t.Add(tui.ANSI_CLEAR_MODIFIER)
	t.Add(tui.ANSI_BOLD)
	t.Add(tui.AnsiMoveDown(1))
	t.Add(m.displayDetails(m.queriedDirectories[m.selection], t))

	t.MoveAt(0, t.Height-3)
	t.Add(tui.AnsiMoveDown(1))
	t.Add(tui.AnsiMoveRight(1))
	t.Add(tui.ANSI_CLEAR_MODIFIER)
	t.Add(tui.AnsiMoveDown(1))
	t.Add(tui.AnsiMoveRight(1))
	t.Add(string(m.queryInput) + strings.Repeat(" ", max(0, m.input.Width-len(m.queryInput)-4)))
	t.Add(tui.AnsiMoveLeft(m.input.Width - len(m.queryInput) - 4))

	t.Flush()
}

var detailColors = []string{
	"33",
	"69",
	"105",
	"141",
	"177",
	"213",
}

func getCleanTitle(title string, rowMaxLen int) string {
	if len(title) > rowMaxLen {
		return string([]byte(title)[:max(0, rowMaxLen-4)]) + "..."
	}

	return title
}

func (m *model) displayDetails(dir string, t *tui.TUI) string {
	y := 3
	x := t.Width/2 + 3

	details := m.projectDetails[dir]
	prevDetails := m.projectDetails[m.prevQueriedDirectories[m.previousSelection]]
	rowMaxLen := max(0, t.Width/2-4)

	prevCleanedTitle := getCleanTitle(prevDetails.Title, rowMaxLen)
	cleanedTitle := getCleanTitle(details.Title, rowMaxLen)

	rowTitle := cleanedTitle + strings.Repeat(" ", max(0, len(prevCleanedTitle)-len(cleanedTitle)))
	path := fmt.Sprintf(tui.ANSI_MOVE_TO, y+1, x) + PATH_LABEL + details.Path + strings.Repeat(" ", max(0, len(prevDetails.Path)-len(details.Path)))

	detailsString := fmt.Sprintf(tui.ANSI_MOVE_TO, y, x) + tui.ANSI_BOLD + rowTitle + tui.ANSI_CLEAR_MODIFIER +
		fmt.Sprintf(tui.ANSI_MOVE_TO, y+1, x) + "\x1b[38;5;243m" + path + tui.ANSI_CLEAR_MODIFIER +
		fmt.Sprintf(tui.ANSI_MOVE_TO, y+3, x) + tui.ANSI_BOLD + INFO_LABEL + tui.ANSI_CLEAR_MODIFIER

	// TODO: Maybe there's a better way to do this instead of cleaning the entire screen?
	for i := min(1, len(details.Rest)); i < len(prevDetails.Rest); i++ {
		detailsString += fmt.Sprintf(tui.ANSI_MOVE_TO, y+4+i, x) + strings.Repeat(" ", rowMaxLen)
	}

	if len(details.Rest) == 0 {
		detailsString += fmt.Sprintf(tui.ANSI_MOVE_TO, y+4, x) + "\x1b[38;5;243m" + INFO_NO_DATA_LABEL + strings.Repeat(" ", max(0, rowMaxLen-len(INFO_NO_DATA_LABEL)))
	} else {
		order := m.detailer.GetRestOrder()

		displayedIndex := 0
		for i, key := range order {
			val, exists := details.Rest[key]
			if !exists {
				continue
			}

			detailsContent := key + val + strings.Repeat(" ", max(0, rowMaxLen-len(key)-len(val)))
			detailsString += fmt.Sprintf(tui.ANSI_MOVE_TO, y+4+displayedIndex, x) + "\x1b[38;5;" + detailColors[i%len(detailColors)] + "m" + detailsContent + tui.ANSI_CLEAR_MODIFIER
			displayedIndex++
		}
	}

	return detailsString
}

func (m *model) Update(e tui.Event, t *tui.TUI) bool {
	result := true

	m.previousSelection = m.selection

	switch typedE := e.(type) {
	case tui.EventResize:
		m.list.Height = typedE.Height - m.input.Height
		m.list.Width = typedE.Width / 2

		m.input.Width = typedE.Width / 2

		m.info.Width = typedE.Width/2 - 1
		m.info.Height = typedE.Height

		t.Clear()

		t.Add(tui.ANSI_CLEAR_MODIFIER + "\x1b[38;5;45m")
		m.list.Render(t)

		t.MoveAt(t.Width/2+1, 0)
		t.Add("\x1b[38;5;45m")
		m.info.Render(t)

		t.MoveAt(0, t.Height-3)
		t.Add("\x1b[38;5;45m")
		m.input.Render(t)
	case tui.EventKeyPress:
		result = m.onKeyPress(typedE)
	}

	if len(m.queryInput) != 0 {
		m.queriedDirectories = m.matcher.Match(m.directories, string(m.queryInput))
	} else {
		m.prevQueriedDirectories = m.queriedDirectories
		m.queriedDirectories = m.directories
	}

	if result {
		m.selection = (m.selection + len(m.queriedDirectories)) % len(m.queriedDirectories)
	}

	return result
}

func (m *model) onKeyPress(e tui.EventKeyPress) bool {
	switch e.ReadBuffer[0] {
	case '\x7F':
		if len(m.queryInput) == 0 {
			break
		}
		m.queryInput = m.queryInput[0 : len(m.queryInput)-1]
	case '\x0E', '\x04':
		m.selection++
	case '\x10', '\x15':
		m.selection--
	case '\x03', '\x18':
		m.selection = -1
		return false
	case '\x0D':
		return false
	case '\x1B':

		if e.ReadBuffer[1] == '\x7F' {
			foundSpace := false
			foundWord := false
			i := len(m.queryInput) - 1
			for i >= 0 && !foundSpace {
				if foundWord && m.queryInput[i] == '\x20' {
					foundSpace = true
				} else {
					foundWord = true
					i--
				}
			}

			m.queryInput = m.queryInput[0 : i+1]
			m.selection = 0
		}

		switch e.ReadBuffer[2] {
		case '\x41':
			m.selection--
		case '\x42':
			m.selection++
		}
	default:
		m.queryInput = append(m.queryInput, e.ReadBuffer[0])
		m.selection = 0
	}

	return true
}

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

	err := setupFiles()
	if err != nil {
		fmt.Println("An unexpected error occured while initializing the config directory:", err.Error())
		os.Exit(EXIT_BAD_PATH)
	}

	projectRoots := gatherProjectPaths()

	directories, err := gatherProjects(projectRoots)
	if err != nil {
		fmt.Printf("There was a problem while collecting the projects: %v\n", err.Error())
		os.Exit(EXIT_BAD_PATH)
	}

	if len(directories) == 0 {
		fmt.Println("There don't exist any projects in the directories: ", projectRoots)
		os.Exit(EXIT_NO_PROJECTS)
	}

	model := &model{
		matcher:  matchers.NewMatcher(matchers.MatcherType(flagMatcher)),
		detailer: detailers.NewDetailer(detailers.DetailerType(flagDetailer)),

		directories:            directories,
		queriedDirectories:     directories,
		prevQueriedDirectories: directories,

		projectDetails: make(map[string]detailers.Details),
	}

	t := tui.NewTUI(model)
	defer t.Close()

	t.Run()
	t.Clear()
	t.Flush()

	selectionPath := ""
	if model.selection != -1 {
		selectionPath = model.queriedDirectories[model.selection]
	}

	err = saveSelectionToDisk(selectionPath)
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

func gatherProjects(roots []string) ([]string, error) {
	directories := []string{}
	for _, root := range roots {
		entries, err := os.ReadDir(root)
		if err != nil {
			return nil, err
		}

		for _, v := range entries {
			if v.Type().IsDir() == false {
				continue
			}

			projectPath := strings.ReplaceAll(root+string(os.PathSeparator)+v.Name(), "//", "/")
			directories = append(directories, projectPath)
		}
	}

	return directories, nil
}

func setupFiles() error {
	baseDir := os.ExpandEnv(DIR_CONFIG)
	return os.MkdirAll(baseDir, 0751)
}

func saveSelectionToDisk(selection string) error {
	baseDir := os.ExpandEnv(DIR_CONFIG)
	file, err := os.Create(baseDir + string(os.PathSeparator) + FILENAME_SELECTION)
	if err == nil {
		file.Write([]byte(selection))
	}

	return err
}

func gatherProjectPaths() []string {
	pathsString := os.ExpandEnv(ENV_PROJECT_PATHS)
	return strings.Split(pathsString, ";")
}
