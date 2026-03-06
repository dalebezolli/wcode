package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
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

type model struct {
	selection          int
	queryInput         []byte
	directories        []string
	queriedDirectories []string
}

func (m *model) View(tui *TUI) {
	tui.Clear()

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
			listContent += fmt.Sprintf(ANSI_MOVE_TO, 3+i, 3)
			listContent += "\x1b[2;1m•" + ANSI_CLEAR_MODIFIER
		}

		listContent += fmt.Sprintf(ANSI_MOVE_TO, 3+i, 5)
		listContent += fmt.Sprintf("\x1b[2%vm", selectedMod) + path + ANSI_CLEAR_MODIFIER

		listContent += fmt.Sprintf(ANSI_MOVE_TO, 3+i, len(path)+6)
		listContent += fmt.Sprintf("\x1b[%vm", selectedMod) + project + ANSI_CLEAR_MODIFIER
	}

	tui.Add(ANSI_CLEAR_MODIFIER + "\x1b[2;36m")
	list := Box{
		Width:   80,
		Height:  tui.Height - 4,
		Title:   "Projects",
		Content: ANSI_CLEAR_MODIFIER + "\x1b[B\x1b[C" + listContent,
	}

	list.Render(tui)

	tui.MoveAt(0, tui.Height-3)
	tui.Add("\x1b[2;36m")
	input := Box{
		Width:   80,
		Height:  4,
		Title:   "What project are you working on today?",
		Content: ANSI_CLEAR_MODIFIER + "\x1b[1m\x1b[B\x1b[C" + string(m.queryInput),
	}

	input.Render(tui)
	tui.Flush()
}

func (m *model) Update(e Event) bool {
	result := true

	switch typedE := e.(type) {
	case EventKeyPress:
		result = m.onKeyPress(typedE)
	}

	if len(m.queryInput) != 0 {
		m.queriedDirectories = getProjectMatchesRG(m.directories, string(m.queryInput))
	} else {
		m.queriedDirectories = m.directories
	}

	m.selection = min(max(m.selection, 0), len(m.queriedDirectories)-1)

	return result
}

func (m *model) onKeyPress(e EventKeyPress) bool {
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

	model := &model{
		directories:        directories,
		queriedDirectories: directories,
	}

	tui := NewTUI(model)
	defer tui.Close()

	tui.Run()
	tui.Clear()
	tui.Flush()

	selectionPath := ""
	if model.selection != -1 {
		selectionPath = model.queriedDirectories[model.selection]
	}

	err = saveSelectionToDisk(selectionPath)
	if err != nil {
		fmt.Println("An unexpected error occured while saving the selection:", err.Error())
		tui.Close()
		os.Exit(EXIT_BAD_PATH)
	}

	if len(selectionPath) == 0 {
		tui.Close()
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

func getProjectMatches(dirs []string, needle string, matchPath bool) []string {
	res := make([]string, 0, len(dirs))

	for _, hay := range dirs {
		if !isMatch(strings.ToLower(hay), strings.ToLower(needle), matchPath) {
			continue
		}

		res = append(res, hay)
	}

	return res
}

func getProjectMatchesRG(dirs []string, needle string) []string {
	echoCmd := exec.Command("echo", strings.Join(dirs, "\n"))
	rgCmd := exec.Command("rg", strings.ReplaceAll(needle, " ", ".*"))

	cmdPipe, err := echoCmd.StdoutPipe()
	if err != nil {
		return dirs
	}

	rgCmd.Stdin = cmdPipe

	echoCmd.Start()
	res, err := rgCmd.CombinedOutput()

	if err != nil {
		return dirs
	}

	return strings.Split(string(res), "\n")
}

func isMatch(haystack string, needle string, matchPath bool) bool {
	for i := len(haystack) - 1; i >= len(needle)-1; i-- {
		if !matchPath && haystack[i] == '/' {
			return false
		}

		j := 0
		for j < len(needle) && haystack[i-j] == needle[len(needle)-j-1] {
			j++
		}

		if j == len(needle) {
			return true
		}
	}

	return false
}
