package main

import (
	"github.com/arbezy/dead-link-checker/internal/tui"
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	m := tui.NewModel()
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Fatalf("unable to run TUI: %v", err)
	}
}
