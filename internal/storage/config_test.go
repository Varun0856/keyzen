package storage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigLoadMissing(t *testing.T) {
	cs := NewConfigStore(filepath.Join(t.TempDir(), "config.json"))
	c, err := cs.Load()
	if err != nil || c.Theme != "" {
		t.Fatalf("Load on missing file = %+v, %v; want zero Config, nil", c, err)
	}
}

func TestConfigRoundTrip(t *testing.T) {
	cs := NewConfigStore(filepath.Join(t.TempDir(), "keyzen", "config.json"))

	if err := cs.Save(Config{Theme: "nord"}); err != nil {
		t.Fatal(err)
	}
	got, err := cs.Load()
	if err != nil {
		t.Fatal(err)
	}
	if got.Theme != "nord" {
		t.Fatalf("Theme = %q, want nord", got.Theme)
	}
}

func TestConfigLoadCorrupt(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte("{oops"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := NewConfigStore(path).Load(); err == nil {
		t.Fatal("expected an error for corrupt JSON")
	}
}
