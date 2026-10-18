package main

import (
	"reflect"
	"testing"
)

func TestWrapLines(t *testing.T) {
	tests := []struct {
		name  string
		text  string
		width int
		want  [][2]int
	}{
		{"fits on one line", "aaa bbb ccc", 100, [][2]int{{0, 11}}},
		{"wraps at a space", "aaa bbb ccc", 7, [][2]int{{0, 8}, {8, 11}}},
		{"long word stays whole", "abcdefghij kl", 5, [][2]int{{0, 11}, {11, 13}}},
		{"empty text", "", 10, [][2]int{{0, 0}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := wrapLines([]rune(tt.text), tt.width)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("wrapLines(%q, %d) = %v, want %v", tt.text, tt.width, got, tt.want)
			}
		})
	}
}
