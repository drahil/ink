package app

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func autosave(version int) tea.Cmd {
	return tea.Tick(300*time.Millisecond, func(time.Time) tea.Msg {
		return autosaveMsg{version: version}
	})
}

func (m Model) handleAutosave(msg autosaveMsg) (Model, tea.Cmd) {
	if msg.version != m.saveVersion {
		return m, nil
	}

	path := m.editor.Path
	content := m.editor.Content
	if path == "" {
		return m, nil
	}

	return m, func() tea.Msg {
		err := m.project.WriteFile(path, content)
		return saveCompleteMsg{version: msg.version, err: err}
	}
}

func (m Model) handleSaveComplete(msg saveCompleteMsg) (Model, tea.Cmd) {
	if msg.version != m.saveVersion {
		return m, nil
	}

	if msg.err != nil {
		m.status = "save failed: " + msg.err.Error()
	}

	return m, nil
}
