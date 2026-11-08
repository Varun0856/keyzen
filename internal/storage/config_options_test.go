package storage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigKeepsOptions(t *testing.T) {
	cs := NewConfigStore(filepath.Join(t.TempDir(), "config.json"))
	want := Config{Theme: "nord", Punctuation: true, Numbers: true}

	if err := cs.Save(want); err != nil {
		t.Fatal(err)
	}
	got, err := cs.Load()
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func TestConfigOldFileWithoutOptions(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte(`{"theme":"nord"}`), 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := NewConfigStore(path).Load()
	if err != nil {
		t.Fatal(err)
	}
	if got.Theme != "nord" || got.Punctuation || got.Numbers {
		t.Fatalf("got %+v, want theme nord with both options off", got)
	}
}
