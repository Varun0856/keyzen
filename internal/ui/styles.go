package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	correctStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	wrongStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Underline(true)
	pendingStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	cursorStyle   = lipgloss.NewStyle().Reverse(true)
	timerStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("3")).Bold(true)
	bigStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Bold(true)
	dimStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	goldStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("3")).Bold(true)
	titleStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Bold(true)
	groupStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("5")).Bold(true)
	selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Bold(true)
)

func place(width, height, contentWidth int, content string) string {
	left := (width - contentWidth) / 2
	if left < 0 {
		left = 0
	}
	pad := strings.Repeat(" ", left)

	lines := strings.Split(content, "\n")
	for i, l := range lines {
		if l != "" {
			lines[i] = pad + l
		}
	}
	return strings.Repeat("\n", height/4) + strings.Join(lines, "\n")
}
