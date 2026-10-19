package ui

import (
	"math"
	"strings"
	"time"
)

func barParts(frac float64, width int) (filled, empty int) {
	if width < 0 {
		width = 0
	}
	filled = int(frac*float64(width) + 0.5)
	if filled < 0 {
		filled = 0
	}
	if filled > width {
		filled = width
	}
	return filled, width - filled
}

func progressBar(frac float64, width int) string {
	f, e := barParts(frac, width)
	return accentStyle.Render(strings.Repeat("━", f)) + dimStyle.Render(strings.Repeat("━", e))
}

var sparkRunes = []rune("▁▂▃▄▅▆▇█")

func sparkline(values []float64) string {
	if len(values) == 0 {
		return ""
	}
	peak := 0.0
	for _, v := range values {
		if v > peak {
			peak = v
		}
	}

	out := make([]rune, len(values))
	for i, v := range values {
		idx := 0
		if peak > 0 {
			idx = int(v/peak*float64(len(sparkRunes)-1) + 0.5)
		}
		if idx > len(sparkRunes)-1 {
			idx = len(sparkRunes) - 1
		}
		out[i] = sparkRunes[idx]
	}
	return string(out)
}

func seriesBucket(elapsed time.Duration) time.Duration {
	secs := int(math.Ceil(elapsed.Seconds() / 30))
	if secs < 1 {
		secs = 1
	}
	return time.Duration(secs) * time.Second
}
