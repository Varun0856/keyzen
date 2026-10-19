package ui

import (
	"fmt"
	"strings"

	"github.com/Varun0856/keyzen/internal/mode"
	"github.com/Varun0856/keyzen/internal/storage"
	tea "github.com/charmbracelet/bubbletea"
)

type menu struct {
	modes         []mode.Mode
	best          map[string]float64
	cursor        int
	width, height int
}

func newMenu(store *storage.Store, w, h int) menu {
	best := map[string]float64{}
	if store != nil {
		if sc, err := store.Load(); err == nil {
			for k, b := range sc {
				best[k] = b.WPM
			}
		}
	}
	return menu{modes: mode.All(), best: best, width: w, height: h}
}

func (m menu) Init() tea.Cmd { return nil }

func (m menu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.modes)-1 {
				m.cursor++
			}
		case "enter":
			return m, send(startTestMsg{md: m.modes[m.cursor]})
		case "s":
			return m, send(showStatsMsg{})
		case "q", "esc":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m menu) View() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("keyzen") + "\n")
	b.WriteString(dimStyle.Render("a typing test for your terminal") + "\n\n")

	prevGroup := ""
	for i, md := range m.modes {
		if md.Group != prevGroup {
			if prevGroup != "" {
				b.WriteString("\n")
			}
			b.WriteString(groupStyle.Render(md.Group) + "\n")
			prevGroup = md.Group
		}

		marker := " "
		label := fmt.Sprintf("%-10s", md.Label)
		if i == m.cursor {
			marker = "▸ "
			label = selectedStyle.Render(label)
		}

		best := ""
		if wpm, ok := m.best[md.Key()]; ok {
			best = dimStyle.Render(fmt.Sprintf("best %.0f wpm", wpm))
		}
		b.WriteString(" " + marker + label + best + "\n")
	}

	b.WriteString("\n" + dimStyle.Render("↑/↓ select • enter start • s stats • q quit") + "\n")
	return place(m.width, m.height, 44, b.String())
}
