package typing

import (
	"math"
	"testing"
	"time"
)

func near(a, b float64) bool { return math.Abs(a-b) < 0.01 }

func TestCompute(t *testing.T) {
	tests := []struct {
		name       string
		correct    int
		keystrokes int
		mistakes   int
		elapsed    time.Duration
		wantWPM    float64
		wantAcc    float64
	}{
		{"50 chars in a minute", 50, 50, 0, time.Minute, 10, 100},
		{"100 chars in 30s", 100, 100, 0, 30 * time.Second, 40, 100},
		{"with mistakes", 90, 100, 10, time.Minute, 18, 90},
		{"no time elapsed", 10, 10, 0, 0, 0, 100},
		{"no keystrokes", 0, 0, 0, time.Minute, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Compute(tt.correct, tt.keystrokes, tt.mistakes, tt.elapsed)
			if !near(got.WPM, tt.wantWPM) {
				t.Errorf("WPM = %.2f, want %.2f", got.WPM, tt.wantWPM)
			}
			if !near(got.Accuracy, tt.wantAcc) {
				t.Errorf("Accuracy = %.2f, want %.2f", got.Accuracy, tt.wantAcc)
			}
		})
	}
}

func TestMistakesCountEventAfterBackspace(t *testing.T) {
	s := NewSession("ab")
	s.Type('x')
	s.Backspace()
	s.Type('a')
	s.Type('b')

	if s.keystrokes != 3 || s.mistakes != 1 {
		t.Fatalf("keystrokes=%d mistakes=%d, want 3 and 1", s.keystrokes, s.mistakes)
	}

	if s.CorrectCount() != 2 {
		t.Fatalf("CorrectCount = %d, want 2", s.CorrectCount())
	}
}
