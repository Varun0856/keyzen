package typing

import "time"

type CharState int

const (
	Pending CharState = iota
	Correct
	Incorrect
)

type KeyEvent struct {
	Expected rune
	Typed    rune
	At       time.Duration
}

type Session struct {
	target     []rune
	typed      []rune
	startedAt  time.Time
	endedAt    time.Time
	keystrokes int
	mistakes   int
	events     []KeyEvent
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
	expected := s.target[len(s.typed)]
	s.keystrokes++
	if r != expected {
		s.mistakes++
	}
	s.events = append(s.events, KeyEvent{
		Expected: expected,
		Typed:    r,
		At:       time.Since(s.startedAt),
	})
	s.typed = append(s.typed, r)

	if s.Done() && s.endedAt.IsZero() {
		s.endedAt = time.Now()
	}
}

func (s *Session) Events() []KeyEvent { return s.events }

func (s *Session) Started() bool { return !s.startedAt.IsZero() }

func (s *Session) Backspace() {
	if len(s.typed) > 0 {
		s.typed = s.typed[:len(s.typed)-1]
	}
}

func (s *Session) Finish() {
	if s.Started() && s.endedAt.IsZero() {
		s.endedAt = time.Now()
	}
}

func (s *Session) Done() bool { return len(s.typed) >= len(s.target) }

func (s *Session) Elapsed() time.Duration {
	if !s.Started() {
		return 0
	}
	if s.endedAt.IsZero() {
		return time.Since(s.startedAt)
	}
	return s.endedAt.Sub(s.startedAt)
}

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

func (s *Session) CorrectCount() int {
	n := 0
	for i := range s.typed {
		if s.State(i) == Correct {
			n++
		}
	}
	return n
}

func (s *Session) Stats() Stats {
	return Compute(s.CorrectCount(), s.keystrokes, s.mistakes, s.Elapsed())
}
