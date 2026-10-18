package mode

import (
	"strconv"
)

type Kind int

const (
	Timed Kind = iota
	WordCount
)

type Mode struct {
	Group   string
	Label   string
	Kind    Kind
	Seconds int
	Words   int
}

func timed(s int) Mode {
	return Mode{Group: "Timed", Label: strconv.Itoa(s) + "s", Kind: Timed, Seconds: s}
}

func wordCount(n int) Mode {
	return Mode{Group: "Word count", Label: strconv.Itoa(n) + " words", Kind: WordCount, Words: n}
}

func (m Mode) Key() string {
	if m.Kind == Timed {
		return strconv.Itoa(m.Seconds) + "s"
	}
	return strconv.Itoa(m.Words) + "w"
}

func (m Mode) Title() string { return m.Group + " . " + m.Label }

func (m Mode) TextLength() int {
	if m.Kind == Timed {
		return m.Seconds * 3
	}
	return m.Words
}

func All() []Mode {
	return []Mode{
		timed(15), timed(30), timed(60), timed(120),
		wordCount(10), wordCount(25), wordCount(50), wordCount(100),
	}
}
