package app

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/drahil/ink/internal/ui"
)

func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return "starting..."
	}

	header := ui.RenderHeader(m.width)
	status := ui.RenderStatusBar(m.width, m.unsavedCount())
	if m.prompt.active() {
		status = ui.RenderPrompt(m.width, m.prompt.message)
	}

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

	mainHeight := max(0, m.height-fixedHeight)

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
		filesPane := m.files.View(filesWidth, mainHeight, m.focused == PaneFiles, m.cursorVisible, m.modifiedPaths())

		return lipgloss.JoinHorizontal(lipgloss.Top, editorPane, filesPane)
	}

	return m.editor.View(max(0, m.width), mainHeight, m.focused == PaneEditor, m.cursorVisible)
}

func (m Model) modifiedPaths() map[string]bool {
	modified := make(map[string]bool)
	if m.editor.IsModified && m.editor.Path != "" {
		modified[m.editor.Path] = true
	}
	for path, pane := range m.allOpenedFiles {
		if pane.IsModified {
			modified[path] = true
		}
	}

	return modified
}

func (m Model) unsavedCount() int {
	count := len(m.modifiedPaths())
	if m.editor.IsModified && m.editor.Path == "" {
		count++
	}

	return count
}
