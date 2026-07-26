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
)

type cursorBlinkMsg time.Time

type autosaveMsg struct {
	version int
}

type saveCompleteMsg struct {
	version int
	err     error
}

const cursorBlinkInterval = 500 * time.Millisecond
const minEditorWidth = 20
const minFilesWidth = 32

type Model struct {
	width         int
	height        int
	focused       Pane
	status        string
	cursorVisible bool
	files         ui.FilesPane
	editor        ui.EditorPane
	filesVisible  bool
	saveVersion   int
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
		filesVisible:  true,
	}
}

func (p Pane) String() string {
	switch p {
	case PaneFiles:
		return "files"
	case PaneEditor:
		return "editor"
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
		m.normalizeFocus()
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
	case autosaveMsg:
		if msg.version != m.saveVersion {
			return m, nil
		}

		path := m.editor.Path
		content := m.editor.Content
		if path == "" {
			return m, nil
		}

		return m, func() tea.Msg {
			err := project.WriteFile(".", path, content)
			return saveCompleteMsg{version: msg.version, err: err}
		}
	case saveCompleteMsg:
		if msg.version != m.saveVersion {
			return m, nil
		}

		if msg.err != nil {
			m.status = "save failed: " + msg.err.Error()
			return m, nil
		}

		return m, nil
	}

	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	m.normalizeFocus()

	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "tab":
		m.focusNextPane()
		return m, nil
	case "alt+1", "cmd+1":
		m.toggleFilesPane()
		return m, nil
	}

	switch m.focused {
	case PaneEditor:
		return m, m.handleEditorKey(msg)
	case PaneFiles:
		m.handleFilesKey(msg)
	}

	return m, nil
}

func (m *Model) handleEditorKey(msg tea.KeyMsg) tea.Cmd {
	before := m.editor.Content

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
			return nil
		}

		m.editor.InsertRune(msg.Runes[0])
	}

	m.cursorVisible = true
	m.status = m.editorCursorStatus()

	if before == m.editor.Content || m.editor.Path == "" {
		return nil
	}

	m.saveVersion++
	return autosave(m.saveVersion)
}

func autosave(version int) tea.Cmd {
	return tea.Tick(300*time.Millisecond, func(time.Time) tea.Msg {
		return autosaveMsg{version: version}
	})
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
		m.saveVersion = 0
		m.focused = PaneEditor
		m.cursorVisible = true
		m.status = m.editorCursorStatus()
		return
	default:
		if len(msg.Runes) == 0 {
			m.status = fmt.Sprintf("pressed %q", msg.String())
			return
		}

		m.files.Query(msg.Runes[0])
	}
}

func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return "starting..."
	}

	header := ui.RenderHeader(m.width)
	status := ui.RenderStatusBar(m.width, m.status)

	headerHeight := lipgloss.Height(header)
	statusHeight := lipgloss.Height(status)
	includeHeader := m.height >= headerHeight
	includeStatus := m.height >= headerHeight+statusHeight

	fixedHeight := 0
	if includeHeader {
		fixedHeight += headerHeight
	}
	if includeStatus {
		fixedHeight += statusHeight
	}

	remainingHeight := max(0, m.height-fixedHeight)
	mainHeight := remainingHeight

	sections := make([]string, 0, 3)
	if includeHeader {
		sections = append(sections, header)
	}
	if mainHeight > 0 {
		sections = append(sections, m.renderMain(mainHeight))
	}
	if includeStatus {
		sections = append(sections, status)
	}

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (m Model) renderMain(mainHeight int) string {
	if m.filesPaneAvailable() {
		filesWidth := min(48, max(minFilesWidth, m.width/3))
		filesWidth = min(filesWidth, m.width-minEditorWidth)
		editorWidth := max(0, m.width-filesWidth)
		editorPane := m.editor.View(editorWidth, mainHeight, m.focused == PaneEditor, m.cursorVisible)
		filesPane := m.files.View(filesWidth, mainHeight, m.focused == PaneFiles, m.cursorVisible)

		return lipgloss.JoinHorizontal(lipgloss.Top, editorPane, filesPane)
	}

	editorWidth := max(0, m.width)
	return m.editor.View(editorWidth, mainHeight, m.focused == PaneEditor, m.cursorVisible)
}

func (m Model) filesPaneAvailable() bool {
	return m.filesVisible && m.width >= minEditorWidth+minFilesWidth
}

func (m *Model) normalizeFocus() {
	if m.focused == PaneFiles && !m.filesPaneAvailable() {
		m.focused = PaneEditor
		m.cursorVisible = true
	}
}

func (m *Model) focusNextPane() {
	switch m.focused {
	case PaneFiles:
		m.focused = PaneEditor
	case PaneEditor:
		if m.filesPaneAvailable() {
			m.focused = PaneFiles
		} else {
			m.focused = PaneEditor
		}
	}

	m.cursorVisible = m.focused == PaneEditor || m.focused == PaneFiles
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

func (m *Model) toggleFilesPane() {
	m.filesVisible = !m.filesVisible
	if !m.filesVisible {
		m.focused = PaneEditor
		m.cursorVisible = true
	} else if m.filesPaneAvailable() {
		m.focused = PaneFiles
		m.cursorVisible = true
	} else {
		m.focused = PaneEditor
		m.cursorVisible = true
	}
}
