package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Varun0856/keyzen/internal/typing"
	"github.com/Varun0856/keyzen/internal/words"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

const (
	maxTextWidth = 60
	visibleLines = 3
)

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
	seconds   int
	remaining int
	finished  bool
	gen       int
	width     int
}

func newModel(gen, seconds, width int) model {
	return model{
		session:   typing.NewSession(words.Random(seconds * 3)),
		seconds:   seconds,
		remaining: seconds,
		gen:       gen,
		width:     width,
	}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc, tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyTab:
			return newModel(m.gen+1, m.seconds, m.width), nil
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

func (m model) renderChar(i int) string {
	s := string(m.session.Char(i))
	switch m.session.State(i) {
	case typing.Correct:
		return correctStyle.Render(s)
	case typing.Incorrect:
		return wrongStyle.Render(s)
	}
	if i == m.session.Cursor() {
		return cursorStyle.Render(s)
	}
	return pendingStyle.Render(s)
}

func (m model) textView() string {
	width := m.width - 4
	if width > maxTextWidth {
		width = maxTextWidth
	}
	if width < 10 {
		width = 10
	}

	runes := make([]rune, m.session.Len())
	for i := range runes {
		runes[i] = m.session.Char(i)
	}
	lines := wrapLines(runes, width)

	cur := len(lines) - 1
	for idx, ln := range lines {
		if m.session.Cursor() >= ln[0] && m.session.Cursor() < ln[1] {
			cur = idx
			break
		}
	}

	first := cur - 1
	if first < 0 {
		first = 0
	}
	last := first + visibleLines
	if last > len(lines) {
		last = len(lines)
	}

	var b strings.Builder
	for _, ln := range lines[first:last] {
		b.WriteString(" ")
		for i := ln[0]; i < ln[1]; i++ {
			b.WriteString(m.renderChar(i))
		}
		b.WriteString("\n")
	}
	return b.String()
}

func (m model) View() string {
	if m.finished {
		return m.resultsView()
	}

	status := "start typing to begin"
	if m.session.Started() {
		status = fmt.Sprintf("%ds", m.remaining)
	}
	return "\n " + timerStyle.Render(status) + "\n\n " + m.textView() + "\n " + dimStyle.Render("tab restart • esc quit") + "\n"
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
	seconds := flag.Int("t", 30, "test duration in seconds (15 30 or 60)")
	flag.Parse()

	switch *seconds {
	case 15, 30, 60:
	default:
		fmt.Fprintln(os.Stderr, "keyzen: -t must be 15, 30 or 60")
		os.Exit(2)
	}
	if _, err := tea.NewProgram(newModel(0, *seconds, 80)).Run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

}
