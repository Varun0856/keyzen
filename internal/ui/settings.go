package ui

import (
	"math/rand"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/Varun0856/keyzen/internal/storage"
	"github.com/Varun0856/keyzen/internal/words"
)

const previewSample = "the quick brown fox jumps over the lazy dog"

type settingsScreen struct {
	cfg           storage.Config
	cursor        int
	width, height int
}

func newSettings(cfg storage.Config, w, h int) settingsScreen {
	return settingsScreen{cfg: cfg, width: w, height: h}
}

func (m settingsScreen) Init() tea.Cmd { return nil }

func (m settingsScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			m.cursor = clamp(m.cursor-1, 0, 1)
		case "down", "j":
			m.cursor = clamp(m.cursor+1, 0, 1)
		case "enter", " ", "left", "right", "h", "l":
			if m.cursor == 0 {
				m.cfg.Punctuation = !m.cfg.Punctuation
			} else {
				m.cfg.Numbers = !m.cfg.Numbers
			}
			return m, send(settingsChangedMsg{cfg: m.cfg})
		case "esc", "m":
			return m, send(showMenuMsg{})
		case "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m settingsScreen) options() words.Options {
	return words.Options{Punctuation: m.cfg.Punctuation, Numbers: m.cfg.Numbers}
}

func (m settingsScreen) preview() string {
	rng := rand.New(rand.NewSource(7))
	return words.Decorate(previewSample, m.options(), rng.Float64)
}

func onOff(on bool) string {
	if on {
		return "● on"
	}
	return "○ off"
}

func (m settingsScreen) View() string {
	const cw = 50

	var b strings.Builder
	b.WriteString(titleStyle.Render("settings") + "\n\n")
	b.WriteString(menuLine("Punctuation", onOff(m.cfg.Punctuation), cw, m.cursor == 0) + "\n")
	b.WriteString(menuLine("Numbers", onOff(m.cfg.Numbers), cw, m.cursor == 1) + "\n")
	b.WriteString("\n" + dimStyle.Render("theme  ") + accentStyle.Render(currentTheme.Name) +
		dimStyle.Render("   (press t to change)") + "\n\n")
	b.WriteString(groupStyle.Render("preview") + "\n" + textStyle.Render(m.preview()))

	foot := hints("↑↓", "move", "enter", "toggle", "t", "theme", "esc", "back")
	return screen(m.width, m.height, panel(b.String(), cw), foot)
}
