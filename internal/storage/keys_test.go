package storage

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/Varun0856/keyzen/internal/typing"
)

func TestKeyStoreLoadMissing(t *testing.T) {
	ks := NewKeyStore(filepath.Join(t.TempDir(), "keys.json"))
	got, err := ks.Load()
	if err != nil || len(got) != 0 {
		t.Fatalf("Load on missing file = %v, %v; want empty, nil", got, err)
	}
}

func TestKeyStoreRecordAccumulates(t *testing.T) {
	ks := NewKeyStore(filepath.Join(t.TempDir(), "keys.json"))
	ms := time.Millisecond

	first := map[rune]typing.KeyStat{
		'a': {Presses: 2, Mistakes: 1, Total: 200 * ms, Timed: 1},
	}
	second := map[rune]typing.KeyStat{
		'a': {Presses: 3, Total: 300 * ms, Timed: 3},
		'b': {Presses: 1},
	}
	for _, s := range []map[rune]typing.KeyStat{first, second} {
		if err := ks.Record(s); err != nil {
			t.Fatal(err)
		}
	}

	got, err := ks.Load()
	if err != nil {
		t.Fatal(err)
	}
	a := got["a"]
	if a.Presses != 5 || a.Mistakes != 1 || a.Timed != 4 || a.Total != 500*ms {
		t.Errorf("a = %+v, want 5 presses, 1 mistake, 4 timed, 500ms", a)
	}
	if got["b"].Presses != 1 {
		t.Errorf("b = %+v, want 1 press", got["b"])
	}
}

func TestRunesSkipsMultiCharKeys(t *testing.T) {
	totals := KeyTotals{"a": {Presses: 1}, "bc": {Presses: 9}}
	got := totals.Runes()
	if len(got) != 1 || got['a'].Presses != 1 {
		t.Fatalf("Runes() = %v, want only 'a'", got)
	}
}
