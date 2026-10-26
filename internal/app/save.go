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
	if msg.version != m.saveVersion || m.editor.Path == "" || !m.editor.IsModified {
		return m, nil
	}

	return m, m.save()
}

func (m Model) handleSaveComplete(msg saveCompleteMsg) (Model, tea.Cmd) {
	quitAfterSave := m.quitAfterSave
	m.quitAfterSave = false

	if msg.err != nil {
		m.status = "save failed: " + msg.err.Error()
		return m, nil
	}

	if msg.version == m.saveVersion {
		m.editor.IsModified = false
	}
	m.status = "saved " + msg.path

	if quitAfterSave && !m.editor.IsModified {
		return m, tea.Quit
	}

	return m, nil
}

func (m *Model) save() tea.Cmd {
	if m.editor.Path == "" {
		m.status = "no file name"
		return nil
	}
	if !m.editor.IsModified {
		m.status = "no changes"
		return nil
	}

	store := m.project
	path := m.editor.Path
	content := m.editor.Content
	version := m.saveVersion

	return func() tea.Msg {
		err := store.WriteFile(path, content)
		return saveCompleteMsg{version: version, path: path, err: err}
	}
}
