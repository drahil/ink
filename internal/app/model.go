package app

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type Pane int

const (
	PaneFiles Pane = iota
	PaneEditor
	PaneCommand
)

type Model struct {
	width   int
	height  int
	focused Pane
	status  string
}

func NewModel() Model {
	return Model{
		focused: PaneEditor,
		status:  "ready",
	}
}

func (p Pane) String() string {
	switch p {
	case PaneFiles:
		return "files"
	case PaneEditor:
		return "editor"
	case PaneCommand:
		return "command"
	default:
		return "unknown"
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.status = fmt.Sprintf("terminal resized to %dx%d", m.width, m.height)
	case tea.KeyMsg:
		switch msg.String() {
		case "crtl+c", "q":
			return m, tea.Quit
		case "tab":
			m.focusNextPane()
		default:
			m.status = fmt.Sprintf("pressed %q", msg.String())
		}
	}

	return m, nil
}

func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return "starting..."
	}

	var b strings.Builder

	fmt.Fprintf(&b, "ink | focused: %s | size: %dx%d\n\n", m.focused, m.width, m.height)

	b.WriteString("[ files ]\n")
	b.WriteString("  file explorer placeholder\n\n")

	b.WriteString("[ editor ]\n")
	b.WriteString("  editor placeholder\n")
	b.WriteString("  later: open file text goes here\n\n")

	b.WriteString("[ command ]\n")
	b.WriteString("  command/status input placeholder\n\n")

	fmt.Fprintf(&b, "status: %s\n", m.status)
	fmt.Fprintf(&b, "keys: tab focus | q quit | ctrl+c quit\n")

	return b.String()
}

func (m Model) focusNextPane() {
	switch m.focused {
	case PaneFiles:
		m.focused = PaneEditor
	case PaneEditor:
		m.focused = PaneCommand
	case PaneCommand:
		m.focused = PaneFiles
	}

	m.status = fmt.Sprintf("focused %s pane", m.focused)
}
