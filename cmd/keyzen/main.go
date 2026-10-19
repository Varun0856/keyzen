package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/Varun0856/keyzen/internal/storage"
	"github.com/Varun0856/keyzen/internal/ui"
)

func main() {
	theme := flag.String("theme", "",
		"color theme: "+strings.Join(ui.ThemeNames(), ", ")+" (default: your last choice)")
	flag.Parse()

	var cfg *storage.ConfigStore
	if path, err := storage.DefaultConfigPath(); err == nil {
		cfg = storage.NewConfigStore(path)
	}

	// An explicit -theme wins; otherwise use the one saved by the t key.
	if *theme != "" {
		if !ui.SetTheme(*theme) {
			fmt.Fprintf(os.Stderr, "keyzen: unknown theme %q\n", *theme)
			os.Exit(2)
		}
	} else if cfg != nil {
		if c, err := cfg.Load(); err == nil && c.Theme != "" {
			ui.SetTheme(c.Theme) // an unknown saved name is simply ignored
		}
	}

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

	p := tea.NewProgram(ui.New(store, keys, cfg), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
