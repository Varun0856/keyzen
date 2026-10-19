package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type resultsScreen struct {
	res           result
	width, height int
}

func newResults(res result, w, h int) resultsScreen {
	return resultsScreen{res: res, width: w, height: h}
}

func (m resultsScreen) Init() tea.Cmd { return nil }

func (m resultsScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "enter":
			return m, send(startTestMsg{md: m.res.md})
		case "esc", "m":
			return m, send(showMenuMsg{})
		case "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m resultsScreen) View() string {
	st := m.res.stats

	scale := 1
	if m.width >= 100 && m.height >= 30 {
		scale = 2
	}
	hero := renderBig(fmt.Sprintf("%.0f", st.WPM), scale, currentTheme.Accent, currentTheme.Accent2) +
		"\n" + dimStyle.Render("words per minute")

	side := []string{
		bigStyle.Render(fmt.Sprintf("%.1f%%", st.Accuracy)) + dimStyle.Render("  accuracy"),
		textStyle.Render(fmt.Sprintf("%d", st.Correct)) + dimStyle.Render("  correct characters"),
		textStyle.Render(fmt.Sprintf("%.1fs", st.Elapsed.Seconds())) + dimStyle.Render("  time"),
	}
	switch {
	case m.res.hasBest && m.res.newBest:
		side = append(side, goldStyle.Render("★ new best"))
	case m.res.hasBest:
		side = append(side, dimStyle.Render(fmt.Sprintf("best  %.0f wpm", m.res.best.WPM)))
	}
	top := lipgloss.JoinHorizontal(lipgloss.Center, hero, "        ", strings.Join(side, "\n\n"))

	var b strings.Builder
	b.WriteString(dimStyle.Render(m.res.md.Title()) + "\n\n")
	b.WriteString(top)

	if line := sparkline(m.res.series); line != "" {
		b.WriteString("\n\n" + accentStyle.Render(stretch(line, 2)) + "\n" + dimStyle.Render("speed over time"))
	}

	if len(m.res.weak) > 0 {
		chips := make([]string, 0, len(m.res.weak))
		for _, k := range m.res.weak {
			chips = append(chips, bigStyle.Render(string(k.Key))+dimStyle.Render(
				fmt.Sprintf(" %.0f%% %dms", (1-k.Stat.ErrorRate())*100, k.Stat.AvgInterval().Milliseconds())))
		}
		b.WriteString("\n\n" + groupStyle.Render("weakest keys") + "\n" + strings.Join(chips, "    "))
	}

	body := b.String()
	foot := hints("tab", "retry", "m", "menu", "t", "theme", "q", "quit")
	return screen(m.width, m.height, panel(body, lipgloss.Width(body)), foot)
}
