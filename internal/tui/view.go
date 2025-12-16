package tui

import (
	"fmt"
	"github.com/arbezy/dead-link-checker/internal/crawling"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

var (
	appNameStyle    = lipgloss.NewStyle().Background(lipgloss.Color("99")).Padding(0, 1)
	faint           = lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Faint(true)
	goodStatusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true)
	badStatusStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Bold(true)
)

// text input styles
var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	noStyle      = lipgloss.NewStyle()
	helpStyle    = blurredStyle

	focusedButton = focusedStyle.Render("[ submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

func (m model) View() string {
	s := appNameStyle.Render("DEAD LINK CHECKER")
	s += "\n"

	if m.state == frontView {
		s += "Welcome!\n\n"
		s += faint.Render("(q)uit | (l)ogin")
	}

	if m.state == loginView {
		s += "Enter proxy username and password:\n\n"

		for i := 0; i < len(m.inputs); i++ {
			m.inputs[i].PromptStyle = noStyle
			m.inputs[i].TextStyle = noStyle
			if i == m.focusIndex {
				m.inputs[i].PromptStyle = focusedStyle
				m.inputs[i].TextStyle = focusedStyle
			}
		}

		s += m.viewTextInputs()

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
		// only show 30 results at a time
		top, bottom := m.keepWithinBounds(m.listIndex, m.listIndex+30)

		for i, n := range m.result[top:bottom] {
			shortUrl := m.getShortUrl(n)
			paddedIndex := fmt.Sprintf("%3d. ", m.listIndex+i)
			if n.Status == "200 OK" {
				s += paddedIndex + shortUrl + " | " + goodStatusStyle.Render(n.Status) + "\n"
			} else {
				s += paddedIndex + shortUrl + " | " + badStatusStyle.Render(n.Status) + "\n"
			}
		}
		s += faint.Render("(q)uit | (j/k) to scroll up / down")
	}

	return s
}

func (m model) getShortUrl(res crawling.CheckedLink) string {
	shortUrl := res.Url
	if len(res.Url) > 100 {
		shortUrl = res.Url[:100] + "..."
	} else {
		// pad to 100 (and leave room for elipsis)
		for len(shortUrl) < 103 {
			shortUrl += " "
		}
	}
	return shortUrl
}

func (m model) keepWithinBounds(top, bottom int) (int, int) {
	if top < 0 {
		top = 0
		bottom = 30
	}
	if bottom > len(m.result) {
		top = len(m.result) - 30
		if top < 0 {
			top = 0
		}
		bottom = len(m.result)
	}
	return top, bottom
}

func (m model) SetTextInputsStyles() {
	for i := range m.inputs {
		m.inputs[i] = textinput.New()
		switch i {
		case 0:
			m.inputs[i].Placeholder = "Username"
			m.inputs[i].Focus()
			m.inputs[i].CharLimit = 156
			m.inputs[i].Width = 30
		case 1:
			m.inputs[i].Placeholder = "Password"
			m.inputs[i].EchoMode = textinput.EchoPassword
			m.inputs[i].EchoCharacter = '*'
			m.inputs[i].CharLimit = 156
			m.inputs[i].Width = 30
		}
	}
}

func (m model) viewTextInputs() string {
	var b strings.Builder

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	return b.String()
}
