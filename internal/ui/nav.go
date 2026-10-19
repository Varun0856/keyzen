package ui

import (
	"github.com/Varun0856/keyzen/internal/mode"
	"github.com/Varun0856/keyzen/internal/storage"
	"github.com/Varun0856/keyzen/internal/typing"
	tea "github.com/charmbracelet/bubbletea"
)

type showMenuMsg struct{}

type startTestMsg struct{ md mode.Mode }

type showResultsMsg struct{ res result }

type showStatsMsg struct{}

type result struct {
	md      mode.Mode
	stats   typing.Stats
	best    storage.Best
	hasBest bool
	newBest bool
	weak    []typing.KeyRank
}

func send(msg tea.Msg) tea.Cmd {
	return func() tea.Msg { return msg }
}
