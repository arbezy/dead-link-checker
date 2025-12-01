package tui

import (
	"log"
	"time"

	"github.com/arbezy/dead-link-checker/internal/crawling"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// TODO: Fix bug with progress bar that causes it to jump in greater increments than 10 after leaving then returning to crawling page
// this might just get fixed by having it track real crawling progress rather than a 10 second timer
// TODO: might be a good idea to seperate model per 'state' or at least create a seperate state for the crawling state

const (
	frontView uint = iota
	loginView
	crawlingView
	resultsView
)

type tickMsg time.Time

type crawlResultMsg []crawling.CheckedLink

type model struct {
	state uint
	// table     table.Model
	textinput textinput.Model
	urllist   []string
	// listIndex int
	progress progress.Model
	percent  float64
	result   []crawling.CheckedLink
}

func NewModel() model {
	urllist, err := crawling.GetUrls()
	if err != nil {
		log.Fatal("Failure gettings urls")
	}
	return model{
		state:     frontView,
		textinput: textinput.New(),
		urllist:   urllist,
		progress:  progress.New(progress.WithScaledGradient("#FF7CCB", "#FDFF8C")),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.textinput, _ = m.textinput.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		return m.handleKeyInput(key)

	case crawlResultMsg:
		m.result = msg
		m.state = resultsView

	}
	return m, nil
}

func startCrawl(urls []string) tea.Cmd {
	return func() tea.Msg {
		res := crawling.CheckLinks(urls)
		return crawlResultMsg(res)
	}
}

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
		switch key {
		case "enter":
			// ask to confirm

			// move to crawl
			m.state = crawlingView
			return m, startCrawl(m.urllist)
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
		//		case "j":
		//			if m.listIndex < len(m.results) {
		//				m.listIndex++
		//			}
		//		case "k":
		//			if m.listIndex > 0 {
		//				m.listIndex--
		//			}
		case "q":
			return m, tea.Quit
		}
	}
	return m, nil
}
