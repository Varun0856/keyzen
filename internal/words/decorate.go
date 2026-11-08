package words

import (
	"strconv"
	"strings"
	"unicode"
)

type Options struct {
	Punctuation bool
	Numbers     bool
}

func (o Options) Suffix() string {
	s := ""
	if o.Punctuation {
		s += "+p"
	}
	if o.Numbers {
		s += "+n"
	}
	return s
}

func (o Options) Describe() string {
	var parts []string
	if o.Punctuation {
		parts = append(parts, "punctuation")
	}
	if o.Numbers {
		parts = append(parts, "numbers")
	}
	return strings.Join(parts, ", ")
}

func Decorate(text string, opts Options, rnd func() float64) string {
	ws := strings.Fields(text)
	if len(ws) == 0 || (!opts.Punctuation && !opts.Numbers) {
		return text
	}

	if opts.Numbers {
		for i := range ws {
			if rnd() < 0.12 {
				ws[i] = number(rnd)
			}
		}
	}

	if opts.Punctuation {
		newSentence := true
		sinceEnd := 0 // words in the current sentence so far
		for i, w := range ws {
			if newSentence {
				w = capitalize(w)
				newSentence = false
			}
			sinceEnd++

			switch {
			case i == len(ws)-1, sinceEnd >= 14, sinceEnd >= 5 && rnd() < 0.25:
				w += endMark(rnd)
				newSentence = true
				sinceEnd = 0
			case sinceEnd >= 2 && rnd() < 0.12:
				w += ","
			}
			ws[i] = w
		}
	}
	return strings.Join(ws, " ")
}

func digit(rnd func() float64, lo int) int {
	d := lo + int(rnd()*float64(10-lo))
	if d > 9 {
		d = 9
	}
	return d
}

func number(rnd func() float64) string {
	n := 1 + int(rnd()*4)
	if n > 4 {
		n = 4
	}
	var b strings.Builder
	for i := 0; i < n; i++ {
		lo := 0
		if i == 0 {
			lo = 1
		}
		b.WriteString(strconv.Itoa(digit(rnd, lo)))
	}
	return b.String()
}

func endMark(rnd func() float64) string {
	r := rnd()
	switch {
	case r < 0.85:
		return "."
	case r < 0.93:
		return "?"
	}
	return "!"
}

func capitalize(w string) string {
	r := []rune(w)
	if len(r) == 0 {
		return w
	}
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}
