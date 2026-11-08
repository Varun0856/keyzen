package ui

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/Varun0856/keyzen/internal/mode"
	"github.com/Varun0856/keyzen/internal/storage"
	"github.com/Varun0856/keyzen/internal/typing"
)

type showMenuMsg struct{}

type showStatsMsg struct{}

type showSettingsMsg struct{}

type startTestMsg struct{ md mode.Mode }

type showResultsMsg struct{ res result }

type settingsChangedMsg struct{ cfg storage.Config }

type result struct {
	md      mode.Mode
	stats   typing.Stats
	best    storage.Best
	hasBest bool
	newBest bool
	weak    []typing.KeyRank
	series  []float64
}

func send(msg tea.Msg) tea.Cmd {
	return func() tea.Msg { return msg }
}
