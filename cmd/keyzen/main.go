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
)

type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

type model struct {
	session   *typing.Session
	remaining int
	finished  bool
}

func newModel() model {
	return model{
		session:   typing.NewSession(words.Random(15)),
		remaining: testSeconds,
	}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyEsc || msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
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
			return m, tick()
		}

	case tickMsg:
		if m.finished {
			return m, nil
		}
		m.remaining--
		if m.remaining <= 0 {
			m.finished = true
			return m, nil
		}
		return m, tick()
	}
	return m, nil

}

func (m model) View() string {
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
	if m.finished {
		status = "time's up (results screen coming next)"
	}
	return "\n " + timerStyle.Render(status) + "\n\n " + b.String() + "\n\n esc to quit\n"
}

func main() {
	if _, err := tea.NewProgram(newModel()).Run(); err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

}
