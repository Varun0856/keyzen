package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Varun0856/keyzen/internal/mode"
	"github.com/Varun0856/keyzen/internal/storage"
)

type menuGroup struct {
	name  string
	modes []mode.Mode
}

func groupModes(ms []mode.Mode) []menuGroup {
	var groups []menuGroup
	index := map[string]int{}
	for _, md := range ms {
		if i, ok := index[md.Group]; ok {
			groups[i].modes = append(groups[i].modes, md)
			continue
		}
		index[md.Group] = len(groups)
		groups = append(groups, menuGroup{name: md.Group, modes: []mode.Mode{md}})
	}
	return groups
}

type menu struct {
	groups        []menuGroup
	best          map[string]float64
	suffix        string
	col, row      int
	width, height int
}

func newMenu(store *storage.Store, suffix string, w, h int) menu {
	best := map[string]float64{}
	if store != nil {
		if sc, err := store.Load(); err == nil {
			for k, b := range sc {
				best[k] = b.WPM
			}
		}
	}
	return menu{groups: groupModes(mode.All()), best: best, suffix: suffix, width: w, height: h}
}

func (m menu) Init() tea.Cmd { return nil }

func (m menu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.row > 0 {
				m.row--
			}
		case "down", "j":
			if m.row < len(m.groups[m.col].modes)-1 {
				m.row++
			}
		case "left", "h":
			if m.col > 0 {
				m.col--
				m.row = clamp(m.row, 0, len(m.groups[m.col].modes)-1)
			}
		case "right", "l":
			if m.col < len(m.groups)-1 {
				m.col++
				m.row = clamp(m.row, 0, len(m.groups[m.col].modes)-1)
			}
		case "enter":
			return m, send(startTestMsg{md: m.groups[m.col].modes[m.row]})
		case "s":
			return m, send(showStatsMsg{})
		case "o":
			return m, send(showSettingsMsg{})
		case "q", "esc":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m menu) colWidth() int {
	return clamp((m.width-12)/len(m.groups)-6, 18, 28)
}

func menuLine(label, value string, w int, selected bool) string {
	left := "  " + label
	if selected {
		left = "▸ " + label
	}
	gap := w - lipgloss.Width(left) - lipgloss.Width(value)
	if gap < 1 {
		gap = 1
	}
	if selected {
		return selectedRowStyle.Render(left + strings.Repeat(" ", gap) + value)
	}
	return textStyle.Render(left) + strings.Repeat(" ", gap) + dimStyle.Render(value)
}

func (m menu) column(gi, cw, height int) string {
	g := m.groups[gi]
	active := gi == m.col

	var b strings.Builder
	b.WriteString(groupStyle.Render(g.name) + "\n\n")
	for ri, md := range g.modes {
		best := ""
		if wpm, ok := m.best[md.Key()+m.suffix]; ok {
			best = fmt.Sprintf("%.0f", wpm)
			if cw >= 24 {
				best += " wpm"
			}
		}
		b.WriteString(menuLine(md.Label, best, cw, active && ri == m.row) + "\n")
	}

	border := currentTheme.Dim
	if active {
		border = currentTheme.Accent
	}
	return panelStyle.
		BorderForeground(border).
		Width(cw + 4).
		Height(height).
		Render(strings.TrimRight(b.String(), "\n"))
}

func (m menu) View() string {
	cw := m.colWidth()

	maxRows := 0
	for _, g := range m.groups {
		if len(g.modes) > maxRows {
			maxRows = len(g.modes)
		}
	}
	height := maxRows + 4

	cols := make([]string, 0, len(m.groups)*2)
	for i := range m.groups {
		if i > 0 {
			cols = append(cols, "  ")
		}
		cols = append(cols, m.column(i, cw, height))
	}
	columns := lipgloss.JoinHorizontal(lipgloss.Top, cols...)

	scale := 1
	if m.width >= 100 {
		scale = 2
	}
	logo := renderBig("keyzen", scale, currentTheme.Accent, currentTheme.Accent2)
	if m.height < 24 {
		logo = titleStyle.Render("k e y z e n")
	}

	tagline := dimStyle.Render("type faster, right in your terminal")
	foot := hints("←↑↓→", "select", "enter", "start", "s", "stats", "o", "options", "t", "theme", "q", "quit")

	return screen(m.width, m.height, logo, tagline, columns, foot)
}
