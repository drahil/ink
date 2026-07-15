package app

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/drahil/ink/internal/project"
	"github.com/drahil/ink/internal/ui"
)

type Pane int

const (
	PaneFiles Pane = iota
	PaneEditor
	PaneCommand
)

type cursorBlinkMsg time.Time

const cursorBlinkInterval = 500 * time.Millisecond

type Model struct {
	width         int
	height        int
	focused       Pane
	status        string
	cursorVisible bool
	files         ui.FilesPane
	editor        ui.EditorPane
	command       ui.CommandPane
}

func NewModel() Model {
	items, err := project.ListFiles(".")
	status := "ready"
	if err != nil {
		status = "could not load project files"
	}

	return Model{
		focused:       PaneEditor,
		status:        status,
		cursorVisible: true,
		files:         ui.NewFilesPane(items),
		editor:        ui.NewEditorPane(),
		command:       ui.NewCommandPane(),
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
	return blinkCursor()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.status = fmt.Sprintf("terminal resized to %dx%d", m.width, m.height)
	case cursorBlinkMsg:
		if m.focused == PaneEditor || m.focused == PaneFiles {
			m.cursorVisible = !m.cursorVisible
		} else {
			m.cursorVisible = false
		}

		return m, blinkCursor()
	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "tab":
		m.focusNextPane()
		return m, nil
	}

	switch m.focused {
	case PaneEditor:
		m.handleEditorKey(msg)
	case PaneFiles:
		m.handleFilesKey(msg)
	case PaneCommand:
		m.handleCommandKey(msg)
	}

	return m, nil
}

func (m *Model) handleEditorKey(msg tea.KeyMsg) {
	switch msg.String() {
	case "up":
		m.editor.MoveCursorUp()
	case "down":
		m.editor.MoveCursorDown()
	case "right":
		m.editor.MoveCursorRight()
	case "left":
		m.editor.MoveCursorLeft()
	case "backspace":
		m.editor.Backspace()
	case "enter":
		m.editor.InsertNewline()
	default:
		if len(msg.Runes) == 0 {
			m.status = fmt.Sprintf("pressed %q", msg.String())
			return
		}

		m.editor.InsertRune(msg.Runes[0])
	}

	m.cursorVisible = true
	m.status = m.editorCursorStatus()
}

func (m *Model) handleFilesKey(msg tea.KeyMsg) {
	switch msg.String() {
	case "backspace":
		m.files.Backspace()
	case "esc":
		m.files.ClearSearch()
	case "up":
		m.files.MoveSelectionUp()
	case "down":
		m.files.MoveSelectionDown()
	case "enter":
		item := m.files.SelectedItem()
		if item == "" {
			m.status = "no file selected"
			return
		}

		content, err := project.ReadFile(".", item)
		if err != nil {
			m.status = "could not open: " + item
			return
		}

		m.editor.OpenContent(item, content)
		m.focused = PaneEditor
		m.cursorVisible = true
		m.status = "opened: " + item
		return
	default:
		if len(msg.Runes) == 0 {
			m.status = fmt.Sprintf("pressed %q", msg.String())
			return
		}

		m.files.Query(msg.Runes[0])
	}

	m.status = "files search: " + m.files.SearchQuery
}

func (m *Model) handleCommandKey(msg tea.KeyMsg) {
	m.status = fmt.Sprintf("command pressed %q", msg.String())
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

	editorPane := m.editor.View(editorWidth, mainHeight, m.focused == PaneEditor, m.cursorVisible)
	filesPane := m.files.View(filesWidth, mainHeight, m.focused == PaneFiles, m.cursorVisible)
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
	m.cursorVisible = m.focused == PaneEditor
}

func (m Model) editorCursorStatus() string {
	return fmt.Sprintf("row:%d --- column:%d", m.editor.Cursor.Row, m.editor.Cursor.Column)
}

func blinkCursor() tea.Cmd {
	return tea.Tick(cursorBlinkInterval, func(t time.Time) tea.Msg {
		return cursorBlinkMsg(t)
	})
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
