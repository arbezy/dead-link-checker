package main

import (
	"log"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/progress"
	tea  "github.com/charmbracelet/bubbletea"
)

// TODO: Fix bug with progress bar that causes it to jump in greater increments than 10 after leaving then returning to crawling page

const (
	frontView uint = iota
	loginView
	crawlingView
	resultsView
)

type tickMsg time.Time

// TODO: Need to make status a bit slicker, i.e. using an enum or something...
type Result struct {
	shortname string
	status string 
}

type model struct {
	state	  uint
	// table     table.Model
	textinput textinput.Model
	urllist   []string
	results   []Result
	listIndex int
	progress  progress.Model
	percent   float64
}

func NewModel() model {
	urllist, err := GetUrls()
	if err != nil {
		log.Fatal("Failure gettings urls")
	}
	return model {
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

	case tickMsg:
		m.percent += 0.1
		if m.percent > 1.0 {
			m.percent = 1.0
			m.state = resultsView
		}
		return m, tickCmd()
	}
	
	return m, nil
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
			m.percent = 0.0
			m.state = crawlingView
			return m, tickCmd()
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
			if m.listIndex < len(m.results) { 
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

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
