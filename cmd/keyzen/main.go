package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Varun0856/keyzen/internal/typing"
	"github.com/Varun0856/keyzen/internal/words"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

const testSeconds = 30

var (
	correctStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	wrongStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Underline(true)
	pendingStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	cursorStyle  = lipgloss.NewStyle().Reverse(true)
	timerStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("3")).Bold(true)
	bigStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Bold(true)
	dimStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
)

type tickMsg struct{ gen int }

func tick(gen int) tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg{gen: gen}
	})
}

type model struct {
	session   *typing.Session
	remaining int
	finished  bool
	gen       int
}

func newModel(gen int) model {
	return model{
		session:   typing.NewSession(words.Random(15)),
		remaining: testSeconds,
		gen:       gen,
	}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc, tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyTab:
			return newModel(m.gen + 1), nil
		}
		if m.finished {
			return m, nil
		}

		wasStarted := m.session.Started()

		switch msg.Type {
		case tea.KeyBackspace:
			m.session.Backspace()
		case tea.KeySpace:
			m.session.Type(' ')
		case tea.KeyRunes:
			for _, r := range msg.Runes {
				m.session.Type(r)
			}
		}

		if m.session.Done() {
			m.finished = true
			return m, nil
		}
		if !wasStarted && m.session.Started() {
			return m, tick(m.gen)
		}

	case tickMsg:
		if msg.gen != m.gen || m.finished {
			return m, nil
		}
		m.remaining--
		if m.remaining <= 0 {
			m.session.Finish()
			m.finished = true
			return m, nil
		}
		return m, tick(m.gen)
	}
	return m, nil
}

func (m model) View() string {
	if m.finished {
		return m.resultsView()
	}
	var b strings.Builder
	for i := 0; i < m.session.Len(); i++ {
		s := string(m.session.Char(i))
		switch m.session.State(i) {
		case typing.Correct:
			b.WriteString(correctStyle.Render(s))
		case typing.Incorrect:
			b.WriteString(wrongStyle.Render(s))
		default:
			if i == m.session.Cursor() {
				b.WriteString(cursorStyle.Render(s))
			} else {
				b.WriteString(pendingStyle.Render(s))
			}
		}
	}

	status := "start typing to begin"
	if m.session.Started() {
		status = fmt.Sprintf("%ds", m.remaining)
	}
	return "\n " + timerStyle.Render(status) + "\n\n " + b.String() + "\n\n " + dimStyle.Render("tab restart • esc quit") + "\n"
}

func (m model) resultsView() string {
	st := m.session.Stats()
	return fmt.Sprintf(
		"\n  %s  %s\n  %s  %s\n  %s  %d\n  %s  %.1fs\n\n  %s\n",
		dimStyle.Render("wpm     "), bigStyle.Render(fmt.Sprintf("%.0f", st.WPM)),
		dimStyle.Render("accuracy"), bigStyle.Render(fmt.Sprintf("%.1f%%", st.Accuracy)),
		dimStyle.Render("correct "), st.Correct,
		dimStyle.Render("time    "), st.Elapsed.Seconds(),
		dimStyle.Render("tab restart • esc quit"),
	)
}

func main() {
	if _, err := tea.NewProgram(newModel(0)).Run(); err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

}
