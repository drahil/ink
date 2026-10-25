package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/drahil/ink/internal/app"
)

var version = "dev"

func main() {
	versionFlag := flag.Bool("v", false, "print version and exit")
	versionFlagLong := flag.Bool("version", false, "print version and exit")
	autosaveFlag := flag.Bool("autosave", false, "save changes automatically")
	flag.Parse()

	if *versionFlag || *versionFlagLong {
		fmt.Printf("ink %s\n", version)
		return
	}

	model := app.NewModel(*autosaveFlag)

	program := tea.NewProgram(
		model,
		tea.WithAltScreen(),
	)

	if _, err := program.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
