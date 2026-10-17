package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

type model struct {
	target string
	typed  []rune
}

var (
	correctStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	wrongStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Underline(true)
	pendingStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	cursorStyle  = lipgloss.NewStyle().Reverse(true)
)

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc, tea.KeyCtrlC:
			return m, tea.Quit

		case tea.KeyBackspace:
			if len(m.typed) > 0 {
				m.typed = m.typed[:len(m.typed)-1]
			}

		case tea.KeyRunes, tea.KeySpace:
			if len(m.typed) < len([]rune(m.target)) {
				if msg.Type == tea.KeySpace {
					m.typed = append(m.typed, ' ')
				} else {
					m.typed = append(m.typed, msg.Runes...)
				}
			}
		}
	}
	return m, nil

}

func (m model) View() string {
	var b strings.Builder
	for i, ch := range []rune(m.target) {
		s := string(ch)
		switch {
		case i < len(m.typed):
			if m.typed[i] == ch {
				b.WriteString(correctStyle.Render(s))
			} else {
				b.WriteString(wrongStyle.Render(s))
			}
		case i == len(m.typed):
			b.WriteString(cursorStyle.Render(s))
		default:
			b.WriteString(pendingStyle.Render(s))
		}
	}
	return "\n  " + b.String() + "\n\n  esc to quit\n"
}

func main() {
	p := tea.NewProgram(model{target: "the quick brown fox"})
	if _, err := p.Run(); err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

}
