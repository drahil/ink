package app

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

var (
	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Background(lipgloss.Color("62")).
			Padding(0, 1)

	paneStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(0, 1)

	focusedPaneStyle = paneStyle.
				BorderForeground(lipgloss.Color("39"))

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Background(lipgloss.Color("236")).
			Padding(0, 1)
)

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
		case "ctrl+c", "q":
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

	header := m.renderHeader()
	status := m.renderStatusBar()

	headerHeight := lipgloss.Height(header)
	statusHeight := lipgloss.Height(status)
	commandHeight := 5
	mainHeight := max(3, m.height-headerHeight-statusHeight-commandHeight)
	filesWidth := min(28, max(18, m.width/4))
	editorWidth := max(20, m.width-filesWidth)

	editorPane := m.renderEditorPane(editorWidth, mainHeight)
	filesPane := m.renderFilesPane(filesWidth, mainHeight)
	main := lipgloss.JoinHorizontal(lipgloss.Top, editorPane, filesPane)
	command := m.renderCommandPane(m.width, commandHeight)

	return lipgloss.JoinVertical(lipgloss.Left, header, main, command, status)
}

func (m Model) renderHeader() string {
	text := fmt.Sprintf("ink | focused: %s | size: %dx%d", m.focused, m.width, m.height)

	return headerStyle.
		Width(max(0, m.width-2)).
		Render(text)
}

func (m Model) renderFilesPane(width, height int) string {
	var b strings.Builder

	b.WriteString(paneTitle("files", m.focused == PaneFiles))
	b.WriteString("\n\n")
	b.WriteString("project root\n")
	b.WriteString("cmd/\n")
	b.WriteString("internal/\n")
	b.WriteString("go.mod")

	return m.renderPane(width, height, b.String(), m.focused == PaneFiles)
}

func (m Model) renderEditorPane(width, height int) string {
	var b strings.Builder

	b.WriteString(paneTitle("editor", m.focused == PaneEditor))
	b.WriteString("\n\n")
	b.WriteString("editor placeholder\n")
	b.WriteString("later: open file text goes here\n\n")
	b.WriteString("tab changes focus")

	return m.renderPane(width, height, b.String(), m.focused == PaneEditor)
}

func (m Model) renderCommandPane(width, height int) string {
	var b strings.Builder

	b.WriteString(paneTitle("command", m.focused == PaneCommand))
	b.WriteString("\n\n")
	b.WriteString("command/status input placeholder")

	return m.renderPane(width, height, b.String(), m.focused == PaneCommand)
}

func (m Model) renderStatusBar() string {
	text := fmt.Sprintf("status: %s | keys: tab focus | q quit | ctrl+c quit", m.status)

	return statusStyle.
		Width(max(0, m.width-2)).
		Render(text)
}

func (m Model) renderPane(width, height int, content string, focused bool) string {
	style := paneStyle
	if focused {
		style = focusedPaneStyle
	}

	return style.
		Width(max(0, width-4)).
		Height(max(0, height-2)).
		Render(content)
}

func paneTitle(name string, active bool) string {
	if active {
		return fmt.Sprintf("[ %s * ]", name)
	}

	return fmt.Sprintf("[ %s ]", name)
}

func (m *Model) focusNextPane() {
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

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}
