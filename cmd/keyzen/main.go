package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/debug"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/Varun0856/keyzen/internal/storage"
	"github.com/Varun0856/keyzen/internal/ui"
)

var version = "dev"

func buildVersion() string {
	if version != "dev" {
		return version
	}
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
		return info.Main.Version
	}
	return version
}

func main() {
	theme := flag.String("theme", "",
		"color theme: "+strings.Join(ui.ThemeNames(), ", ")+" (default: your last choice)")
	showVersion := flag.Bool("version", false, "print the version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Println("keyzen", buildVersion())
		return
	}

	var cfg *storage.ConfigStore
	if path, err := storage.DefaultConfigPath(); err == nil {
		cfg = storage.NewConfigStore(path)
	}

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
