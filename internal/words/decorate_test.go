package words

import (
	"math/rand"
	"strings"
	"testing"
	"unicode"
)

func constRnd(v float64) func() float64 { return func() float64 { return v } }

func TestDecorateNoOptionsChangesNothing(t *testing.T) {
	in := "the quick brown fox"
	if got := Decorate(in, Options{}, constRnd(0.5)); got != in {
		t.Errorf("got %q, want unchanged %q", got, in)
	}
}

func TestDecoratePunctuation(t *testing.T) {
	tests := []struct {
		name string
		in   string
		rnd  float64
		want string
	}{
		{"capital and full stop", "the quick brown fox", 0.5, "The quick brown fox."},
		{"single word", "go", 0.5, "Go."},
		{"commas and an early sentence end", "a b c d e f", 0.1, "A b, c, d, e. F."},
		{"question mark", "is it ok", 0.9, "Is it ok?"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Decorate(tt.in, Options{Punctuation: true}, constRnd(tt.rnd))
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDecorateLongSentenceIsSplit(t *testing.T) {
	in := strings.TrimSpace(strings.Repeat("a ", 15))
	got := Decorate(in, Options{Punctuation: true}, constRnd(0.5))

	if strings.Count(got, ".") != 2 || !strings.HasPrefix(got, "A a") || !strings.HasSuffix(got, "A.") {
		t.Errorf("a 15 word run should become two sentences, got %q", got)
	}
}

func TestDecorateNumbers(t *testing.T) {
	tests := []struct {
		name string
		rnd  float64
		want string
	}{
		{"low roll replaces words", 0.05, "1 1"},
		{"high roll leaves words alone", 0.5, "the quick"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Decorate("the quick", Options{Numbers: true}, constRnd(tt.rnd))
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDecorateKeepsWordCount(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	in := strings.Repeat("word ", 200)
	got := Decorate(in, Options{Punctuation: true, Numbers: true}, rng.Float64)

	fields := strings.Fields(got)
	if len(fields) != 200 {
		t.Fatalf("got %d words, want 200", len(fields))
	}
	if last := fields[len(fields)-1]; !strings.ContainsAny(last[len(last)-1:], ".?!") {
		t.Errorf("text should end with a sentence mark, got %q", last)
	}
	for _, w := range fields {
		core := strings.TrimRight(w, ".,?!")
		allDigits := core != ""
		for _, r := range core {
			if !unicode.IsDigit(r) {
				allDigits = false
			}
		}
		if allDigits && len(core) > 1 && core[0] == '0' {
			t.Errorf("number %q has a leading zero", w)
		}
	}
}

func TestOptionsSuffixAndDescribe(t *testing.T) {
	tests := []struct {
		name         string
		opts         Options
		suffix, desc string
	}{
		{"none", Options{}, "", ""},
		{"punctuation", Options{Punctuation: true}, "+p", "punctuation"},
		{"numbers", Options{Numbers: true}, "+n", "numbers"},
		{"both", Options{Punctuation: true, Numbers: true}, "+p+n", "punctuation, numbers"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.opts.Suffix(); got != tt.suffix {
				t.Errorf("Suffix() = %q, want %q", got, tt.suffix)
			}
			if got := tt.opts.Describe(); got != tt.desc {
				t.Errorf("Describe() = %q, want %q", got, tt.desc)
			}
		})
	}
}
