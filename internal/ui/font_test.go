package ui

import (
	"testing"
	"unicode/utf8"
)

func TestGlyphsAreRectangular(t *testing.T) {
	for r, g := range glyphs {
		w := utf8.RuneCountInString(g[0])
		for i, row := range g {
			if got := utf8.RuneCountInString(row); got != w {
				t.Errorf("glyph %q: row %d is %d wide, row 0 is %d", r, i, got, w)
			}
		}
	}
}

func TestGlyphsExist(t *testing.T) {
	for _, r := range "keyzen0123456789 " {
		if _, ok := glyphs[r]; !ok {
			t.Errorf("missing glyph for %q", r)
		}
	}
}

func TestBigWidth(t *testing.T) {
	tests := []struct {
		name  string
		text  string
		scale int
		want  int
	}{
		{"one letter", "k", 1, 5},
		{"two letters plus gap", "ke", 1, 11},
		{"scaled", "ke", 2, 22},
		{"unknown runes are skipped", "k?", 1, 5},
		{"zero scale acts as 1", "k", 0, 5},
		{"empty", "", 1, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := bigWidth(tt.text, tt.scale); got != tt.want {
				t.Errorf("bigWidth(%q, %d) = %d, want %d", tt.text, tt.scale, got, tt.want)
			}
		})
	}
}

func TestBigRowsHeight(t *testing.T) {
	if got := len(bigRows("keyzen", 2)); got != glyphRows {
		t.Fatalf("got %d rows, want %d", got, glyphRows)
	}
}
