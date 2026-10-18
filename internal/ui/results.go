package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
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

	var b strings.Builder
	b.WriteString(dimStyle.Render(m.res.md.Title()) + "\n\n")
	fmt.Fprintf(&b, "%s  %s\n", dimStyle.Render("wpm     "), bigStyle.Render(fmt.Sprintf("%.0f", st.WPM)))
	fmt.Fprintf(&b, "%s  %s\n", dimStyle.Render("accuracy"), bigStyle.Render(fmt.Sprintf("%.1f%%", st.Accuracy)))
	fmt.Fprintf(&b, "%s  %d\n", dimStyle.Render("correct "), st.Correct)
	fmt.Fprintf(&b, "%s  %.1fs\n", dimStyle.Render("time    "), st.Elapsed.Seconds())

	switch {
	case m.res.hasBest && m.res.newBest:
		b.WriteString("\n" + goldStyle.Render("new best!") + "\n")
	case m.res.hasBest:
		b.WriteString("\n" + dimStyle.Render(fmt.Sprintf("best: %.0f wpm", m.res.best.WPM)) + "\n")
	}

	b.WriteString("\n" + dimStyle.Render("tab retry • m menu • q quit") + "\n")
	return place(m.width, m.height, 30, b.String())
}
