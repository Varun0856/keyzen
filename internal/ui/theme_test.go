package ui

import "testing"

func TestSetTheme(t *testing.T) {
	defer applyTheme(themes[0]) // put the default back for other tests

	if !SetTheme("nord") || currentTheme.Name != "nord" {
		t.Fatalf("SetTheme(nord) failed, current = %q", currentTheme.Name)
	}
	if SetTheme("does-not-exist") {
		t.Fatal("SetTheme accepted an unknown theme")
	}
	if currentTheme.Name != "nord" {
		t.Fatalf("an unknown theme must not change the current one, got %q", currentTheme.Name)
	}
}

func TestThemesAreComplete(t *testing.T) {
	seen := map[string]bool{}
	for _, th := range themes {
		if seen[th.Name] {
			t.Errorf("duplicate theme name %q", th.Name)
		}
		seen[th.Name] = true

		for field, c := range map[string]string{
			"Accent": string(th.Accent), "Accent2": string(th.Accent2),
			"Text": string(th.Text), "Dim": string(th.Dim),
			"Good": string(th.Good), "Bad": string(th.Bad), "Warn": string(th.Warn),
		} {
			if c == "" {
				t.Errorf("theme %q has an empty %s", th.Name, field)
			}
		}
	}
}

func TestCycleTheme(t *testing.T) {
	applyTheme(themes[0])
	defer applyTheme(themes[0])

	seen := map[string]bool{}
	for range themes {
		seen[cycleTheme()] = true
	}
	if len(seen) != len(themes) {
		t.Errorf("cycling visited %d themes, want %d", len(seen), len(themes))
	}
	if currentTheme.Name != themes[0].Name {
		t.Errorf("a full cycle should end where it started, got %q", currentTheme.Name)
	}
}

func TestThemeNames(t *testing.T) {
	names := ThemeNames()
	if len(names) != len(themes) || names[0] != themes[0].Name {
		t.Fatalf("ThemeNames() = %v", names)
	}
}
