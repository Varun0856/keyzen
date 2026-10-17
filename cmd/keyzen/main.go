package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	target string
	typed  []rune
}

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
	return m.target + "\n" + string(m.typed) + "\n\nesc to quit\n"
}

func main() {
	p := tea.NewProgram(model{target: "the quick brown fox"})
	if _, err := p.Run(); err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

}
