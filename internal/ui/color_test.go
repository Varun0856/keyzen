package ui

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestLerpColor(t *testing.T) {
	tests := []struct {
		name string
		a, b lipgloss.Color
		t    float64
		want lipgloss.Color
	}{
		{"start", "#000000", "#ffffff", 0, "#000000"},
		{"end", "#000000", "#ffffff", 1, "#ffffff"},
		{"middle", "#000000", "#ffffff", 0.5, "#808080"},
		{"t above 1 is clamped", "#000000", "#ffffff", 2, "#ffffff"},
		{"t below 0 is clamped", "#000000", "#ffffff", -1, "#000000"},
		{"invalid input returns a", "red", "#ffffff", 0.5, "red"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := lerpColor(tt.a, tt.b, tt.t); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}
