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
