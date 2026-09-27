package tui

import (
	"fmt"
	"strings"

	"github.com/dalebezolli/wcode/internal/detailers"
	"github.com/dalebezolli/wcode/internal/matchers"
)

func NewWcodeProject(matcher matchers.Matcher, detailer detailers.Detailer, directories []string) *ProjectModel {
	return &ProjectModel{
		matcher:                matcher,
		detailer:               detailer,
		directories:            directories,
		queriedDirectories:     directories,
		prevQueriedDirectories: directories,
		projectDetails:         make(map[string]detailers.Details),
	}
}

func (m *ProjectModel) SelectedPath() string {
	if m.selection < 0 || m.selection >= len(m.queriedDirectories) {
		return ""
	}
	return m.queriedDirectories[m.selection]
}

const SELECTED_LINE_INDICATOR = "•"
const PATH_LABEL = "Path: "
const INFO_LABEL = "Info"
const INFO_NO_DATA_LABEL = "No info available"

type ProjectModel struct {
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

	list  Box
	input Box
	info  Box
}

func (m *ProjectModel) Start(t *TUI) {
	m.input = Box{
		Title:  "What project are you working on today?",
		Width:  t.Width / 2,
		Height: 4,
	}

	m.list = Box{
		Title:  "Projects",
		Width:  t.Width / 2,
		Height: t.Height - m.input.Height,
	}

	m.info = Box{
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

	t.Add(ANSI_CLEAR_MODIFIER + "\x1b[38;5;45m")
	m.list.Render(t)

	t.MoveAt(t.Width/2+1, 0)
	t.Add("\x1b[38;5;45m")
	m.info.Render(t)

	t.MoveAt(0, t.Height-3)
	t.Add("\x1b[38;5;45m")
	m.input.Render(t)
}

func (m *ProjectModel) View(t *TUI) {
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
		listBuilder.WriteString(AnsiMoveTo(3, 3+i))
		if m.selection == i {
			selectedMod = ";1"
			listBuilder.WriteString("\x1b[2;1m")
			listBuilder.WriteString(SELECTED_LINE_INDICATOR)
			listBuilder.WriteString(ANSI_CLEAR_MODIFIER)
		} else {
			listBuilder.WriteString(" ")
		}

		listBuilder.WriteString(AnsiMoveTo(5, 3+i))
		listBuilder.WriteString(fmt.Sprintf("\x1b[2%vm", selectedMod))
		listBuilder.WriteString(path)
		listBuilder.WriteString(ANSI_CLEAR_MODIFIER)

		listBuilder.WriteString(fmt.Sprintf("\x1b[%vm ", selectedMod))
		listBuilder.WriteString(project)
		listBuilder.WriteString(ANSI_CLEAR_MODIFIER)
		listBuilder.WriteString(strings.Repeat(" ", max(0, (t.Width/2-6)-(len(project)+len(path)))))
	}

	for i := 0; i < len(m.prevQueriedDirectories)-len(m.queriedDirectories); i++ {
		listBuilder.WriteString(AnsiMoveTo(3, 3+i+len(m.queriedDirectories)))
		listBuilder.WriteString(strings.Repeat(" ", max(0, (t.Width/2-6))))
	}

	t.Add(ANSI_CLEAR_MODIFIER)
	t.Add(AnsiMoveDown(1))
	t.Add(AnsiMoveRight(1))
	t.Add(listBuilder.String())

	t.MoveAt(t.Width/2+1, 0)
	t.Add(AnsiMoveDown(1))
	t.Add(AnsiMoveRight(1))
	t.Add(ANSI_CLEAR_MODIFIER)
	t.Add(ANSI_BOLD)
	t.Add(AnsiMoveDown(1))
	t.Add(m.displayDetails(m.queriedDirectories[m.selection], t))

	t.MoveAt(0, t.Height-3)
	t.Add(AnsiMoveDown(1))
	t.Add(AnsiMoveRight(1))
	t.Add(ANSI_CLEAR_MODIFIER)
	t.Add(AnsiMoveDown(1))
	t.Add(AnsiMoveRight(1))
	t.Add(string(m.queryInput) + strings.Repeat(" ", max(0, m.input.Width-len(m.queryInput)-4)))
	t.Add(AnsiMoveLeft(m.input.Width - len(m.queryInput) - 4))

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

func (m *ProjectModel) displayDetails(dir string, t *TUI) string {
	y := 3
	x := t.Width/2 + 3

	details := m.projectDetails[dir]
	prevDetails := m.projectDetails[m.prevQueriedDirectories[m.previousSelection]]
	rowMaxLen := max(0, t.Width/2-4)

	prevCleanedTitle := getCleanTitle(prevDetails.Title, rowMaxLen)
	cleanedTitle := getCleanTitle(details.Title, rowMaxLen)

	rowTitle := cleanedTitle + strings.Repeat(" ", max(0, len(prevCleanedTitle)-len(cleanedTitle)))
	path := fmt.Sprintf(ANSI_MOVE_TO, y+1, x) + PATH_LABEL + details.Path + strings.Repeat(" ", max(0, len(prevDetails.Path)-len(details.Path)))

	detailsString := fmt.Sprintf(ANSI_MOVE_TO, y, x) + ANSI_BOLD + rowTitle + ANSI_CLEAR_MODIFIER +
		fmt.Sprintf(ANSI_MOVE_TO, y+1, x) + "\x1b[38;5;243m" + path + ANSI_CLEAR_MODIFIER +
		fmt.Sprintf(ANSI_MOVE_TO, y+3, x) + ANSI_BOLD + INFO_LABEL + ANSI_CLEAR_MODIFIER

	// TODO: Maybe there's a better way to do this instead of cleaning the entire screen?
	for i := min(1, len(details.Rest)); i < len(prevDetails.Rest); i++ {
		detailsString += fmt.Sprintf(ANSI_MOVE_TO, y+4+i, x) + strings.Repeat(" ", rowMaxLen)
	}

	if len(details.Rest) == 0 {
		detailsString += fmt.Sprintf(ANSI_MOVE_TO, y+4, x) + "\x1b[38;5;243m" + INFO_NO_DATA_LABEL + strings.Repeat(" ", max(0, rowMaxLen-len(INFO_NO_DATA_LABEL)))
	} else {
		order := m.detailer.GetRestOrder()

		displayedIndex := 0
		for i, key := range order {
			val, exists := details.Rest[key]
			if !exists {
				continue
			}

			detailsContent := key + val + strings.Repeat(" ", max(0, rowMaxLen-len(key)-len(val)))
			detailsString += fmt.Sprintf(ANSI_MOVE_TO, y+4+displayedIndex, x) + "\x1b[38;5;" + detailColors[i%len(detailColors)] + "m" + detailsContent + ANSI_CLEAR_MODIFIER
			displayedIndex++
		}
	}

	return detailsString
}

func (m *ProjectModel) Update(e Event, t *TUI) bool {
	result := true

	m.previousSelection = m.selection

	switch typedE := e.(type) {
	case EventResize:
		m.list.Height = typedE.Height - m.input.Height
		m.list.Width = typedE.Width / 2

		m.input.Width = typedE.Width / 2

		m.info.Width = typedE.Width/2 - 1
		m.info.Height = typedE.Height

		t.Clear()

		t.Add(ANSI_CLEAR_MODIFIER + "\x1b[38;5;45m")
		m.list.Render(t)

		t.MoveAt(t.Width/2+1, 0)
		t.Add("\x1b[38;5;45m")
		m.info.Render(t)

		t.MoveAt(0, t.Height-3)
		t.Add("\x1b[38;5;45m")
		m.input.Render(t)
	case EventKeyPress:
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

func (m *ProjectModel) onKeyPress(e EventKeyPress) bool {
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
