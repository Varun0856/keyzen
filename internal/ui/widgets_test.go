package ui

import (
	"testing"
	"time"
)

func TestBarParts(t *testing.T) {
	tests := []struct {
		name       string
		frac       float64
		width      int
		wantFilled int
		wantEmpty  int
	}{
		{"empty", 0, 10, 0, 10},
		{"half", 0.5, 10, 5, 5},
		{"full", 1, 10, 10, 0},
		{"over 100%", 1.5, 10, 10, 0},
		{"negative", -1, 10, 0, 10},
		{"zero width", 0.5, 0, 0, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, e := barParts(tt.frac, tt.width)
			if f != tt.wantFilled || e != tt.wantEmpty {
				t.Errorf("got (%d, %d), want (%d, %d)", f, e, tt.wantFilled, tt.wantEmpty)
			}
		})
	}
}

func TestSparkline(t *testing.T) {
	tests := []struct {
		name string
		in   []float64
		want string
	}{
		{"nil", nil, ""},
		{"low and high", []float64{0, 10}, "▁█"},
		{"all equal", []float64{10, 10}, "██"},
		{"all zero", []float64{0, 0}, "▁▁"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sparkline(tt.in); got != tt.want {
				t.Errorf("sparkline(%v) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestSeriesBucket(t *testing.T) {
	tests := []struct {
		elapsed time.Duration
		want    time.Duration
	}{
		{15 * time.Second, time.Second},
		{30 * time.Second, time.Second},
		{31 * time.Second, 2 * time.Second},
		{120 * time.Second, 4 * time.Second},
	}
	for _, tt := range tests {
		if got := seriesBucket(tt.elapsed); got != tt.want {
			t.Errorf("seriesBucket(%v) = %v, want %v", tt.elapsed, got, tt.want)
		}
	}
}
