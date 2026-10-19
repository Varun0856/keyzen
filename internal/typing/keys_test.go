package typing

import (
	"testing"
	"time"
)

func TestAnalyzeKeys(t *testing.T) {
	ms := time.Millisecond
	events := []KeyEvent{
		{'a', 'a', 0},
		{'b', 'x', 200 * ms}, // mistake on b
		{'b', 'b', 500 * ms}, // b correct, 300ms after previous key
		{'a', 'a', 600 * ms}, // a correct, 100ms after previous key
	}

	got := AnalyzeKeys(events)

	a := got['a']
	if a.Presses != 2 || a.Mistakes != 0 || a.AvgInterval() != 100*ms {
		t.Errorf("a = %+v (avg %v), want 2 presses, 0 mistakes, 100ms", a, a.AvgInterval())
	}
	b := got['b']
	if b.Presses != 2 || b.Mistakes != 1 || b.AvgInterval() != 300*ms {
		t.Errorf("b = %+v (avg %v), want 2 presses, 1 mistake, 300ms", b, b.AvgInterval())
	}
}

func TestWeakKeys(t *testing.T) {
	stats := map[rune]KeyStat{
		'a': {Presses: 10, Mistakes: 0},
		'b': {Presses: 10, Mistakes: 5},
		'c': {Presses: 10, Mistakes: 2},
		' ': {Presses: 10, Mistakes: 9}, // spaces are skipped
		'z': {Presses: 1, Mistakes: 1},  // too few presses
	}

	tests := []struct {
		name string
		n    int
		want []rune
	}{
		{"top three", 3, []rune{'b', 'c', 'a'}},
		{"top two", 2, []rune{'b', 'c'}},
		{"more than available", 10, []rune{'b', 'c', 'a'}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WeakKeys(stats, tt.n)
			if len(got) != len(tt.want) {
				t.Fatalf("got %d keys, want %d", len(got), len(tt.want))
			}
			for i, r := range tt.want {
				if got[i].Key != r {
					t.Errorf("rank %d = %c, want %c", i, got[i].Key, r)
				}
			}
		})
	}
}

func TestKeyStatAdd(t *testing.T) {
	a := KeyStat{Presses: 2, Mistakes: 1, Total: 100 * time.Millisecond, Timed: 1}
	b := KeyStat{Presses: 3, Mistakes: 0, Total: 200 * time.Millisecond, Timed: 2}

	got := a.Add(b)
	want := KeyStat{Presses: 5, Mistakes: 1, Total: 300 * time.Millisecond, Timed: 3}
	if got != want {
		t.Errorf("Add = %+v, want %+v", got, want)
	}
}

func TestWeights(t *testing.T) {
	ms := time.Millisecond
	stats := map[rune]KeyStat{
		'a': {Presses: 100, Mistakes: 10, Total: 9000 * ms, Timed: 90},
		'b': {Presses: 100, Mistakes: 0, Total: 10000 * ms, Timed: 100},
		'd': {Presses: 100, Mistakes: 0, Total: 20000 * ms, Timed: 100},
		'c': {Presses: 4, Mistakes: 4},
		' ': {Presses: 100, Mistakes: 50},
	}

	got := Weights(stats)

	if !near(got['a'], 2.0) {
		t.Errorf("a = %.3f, want 2.0 (10%% errors)", got['a'])
	}
	if !near(got['d'], 1.0) {
		t.Errorf("d = %.3f, want 1.0 (50%% slower than the mean)", got['d'])
	}
	for _, r := range []rune{'b', 'c', ' '} {
		if _, ok := got[r]; ok {
			t.Errorf("%q should have no weight, got %.3f", r, got[r])
		}
	}
}
