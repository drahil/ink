package app

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/drahil/ink/internal/project"
	"github.com/drahil/ink/internal/ui"
)

const (
	minEditorWidth = 20
	minFilesWidth  = 32
)

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
	project       project.Store
}

func NewModel() Model {
	store := project.NewStore(".")
	items, err := store.ListFiles()
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
		project:       store,
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
		return m.handleAutosave(msg)
	case saveCompleteMsg:
		return m.handleSaveComplete(msg)
	}

	return m, nil
}

func (m Model) editorCursorStatus() string {
	return fmt.Sprintf("row:%d --- column:%d", m.editor.Cursor.Row, m.editor.Cursor.Column)
}
