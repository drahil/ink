package app

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	m.normalizeFocus()

	if m.prompt.active() {
		return m, m.handlePromptKey(msg)
	}

	switch msg.String() {
	case "ctrl+c", "ctrl+q":
		return m, m.requestQuit()
	}

	// Search is a modal input: its keys must not trigger global shortcuts or
	// mutate the document underneath it.
	if m.focused == PaneEditor && m.editor.IsSearchActive {
		return m, m.handleEditorKey(msg)
	}

	switch msg.String() {
	case "alt+1":
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
	if m.editor.IsSearchActive {
		switch msg.String() {
		case "esc":
			m.editor.CloseSearch()
			m.status = m.editorCursorStatus()
		case "enter", "ctrl+f":
			if m.editor.Search() {
				m.status = m.editorCursorStatus()
			} else {
				m.status = fmt.Sprintf("no matches for %q", m.editor.SearchQuery)
			}
		case "backspace":
			m.editor.BackspaceInSearch()
			m.status = "find: " + m.editor.SearchQuery
		default:
			for _, r := range msg.Runes {
				m.editor.AddToSearchQuery(r)
			}
			m.status = "find: " + m.editor.SearchQuery
		}
		m.cursorVisible = true
		return nil
	}

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
	case "tab":
		m.editor.Tab()
	case "enter":
		m.editor.InsertNewline()
	case "ctrl+f":
		m.editor.InitSearch()
		m.status = "find: " + m.editor.SearchQuery
	case "ctrl+s":
		return m.save()
	default:
		if len(msg.Runes) == 0 {
			m.status = fmt.Sprintf("pressed %q", msg.String())
			return nil
		}

		m.editor.InsertRune(msg.Runes[0])
	}

	m.cursorVisible = true
	m.status = m.editorCursorStatus()

	if before == m.editor.Content {
		return nil
	}

	m.editor.IsModified = true
	m.saveVersion++
	if !m.autosave {
		return nil
	}

	return autosave(m.saveVersion)
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
		row, ok := m.files.SelectedRow()
		if !ok {
			m.status = "no file selected"
			return
		}
		if row.IsDir {
			m.files.ToggleExpanded(row)
			m.status = "toggled directory: " + row.Path
			return
		}

		content, err := m.project.ReadFile(row.Path)
		if err != nil {
			m.status = "could not open: " + row.Path
			return
		}

		m.editor.OpenContent(row.Path, content)
		m.focused = PaneEditor
		m.cursorVisible = true
		m.status = m.editorCursorStatus()
		return
	default:
		if len(msg.Runes) == 0 {
			m.status = fmt.Sprintf("pressed %q", msg.String())
			return
		}

		m.files.InsertSearchRune(msg.Runes[0])
	}
}
