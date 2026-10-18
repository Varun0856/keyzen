package ui

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/Varun0856/keyzen/internal/storage"
)

type App struct {
	store         *storage.Store
	current       tea.Model
	width, height int
	gen           int
}

func New(store *storage.Store) App {
	a := App{store: store, width: 80, height: 24}
	a.current = newMenu(store, a.width, a.height)
	return a
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

	case showMenuMsg:
		a.current = newMenu(a.store, a.width, a.height)
		return a, a.current.Init()

	case startTestMsg:
		a.gen++
		a.current = newTest(msg.md, a.store, a.gen, a.width, a.height)
		return a, a.current.Init()

	case showResultsMsg:
		a.current = newResults(msg.res, a.width, a.height)
		return a, a.current.Init()
	}

	var cmd tea.Cmd
	a.current, cmd = a.current.Update(msg)
	return a, cmd
}

func (a App) View() string { return a.current.View() }
