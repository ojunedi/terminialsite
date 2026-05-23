package main

import (
	_ "embed"
	"strings"
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/common-nighthawk/go-figure"
	zone "github.com/lrstanley/bubblezone"
)

//go:embed assets/portrait.txt
var portrait string

// page identifies what's on screen. pageHome is the landing screen; 0..N map
// to entries in the nav slice (see content.go).
const pageHome = -1

// nameBanner is the figlet rendering of the configured name, computed once.
var nameBanner = strings.TrimRight(figure.NewFigure(name, "small", true).String(), "\n")

// model holds the per-session UI state.
type model struct {
	term      string
	width     int
	height    int
	user      string
	page      int // pageHome or 0..len(footerLabels)-1
	footerIdx int // highlighted footer item on the landing screen
	cursor    int // selected row within the current section page
	opened    string
	frame     int // animation frame counter, advanced by tickMsg
}

// tickMsg drives the name-banner shimmer animation.
type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(time.Second/15, func(t time.Time) tea.Msg { return tickMsg(t) })
}

func (m model) Init() tea.Cmd { return tick() }

// selectables returns the openable rows for the current section page.
func (m model) selectables() []item {
	if m.page < 0 || m.page >= len(nav) {
		return nil
	}
	e := nav[m.page]
	if e.contacts {
		out := make([]item, len(contacts))
		for i, c := range contacts {
			out[i] = item{title: c.label + "  " + c.handle, url: c.url}
		}
		return out
	}
	return flatten(*e.section)
}

func flatten(s section) []item {
	var out []item
	for _, g := range s.groups {
		out = append(out, g.items...)
	}
	return out
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		m.frame++
		return m, tick() // re-arm so the animation keeps running

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.MouseMsg:
		if msg.Action != tea.MouseActionRelease || msg.Button != tea.MouseButtonLeft {
			return m, nil
		}
		// Footer nav clicks (available on every screen).
		for i := range nav {
			if zone.Get(footerZoneID(i)).InBounds(msg) {
				m.page = i
				m.footerIdx = i
				m.cursor = 0
				m.opened = ""
				return m, nil
			}
		}
		// Row clicks within a section page.
		if m.page != pageHome {
			rows := m.selectables()
			for i := range rows {
				if zone.Get(rowZoneID(i)).InBounds(msg) {
					m.cursor = i
					m.opened = rows[i].url
					return m, nil
				}
			}
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "esc":
			if m.opened != "" {
				m.opened = ""
			} else if m.page != pageHome {
				m.page = pageHome
			}
		case "left", "h":
			if m.page == pageHome {
				m.footerIdx = (m.footerIdx + len(nav) - 1) % len(nav)
			}
		case "right", "l", "tab":
			if m.page == pageHome {
				m.footerIdx = (m.footerIdx + 1) % len(nav)
			}
		case "up", "k":
			if m.page != pageHome && m.cursor > 0 {
				m.cursor--
				m.opened = ""
			}
		case "down", "j":
			if m.page != pageHome && m.cursor < len(m.selectables())-1 {
				m.cursor++
				m.opened = ""
			}
		case "enter":
			if m.page == pageHome {
				m.page = m.footerIdx
				m.cursor = 0
				m.opened = ""
			} else if rows := m.selectables(); m.cursor < len(rows) {
				m.opened = rows[m.cursor].url
			}
		}
	}
	return m, nil
}

func footerZoneID(i int) string { return "footer-" + nav[i].label }
func rowZoneID(i int) string    { return "row-" + string(rune('a'+i)) }
