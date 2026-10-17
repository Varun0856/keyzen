package storage

import (
	"os"
	"path/filepath"
	"testing"
)

func newTestStore(t *testing.T) *Store {
	t.Helper()
	return New(filepath.Join(t.TempDir(), "keyzen", "scores.json"))
}

func TestLoadMissingFile(t *testing.T) {
	sc, err := newTestStore(t).Load()
	if err != nil {
		t.Fatalf("Load on missing file returned error: %v", err)
	}
	if len(sc) != 0 {
		t.Fatalf("expected empty scores, got %v", sc)
	}
}

func TestRecord(t *testing.T) {
	st := newTestStore(t)

	steps := []struct {
		name    string
		seconds int
		wpm     float64
		wantNew bool
		wantWPM float64
	}{
		{"first score is a best", 30, 50, true, 50},
		{"lower score is not", 30, 40, false, 50},
		{"equal score is not", 30, 50, false, 50},
		{"higher score is", 30, 60, true, 60},
		{"other mode is separate", 60, 20, true, 20},
	}

	for _, s := range steps {
		t.Run(s.name, func(t *testing.T) {
			best, isNew, err := st.Record(s.seconds, s.wpm, 95)
			if err != nil {
				t.Fatal(err)
			}
			if isNew != s.wantNew || best.WPM != s.wantWPM {
				t.Errorf("got best=%.0f new=%v, want best=%.0f new=%v",
					best.WPM, isNew, s.wantWPM, s.wantNew)
			}
		})
	}
}

func TestPersistsAcrossStores(t *testing.T) {
	path := filepath.Join(t.TempDir(), "scores.json")

	if _, _, err := New(path).Record(15, 70, 98); err != nil {
		t.Fatal(err)
	}

	sc, err := New(path).Load() // a brand new Store, same file
	if err != nil {
		t.Fatal(err)
	}
	if sc["15s"].WPM != 70 {
		t.Fatalf("expected 70 wpm saved for 15s, got %v", sc)
	}
}

func TestLoadCorruptFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "scores.json")
	if err := os.WriteFile(path, []byte("{not json"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := New(path).Load(); err == nil {
		t.Fatal("expected an error for corrupt JSON")
	}
}
