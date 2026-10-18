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
