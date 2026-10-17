package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Varun0856/keyzen/internal/typing"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

type model struct {
	session *typing.Session
}

var (
	correctStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	wrongStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Underline(true)
	pendingStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	cursorStyle  = lipgloss.NewStyle().Reverse(true)
)

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.Type {
		case tea.KeyEsc, tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyBackspace:
			m.session.Backspace()
		case tea.KeySpace:
			m.session.Type(' ')
		case tea.KeyRunes:
			for _, r := range msg.Runes {
				m.session.Type(r)
			}
		}
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
	return "\n " + b.String() + "\n\n  esc to quit\n"
}

func main() {
	m := model{session: typing.NewSession("The quick brown fox jumps over the lazy dog")}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

}
