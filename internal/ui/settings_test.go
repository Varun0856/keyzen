package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/Varun0856/keyzen/internal/storage"
)

func press(m settingsScreen, k tea.KeyType) (settingsScreen, tea.Msg) {
	model, cmd := m.Update(tea.KeyMsg{Type: k})
	var msg tea.Msg
	if cmd != nil {
		msg = cmd()
	}
	return model.(settingsScreen), msg
}

func TestSettingsToggle(t *testing.T) {
	m := newSettings(storage.Config{}, 80, 24)

	m, msg := press(m, tea.KeyEnter)
	if !m.cfg.Punctuation {
		t.Fatal("enter on the first row should turn punctuation on")
	}
	if changed, ok := msg.(settingsChangedMsg); !ok || !changed.cfg.Punctuation {
		t.Fatalf("expected settingsChangedMsg with punctuation on, got %#v", msg)
	}

	m, _ = press(m, tea.KeyDown)
	m, _ = press(m, tea.KeyEnter)
	if !m.cfg.Numbers {
		t.Fatal("enter on the second row should turn numbers on")
	}

	m, _ = press(m, tea.KeyEnter)
	if m.cfg.Numbers {
		t.Fatal("a second enter should turn numbers back off")
	}
}

func TestSettingsCursorStaysInRange(t *testing.T) {
	m := newSettings(storage.Config{}, 80, 24)

	for i := 0; i < 5; i++ {
		m, _ = press(m, tea.KeyDown)
	}
	if m.cursor != 1 {
		t.Fatalf("cursor = %d after pressing down 5 times, want 1", m.cursor)
	}
	for i := 0; i < 5; i++ {
		m, _ = press(m, tea.KeyUp)
	}
	if m.cursor != 0 {
		t.Fatalf("cursor = %d after pressing up 5 times, want 0", m.cursor)
	}
}

func TestSettingsEscGoesToMenu(t *testing.T) {
	m := newSettings(storage.Config{}, 80, 24)
	_, msg := press(m, tea.KeyEsc)
	if _, ok := msg.(showMenuMsg); !ok {
		t.Fatalf("esc should return to the menu, got %#v", msg)
	}
}
