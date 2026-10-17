package typing

import "time"

type CharState int

const (
	Pending CharState = iota
	Correct
	Incorrect
)

type Session struct {
	target    []rune
	typed     []rune
	startedAt time.Time
}

func NewSession(target string) *Session {
	return &Session{target: []rune(target)}
}

func (s *Session) Type(r rune) {
	if s.Done() {
		return
	}
	if s.startedAt.IsZero() {
		s.startedAt = time.Now()
	}
	s.typed = append(s.typed, r)
}

func (s *Session) Started() bool { return !s.startedAt.IsZero() }

func (s *Session) Backspace() {
	if len(s.typed) > 0 {
		s.typed = s.typed[:len(s.typed)-1]
	}
}

func (s *Session) Done() bool { return len(s.typed) >= len(s.target) }

func (s *Session) Len() int { return len(s.target) }

func (s *Session) Cursor() int { return len(s.typed) }

func (s *Session) Char(i int) rune { return s.target[i] }

func (s *Session) State(i int) CharState {
	if i >= len(s.typed) {
		return Pending
	}
	if s.typed[i] == s.target[i] {
		return Correct
	}
	return Incorrect
}
