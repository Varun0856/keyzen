package ui

import (
	"sort"
	"time"

	"github.com/charmbracelet/lipgloss"
)

type level int

const (
	levelNone level = iota
	levelGood
	levelOK
	levelBad
	levelWorst
)

func levelForError(rate float64, presses int) level {
	switch {
	case presses < 5:
		return levelNone
	case rate < 0.02:
		return levelGood
	case rate < 0.05:
		return levelOK
	case rate < 0.10:
		return levelBad
	}
	return levelWorst
}

func levelForSpeed(avg, med time.Duration, timed int) level {
	if timed < 5 || med == 0 {
		return levelNone
	}
	ratio := float64(avg) / float64(med)
	switch {
	case ratio <= 0.9:
		return levelGood
	case ratio <= 1.2:
		return levelOK
	case ratio <= 1.6:
		return levelBad
	}
	return levelWorst
}

func median(ds []time.Duration) time.Duration {
	if len(ds) == 0 {
		return 0
	}
	sorted := make([]time.Duration, len(ds))
	copy(sorted, ds)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })

	mid := len(sorted) / 2
	if len(sorted)%2 == 1 {
		return sorted[mid]
	}
	return (sorted[mid-1] + sorted[mid]) / 2
}

func levelStyle(l level) lipgloss.Style {
	base := lipgloss.NewStyle()
	switch l {
	case levelGood:
		return base.Background(lipgloss.Color("2")).Foreground(lipgloss.Color("0"))
	case levelOK:
		return base.Background(lipgloss.Color("3")).Foreground(lipgloss.Color("0"))
	case levelBad:
		return base.Background(lipgloss.Color("208")).Foreground(lipgloss.Color("0"))
	case levelWorst:
		return base.Background(lipgloss.Color("1")).Foreground(lipgloss.Color("15"))
	}
	return base.Background(lipgloss.Color("236")).Foreground(lipgloss.Color("245"))
}
