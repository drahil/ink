package ui

import "github.com/charmbracelet/x/ansi"

func truncateCells(value string, maxWidth int) string {
	return truncateCellsWithTail(value, maxWidth, "")
}

func truncateCellsWithTail(value string, maxWidth int, tail string) string {
	if maxWidth <= 0 {
		return ""
	}

	return ansi.Truncate(value, maxWidth, tail)
}
