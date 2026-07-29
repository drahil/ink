package app

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

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

		content, err := m.project.ReadFile(item)
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

		m.files.InsertSearchRune(msg.Runes[0])
	}
}
