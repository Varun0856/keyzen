package ui

import (
	"testing"
	"time"
)

func TestLevelForError(t *testing.T) {
	tests := []struct {
		name    string
		rate    float64
		presses int
		want    level
	}{
		{"too few presses", 0.5, 4, levelNone},
		{"clean", 0.0, 50, levelGood},
		{"just under 2%", 0.019, 50, levelGood},
		{"2%", 0.02, 50, levelOK},
		{"5%", 0.05, 50, levelBad},
		{"10%", 0.10, 50, levelWorst},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := levelForError(tt.rate, tt.presses); got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLevelForSpeed(t *testing.T) {
	ms := time.Millisecond
	tests := []struct {
		name  string
		avg   time.Duration
		med   time.Duration
		timed int
		want  level
	}{
		{"too few samples", 100 * ms, 200 * ms, 4, levelNone},
		{"no median yet", 100 * ms, 0, 50, levelNone},
		{"much faster than median", 100 * ms, 200 * ms, 50, levelGood},
		{"about median", 200 * ms, 200 * ms, 50, levelOK},
		{"1.5x median", 300 * ms, 200 * ms, 50, levelBad},
		{"2x median", 400 * ms, 200 * ms, 50, levelWorst},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := levelForSpeed(tt.avg, tt.med, tt.timed); got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMedian(t *testing.T) {
	ms := time.Millisecond
	tests := []struct {
		name string
		in   []time.Duration
		want time.Duration
	}{
		{"empty", nil, 0},
		{"one", []time.Duration{5 * ms}, 5 * ms},
		{"odd", []time.Duration{30 * ms, 10 * ms, 20 * ms}, 20 * ms},
		{"even", []time.Duration{10 * ms, 40 * ms, 20 * ms, 30 * ms}, 25 * ms},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := median(tt.in); got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
