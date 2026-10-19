package ui

import "github.com/charmbracelet/lipgloss"

type Theme struct {
	Name    string
	Accent  lipgloss.Color
	Accent2 lipgloss.Color
	Text    lipgloss.Color
	Dim     lipgloss.Color
	Good    lipgloss.Color
	Bad     lipgloss.Color
	Warn    lipgloss.Color
}

var themes = []Theme{
	{
		Name: "midnight", Accent: "#7aa2f7", Accent2: "#bb9af7", Text: "#c0caf5",
		Dim: "#565f89", Good: "#9ece6a", Bad: "#f7768e", Warn: "#e0af68",
	},
	{
		Name: "nord", Accent: "#88c0d0", Accent2: "#b48ead", Text: "#d8dee9",
		Dim: "#4c566a", Good: "#a3be8c", Bad: "#bf616a", Warn: "#ebcb8b",
	},
	{
		Name: "gruvbox", Accent: "#fabd2f", Accent2: "#d3869b", Text: "#ebdbb2",
		Dim: "#7c6f64", Good: "#b8bb26", Bad: "#fb4934", Warn: "#fe8019",
	},
	{
		Name: "dracula", Accent: "#bd93f9", Accent2: "#ff79c6", Text: "#f8f8f2",
		Dim: "#6272a4", Good: "#50fa7b", Bad: "#ff5555", Warn: "#f1fa8c",
	},
	{
		Name: "catppuccin", Accent: "#89b4fa", Accent2: "#cba6f7", Text: "#cdd6f4",
		Dim: "#6c7086", Good: "#a6e3a1", Bad: "#f38ba8", Warn: "#f9e2af",
	},
}

func SetTheme(name string) bool {
	for _, t := range themes {
		if t.Name == name {
			applyTheme(t)
			return true
		}
	}
	return false
}

func ThemeNames() []string {
	names := make([]string, len(themes))
	for i, t := range themes {
		names[i] = t.Name
	}
	return names
}

func cycleTheme() string {
	next := 0
	for i, t := range themes {
		if t.Name == currentTheme.Name {
			next = (i + 1) % len(themes)
			break
		}
	}
	applyTheme(themes[next])
	return themes[next].Name
}
