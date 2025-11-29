package tui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	appNameStyle = lipgloss.NewStyle().Background(lipgloss.Color("99")).Padding(0, 1)
	faint        = lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Faint(true)
	//listEnumeratorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("99")).MarginRight(1)
)

func (m model) View() string {
	s := appNameStyle.Render("URL LINK CHECKER")
	s += "\n"

	if m.state == frontView {
		s += "Welcome!\n\n"
		s += faint.Render("(q)uit | (l)ogin")
	}

	if m.state == loginView {
		s += "Enter proxy username and password:\n\n"
		s += faint.Render("go (b)ack to front menu, or (Enter) to continue")
	}

	if m.state == crawlingView {
		s += "Crawling...\n\n"
		s += m.progress.ViewAs(m.percent)
		s += "\n\n"
		s += faint.Render("(q)uit")
	}

	if m.state == resultsView {
		s += "Crawl Results:\n\n"
		s += faint.Render("(q)uit")
	}

	return s
}
