package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	currentTheme Theme

	correctStyle     lipgloss.Style
	wrongStyle       lipgloss.Style
	pendingStyle     lipgloss.Style
	cursorStyle      lipgloss.Style
	timerStyle       lipgloss.Style
	bigStyle         lipgloss.Style
	dimStyle         lipgloss.Style
	goldStyle        lipgloss.Style
	titleStyle       lipgloss.Style
	groupStyle       lipgloss.Style
	selectedStyle    lipgloss.Style
	selectedRowStyle lipgloss.Style
	accentStyle      lipgloss.Style
	textStyle        lipgloss.Style
	keyStyle         lipgloss.Style
	panelStyle       lipgloss.Style
)

func init() { applyTheme(themes[0]) }

func applyTheme(t Theme) {
	currentTheme = t

	correctStyle = lipgloss.NewStyle().Foreground(t.Text)
	wrongStyle = lipgloss.NewStyle().Foreground(t.Bad).Underline(true)
	pendingStyle = lipgloss.NewStyle().Foreground(t.Dim)
	cursorStyle = lipgloss.NewStyle().Background(t.Accent).Foreground(lipgloss.Color("0"))
	timerStyle = lipgloss.NewStyle().Foreground(t.Warn).Bold(true)
	bigStyle = lipgloss.NewStyle().Foreground(t.Accent).Bold(true)
	dimStyle = lipgloss.NewStyle().Foreground(t.Dim)
	goldStyle = lipgloss.NewStyle().Foreground(t.Warn).Bold(true)
	titleStyle = lipgloss.NewStyle().Foreground(t.Accent).Bold(true)
	groupStyle = lipgloss.NewStyle().Foreground(t.Accent2).Bold(true)
	selectedStyle = lipgloss.NewStyle().Foreground(t.Accent).Bold(true)
	selectedRowStyle = lipgloss.NewStyle().Background(t.Accent).Foreground(lipgloss.Color("0")).Bold(true)
	accentStyle = lipgloss.NewStyle().Foreground(t.Accent)
	textStyle = lipgloss.NewStyle().Foreground(t.Text)
	keyStyle = lipgloss.NewStyle().Foreground(t.Accent).Bold(true)
	panelStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Dim).
		Padding(1, 2)
}

const (
	panelExtra    = 6
	panelExtraBig = 10
)

func panel(content string, contentWidth int) string {
	return panelStyle.Width(contentWidth + 4).Render(content)
}

func bigPanel(content string, contentWidth int) string {
	return panelStyle.Padding(2, 4).Width(contentWidth + 8).Render(content)
}

func hints(pairs ...string) string {
	parts := make([]string, 0, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		parts = append(parts, keyStyle.Render(pairs[i])+" "+dimStyle.Render(pairs[i+1]))
	}
	return strings.Join(parts, dimStyle.Render("  •  "))
}

func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func centerBlock(width int, block string) string {
	pad := strings.Repeat(" ", clamp((width-lipgloss.Width(block))/2, 0, width))
	lines := strings.Split(block, "\n")
	for i, l := range lines {
		if l != "" {
			lines[i] = pad + l
		}
	}
	return strings.Join(lines, "\n")
}

func vertical(height int, content string) string {
	lines := strings.Count(content, "\n") + 1
	return strings.Repeat("\n", clamp((height-lines)/3, 0, height)) + content
}

func screen(width, height int, blocks ...string) string {
	centered := make([]string, len(blocks))
	for i, b := range blocks {
		centered[i] = centerBlock(width, b)
	}
	return vertical(height, strings.Join(centered, "\n\n"))
}

func place(width, height, contentWidth int, content string) string {
	pad := strings.Repeat(" ", clamp((width-contentWidth)/2, 0, width))
	lines := strings.Split(content, "\n")
	for i, l := range lines {
		if l != "" {
			lines[i] = pad + l
		}
	}
	return vertical(height, strings.Join(lines, "\n"))
}
