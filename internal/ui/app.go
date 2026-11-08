package ui

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/Varun0856/keyzen/internal/storage"
	"github.com/Varun0856/keyzen/internal/words"
)

type App struct {
	store         *storage.Store
	keys          *storage.KeyStore
	cfgStore      *storage.ConfigStore
	settings      storage.Config
	current       tea.Model
	width, height int
	gen           int
}

func New(store *storage.Store, keys *storage.KeyStore, cfg *storage.ConfigStore) App {
	a := App{store: store, keys: keys, cfgStore: cfg, width: 80, height: 24}
	if cfg != nil {
		if c, err := cfg.Load(); err == nil {
			a.settings = c
		}
	}
	a.current = newMenu(store, a.options().Suffix(), a.width, a.height)
	return a
}

func (a App) options() words.Options {
	return words.Options{Punctuation: a.settings.Punctuation, Numbers: a.settings.Numbers}
}

func (a App) saveSettings() {
	if a.cfgStore != nil {
		_ = a.cfgStore.Save(a.settings)
	}
}

func (a App) Init() tea.Cmd { return a.current.Init() }

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width, a.height = msg.Width, msg.Height

	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return a, tea.Quit
		}
		if msg.String() == "t" {
			if _, inTest := a.current.(testScreen); !inTest {
				a.settings.Theme = cycleTheme()
				a.saveSettings()
				return a, nil
			}
		}

	case showMenuMsg:
		a.current = newMenu(a.store, a.options().Suffix(), a.width, a.height)
		return a, a.current.Init()

	case startTestMsg:
		a.gen++
		a.current = newTest(msg.md, a.store, a.keys, a.options(), a.gen, a.width, a.height)
		return a, a.current.Init()

	case showResultsMsg:
		a.current = newResults(msg.res, a.width, a.height)
		return a, a.current.Init()

	case showStatsMsg:
		a.current = newStats(a.keys, a.width, a.height)
		return a, a.current.Init()

	case showSettingsMsg:
		a.current = newSettings(a.settings, a.width, a.height)
		return a, a.current.Init()

	case settingsChangedMsg:
		a.settings.Punctuation = msg.cfg.Punctuation
		a.settings.Numbers = msg.cfg.Numbers
		a.saveSettings()
		return a, nil
	}

	var cmd tea.Cmd
	a.current, cmd = a.current.Update(msg)
	return a, cmd
}

func (a App) View() string { return a.current.View() }
