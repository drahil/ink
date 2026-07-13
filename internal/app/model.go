package app

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/drahil/ink/internal/ui"
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
	files   ui.FilesPane
	editor  ui.EditorPane
	command ui.CommandPane
}

func NewModel() Model {
	return Model{
		focused: PaneEditor,
		status:  "ready",
		files:   ui.NewFilesPane(),
		editor:  ui.NewEditorPane(),
		command: ui.NewCommandPane(),
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
		case "up":
			if m.focused == PaneEditor {
				m.editor.MoveCursorUp()
				m.status = m.editorCursorStatus()
			}
		case "down":
			if m.focused == PaneEditor {
				m.editor.MoveCursorDown()
				m.status = m.editorCursorStatus()
			}
		case "right":
			if m.focused == PaneEditor {
				m.editor.MoveCursorRight()
				m.status = m.editorCursorStatus()
			}
		case "left":
			if m.focused == PaneEditor {
				m.editor.MoveCursorLeft()
				m.status = m.editorCursorStatus()
			}
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

	header := ui.RenderHeader(m.width, m.focused.String())
	status := ui.RenderStatusBar(m.width, m.status)

	headerHeight := lipgloss.Height(header)
	statusHeight := lipgloss.Height(status)
	commandHeight := 5
	mainHeight := max(3, m.height-headerHeight-statusHeight-commandHeight)
	filesWidth := min(28, max(18, m.width/4))
	editorWidth := max(20, m.width-filesWidth)

	editorPane := m.editor.View(editorWidth, mainHeight, m.focused == PaneEditor)
	filesPane := m.files.View(filesWidth, mainHeight, m.focused == PaneFiles)
	main := lipgloss.JoinHorizontal(lipgloss.Top, editorPane, filesPane)
	command := m.command.View(m.width, commandHeight, m.focused == PaneCommand)

	return lipgloss.JoinVertical(lipgloss.Left, header, main, command, status)
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

func (m Model) editorCursorStatus() string {
	return fmt.Sprintf("row:%d --- column:%d", m.editor.Cursor.Row, m.editor.Cursor.Column)
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
