package ui

import (
	"strings"
	"unicode/utf8"

	"github.com/charmbracelet/lipgloss"
)

const glyphRows = 5

var glyphs = map[rune][glyphRows]string{
	'k': {"█   █", "█  █ ", "███  ", "█  █ ", "█   █"},
	'e': {"█████", "█    ", "████ ", "█    ", "█████"},
	'y': {"█   █", " █ █ ", "  █  ", "  █  ", "  █  "},
	'z': {"█████", "   █ ", "  █  ", " █   ", "█████"},
	'n': {"█   █", "██  █", "█ █ █", "█  ██", "█   █"},
	' ': {"   ", "   ", "   ", "   ", "   "},

	'0': {"█████", "█   █", "█   █", "█   █", "█████"},
	'1': {" ██  ", "  █  ", "  █  ", "  █  ", " ███ "},
	'2': {"█████", "    █", "█████", "█    ", "█████"},
	'3': {"█████", "    █", " ████", "    █", "█████"},
	'4': {"█   █", "█   █", "█████", "    █", "    █"},
	'5': {"█████", "█    ", "█████", "    █", "█████"},
	'6': {"█████", "█    ", "█████", "█   █", "█████"},
	'7': {"█████", "    █", "   █ ", "  █  ", "  █  "},
	'8': {"█████", "█   █", "█████", "█   █", "█████"},
	'9': {"█████", "█   █", "█████", "    █", "█████"},
}

func stretch(s string, n int) string {
	var b strings.Builder
	for _, r := range s {
		b.WriteString(strings.Repeat(string(r), n))
	}
	return b.String()
}

func bigRows(s string, scale int) []string {
	if scale < 1 {
		scale = 1
	}
	rows := make([]string, glyphRows)
	first := true
	for _, r := range s {
		g, ok := glyphs[r]
		if !ok {
			continue
		}
		for i := range rows {
			if !first {
				rows[i] += strings.Repeat(" ", scale)
			}
			rows[i] += stretch(g[i], scale)
		}
		first = false
	}
	return rows
}

func bigWidth(s string, scale int) int {
	return utf8.RuneCountInString(bigRows(s, scale)[0])
}

func renderBig(s string, scale int, from, to lipgloss.Color) string {
	rows := bigRows(s, scale)
	out := make([]string, len(rows))
	for i, r := range rows {
		t := 0.0
		if len(rows) > 1 {
			t = float64(i) / float64(len(rows)-1)
		}
		out[i] = lipgloss.NewStyle().Foreground(lerpColor(from, to, t)).Render(r)
	}
	return strings.Join(out, "\n")
}
