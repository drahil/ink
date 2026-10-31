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
	if msg.err != nil {
		m.quitAfterSave = false
		m.status = "save failed: " + msg.err.Error()
		return m, nil
	}

	if msg.version == m.saveVersion {
		m.markSaved(msg.path)
	}
	m.status = "saved " + msg.path

	if m.quitAfterSave && !m.anyModified() {
		return m, tea.Quit
	}

	return m, nil
}

func (m *Model) markSaved(path string) {
	if path == m.editor.Path {
		m.editor.IsModified = false
		return
	}

	pane, ok := m.allOpenedFiles[path]
	if !ok {
		return
	}
	pane.IsModified = false
	m.allOpenedFiles[path] = pane
}

func (m Model) anyModified() bool {
	if m.editor.IsModified {
		return true
	}
	for _, pane := range m.allOpenedFiles {
		if pane.IsModified {
			return true
		}
	}

	return false
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

	return m.writeFile(m.editor.Path, m.editor.Content)
}

func (m *Model) saveAll() tea.Cmd {
	var cmds []tea.Cmd
	if m.editor.IsModified && m.editor.Path != "" {
		cmds = append(cmds, m.writeFile(m.editor.Path, m.editor.Content))
	}
	for path, pane := range m.allOpenedFiles {
		if pane.IsModified {
			cmds = append(cmds, m.writeFile(path, pane.Content))
		}
	}

	if len(cmds) == 0 {
		m.status = "no changes"
		return nil
	}

	return tea.Batch(cmds...)
}

func (m *Model) writeFile(path, content string) tea.Cmd {
	store := m.project
	version := m.saveVersion

	return func() tea.Msg {
		err := store.WriteFile(path, content)
		return saveCompleteMsg{version: version, path: path, err: err}
	}
}
