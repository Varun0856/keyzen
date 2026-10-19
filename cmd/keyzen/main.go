package main

import (
	"fmt"
	"os"

	"github.com/Varun0856/keyzen/internal/storage"
	"github.com/Varun0856/keyzen/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	var store *storage.Store
	if path, err := storage.DefaultPath(); err != nil {
		fmt.Fprintln(os.Stderr, "keyzen: scores won't be saved:", err)
	} else {
		store = storage.New(path)
	}

	var keys *storage.KeyStore
	if path, err := storage.DefaultKeysPath(); err == nil {
		keys = storage.NewKeyStore(path)
	}

	p := tea.NewProgram(ui.New(store, keys), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
