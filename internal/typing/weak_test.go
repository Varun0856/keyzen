package typing

import "testing"

func TestTroublesome(t *testing.T) {
	stats := map[rune]KeyStat{
		'a': {Presses: 10, Mistakes: 0}, // clean, so never listed
		'b': {Presses: 10, Mistakes: 3},
		'c': {Presses: 10, Mistakes: 1},
	}

	got := Troublesome(stats, 5)
	want := []rune{'b', 'c'}

	if len(got) != len(want) {
		t.Fatalf("got %d keys, want %d", len(got), len(want))
	}
	for i, r := range want {
		if got[i].Key != r {
			t.Errorf("rank %d = %c, want %c", i, got[i].Key, r)
		}
	}
}

func TestTroublesomeNoMistakes(t *testing.T) {
	stats := map[rune]KeyStat{'a': {Presses: 10}}
	if got := Troublesome(stats, 3); len(got) != 0 {
		t.Fatalf("expected no keys, got %v", got)
	}
}
