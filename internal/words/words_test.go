package words

import (
	"strings"
	"testing"
)

func TestRandom(t *testing.T) {
	tests := []struct {
		name string
		n    int
	}{
		{"one word", 1},
		{"ten word", 10},
		{"fifty words", 50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := strings.Fields(Random(tt.n))
			if len(got) != tt.n {
				t.Errorf("RAndom(%d) returned %d words", tt.n, len(got))
			}
		})
	}
}

func TestListLoaded(t *testing.T) {
	if len(list) == 0 {
		t.Fatal("embedded word list is empty")
	}
}

func TestPickWeighted(t *testing.T) {
	list := []string{"aaa", "bbb"}
	weights := map[rune]float64{'a': 10}

	tests := []struct {
		name string
		roll float64
		want string
	}{
		{"low roll", 0.0, "aaa"},
		{"middle roll", 0.5, "aaa"},
		{"top of range", 0.99, "bbb"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pickWeighted(list, 1, weights, func() float64 { return tt.roll })
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestPickWeightedBias(t *testing.T) {
	list := []string{"aaa", "bbb"}
	weights := map[rune]float64{'a': 10}

	k := 0
	pick := func() float64 {
		v := float64(k) / 1000
		k++
		return v
	}

	aaa := 0
	for _, w := range strings.Fields(pickWeighted(list, 1000, weights, pick)) {
		if w == "aaa" {
			aaa++
		}
	}
	if aaa < 900 {
		t.Errorf("aaa picked %d/1000 times, expected well over 900", aaa)
	}
}
