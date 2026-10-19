package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func lerpColor(a, b lipgloss.Color, t float64) lipgloss.Color {
	var ar, ag, ab, br, bg, bb int
	if _, err := fmt.Sscanf(string(a), "#%02x%02x%02x", &ar, &ag, &ab); err != nil {
		return a
	}
	if _, err := fmt.Sscanf(string(b), "#%02x%02x%02x", &br, &bg, &bb); err != nil {
		return a
	}
	if t < 0 {
		t = 0
	}
	if t > 1 {
		t = 1
	}

	mix := func(x, y int) int {
		return int(float64(x) + (float64(y)-float64(x))*t + 0.5)
	}
	return lipgloss.Color(fmt.Sprintf("#%02x%02x%02x", mix(ar, br), mix(ag, bg), mix(ab, bb)))
}
