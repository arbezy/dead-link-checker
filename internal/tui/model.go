package tui

import (
	"log"

	"github.com/arbezy/dead-link-checker/internal/crawling"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// FIX: text inputs take inputs on all screens when they should onlt take input during login screen
// this should be simple enough to fix, just throw an if stmt around textinput handling

const (
	frontView uint = iota
	loginView
	crawlingView
	resultsView
)

type crawlResultMsg []crawling.CheckedLink

type updateProgressBarMsg int

type model struct {
	state      uint
	inputs     []textinput.Model
	focusIndex int
	urllist    []string
	listIndex  int
	progress   progress.Model
	percent    float64
	result     []crawling.CheckedLink
}

func NewModel() model {
	urllist, err := crawling.GetUrls()
	if err != nil {
		log.Fatal("Failure gettings urls")
	}
	m := model{
		state:    frontView,
		inputs:   make([]textinput.Model, 2),
		urllist:  urllist,
		progress: progress.New(progress.WithScaledGradient("#FF7CCB", "#FDFF8C")),
	}
	m.SetTextInputsStyles()

	return m
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmd := m.updateInputs(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		return m.handleKeyInput(key)

	case crawlResultMsg:
		m.result = msg

	case updateProgressBarMsg:
		// while not all links crawled, keep updating progress bar
		if crawling.LinksCrawled < len(m.urllist) {
			m.percent = float64(crawling.LinksCrawled) / float64(len(m.urllist))
			return m, updateProgressBar
		}
		// crawl finished, move onto results
		m.state = resultsView

	}

	return m, cmd
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func startCrawl(urls []string) tea.Cmd {
	return func() tea.Msg {
		res := crawling.CheckLinks(urls)
		return crawlResultMsg(res)
	}
}

func updateProgressBar() tea.Msg {
	return updateProgressBarMsg(0)
}

// TODO: this is starting to get a bit complicated -> could probs move the login logic out of this function
func (m model) handleKeyInput(key string) (tea.Model, tea.Cmd) {
	switch m.state {
	case frontView:
		switch key {
		case "q":
			return m, tea.Quit
		case "l":
			m.state = loginView
		}
	case loginView:
		// TODO: this flow needs a bit of work: shld be able to move up and down, and should chagnge back button
		switch key {
		case "tab":
			// jump from username to password
			m.focusIndex++
			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i < len(m.inputs); i++ {
				for i == m.focusIndex {
					// set focused state
					cmds[i] = m.inputs[i].Focus()
					continue
				}
				// remove focussed states from inputs
				m.inputs[i].Blur()
			}
			return m, tea.Batch(cmds...)

		// TODO: think about making so they need to press enter twice
		case "enter":
			// ask to confirm

			// record input values
			username := m.inputs[0].Value()
			password := m.inputs[1].Value()

			// set proxy
			crawling.SetProxy(username, password)

			// move to crawl
			m.state = crawlingView
			return m, tea.Batch(startCrawl(m.urllist), updateProgressBar)
		case "b":
			m.state = frontView
		}
	case crawlingView:
		// wait until crawl is done or quit early
		switch key {
		case "q":
			// quit early
			m.state = loginView
		}
	case resultsView:
		switch key {
		case "j":
			// showing 30 results at a time, so stop index going above len-30
			if m.listIndex < len(m.result)-30 {
				m.listIndex++
			}
		case "k":
			if m.listIndex > 0 {
				m.listIndex--
			}
		case "q":
			return m, tea.Quit
		}
	}
	return m, nil
}
