package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/dalebezolli/wcode/internal/matchers"
	"github.com/dalebezolli/wcode/internal/tui"
)

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

type model struct {
	matcher matchers.Matcher

	selection          int
	queryInput         []byte
	directories        []string
	queriedDirectories []string

	list  tui.Box
	input tui.Box
	info  tui.Box
}

func (m *model) View(t *tui.TUI) {
	t.Clear()

	listContent := ""
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
		if m.selection == i {
			selectedMod = ";1"
			listContent += tui.AnsiMoveTo(3, 3+i)
			listContent += "\x1b[2;1m" + SELECTED_LINE_INDICATOR + tui.ANSI_CLEAR_MODIFIER
		}

		listContent += tui.AnsiMoveTo(5, 3+i)
		listContent += fmt.Sprintf("\x1b[2%vm", selectedMod) + path + tui.ANSI_CLEAR_MODIFIER

		listContent += tui.AnsiMoveTo(len(path)+6, 3+i)
		listContent += fmt.Sprintf("\x1b[%vm", selectedMod) + project + tui.ANSI_CLEAR_MODIFIER
	}

	m.list.Content = tui.ANSI_CLEAR_MODIFIER + tui.AnsiMoveDown(1) + tui.AnsiMoveRight(1) + listContent
	t.Add(tui.ANSI_CLEAR_MODIFIER + "\x1b[2;36m")
	m.list.Render(t)

	t.MoveAt(t.Width/2+1, 0)
	m.info.Render(t)

	t.MoveAt(0, t.Height-3)
	m.input.Content = tui.ANSI_CLEAR_MODIFIER + tui.ANSI_BOLD + tui.AnsiMoveDown(1) + tui.AnsiMoveRight(1) + string(m.queryInput)
	t.Add("\x1b[2;36m")
	m.input.Render(t)

	t.Flush()
}

func (m *model) Update(e tui.Event) bool {
	result := true

	switch typedE := e.(type) {
	case tui.EventResize:
		m.list.Height = typedE.Height - m.input.Height
		m.list.Width = typedE.Width / 2

		m.input.Width = typedE.Width / 2

		m.info.Width = typedE.Width/2 - 1
		m.info.Height = typedE.Height
	case tui.EventKeyPress:
		result = m.onKeyPress(typedE)
	}

	if len(m.queryInput) != 0 {
		m.queriedDirectories = m.matcher.Match(m.directories, string(m.queryInput))
	} else {
		m.queriedDirectories = m.directories
	}

	if result {
		m.selection = min(max(m.selection, 0), len(m.queriedDirectories)-1)
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

func main() {
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

	var matcher matchers.Matcher

	if matchers.IsMatcherRGAvailable() {
		matcher = matchers.MatcherRG{}
	} else {
		matcher = matchers.MatcherLinear{}
	}

	model := &model{
		directories:        directories,
		queriedDirectories: directories,
		matcher:            matcher,
	}

	t := tui.NewTUI(model)

	model.input = tui.Box{
		Title:  "What project are you working on today?",
		Width:  t.Width / 2,
		Height: 4,
	}

	model.list = tui.Box{
		Title:  "Projects",
		Width:  t.Width / 2,
		Height: t.Height - model.input.Height,
	}

	model.info = tui.Box{
		Title:  "Info",
		Width:  t.Width/2 - 1,
		Height: t.Height,
	}

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
	return strings.Split(pathsString, " ")
}
