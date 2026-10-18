package typing

import "testing"

func TestSessionState(t *testing.T) {
	tests := []struct {
		name   string
		target string
		input  string
		want   []CharState
	}{
		{"nothing typed", "abc", "", []CharState{Pending, Pending, Pending}},
		{"partial", "abc", "a", []CharState{Correct, Pending, Pending}},
		{"all correct", "abc", "abc", []CharState{Correct, Correct, Correct}},
		{"one wrong", "abc", "axc", []CharState{Correct, Incorrect, Correct}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSession(tt.target)
			for _, r := range tt.input {
				s.Type(r)
			}
			for i, want := range tt.want {
				if got := s.State(i); got != want {
					t.Errorf("State(%d) = %v, want %v", i, got, want)
				}
			}
		})
	}
}

func TestSessionBackSpaceAndOverflow(t *testing.T) {
	s := NewSession("ab")

	s.Backspace()
	s.Type('a')
	s.Type('b')
	s.Type('c')

	if !s.Done() || s.Cursor() != 2 {
		t.Fatalf("expected done with cursor 2, got done=%v cursor=%d", s.Done(), s.Cursor())
	}

	s.Backspace()
	if s.Done() || s.Cursor() != 1 {
		t.Fatalf("after backspace expected cursor 1, got %d", s.Cursor())
	}
}

func TestSessionStartsOnFirstKey(t *testing.T) {
	s := NewSession("ab")
	if s.Started() {
		t.Fatal("should not be started before any input")
	}
	s.Backspace()
	if s.Started() {
		t.Fatal("backspace should not started the clock")
	}

	s.Type('a')
	if !s.Started() {
		t.Fatal("first keystroke should start the clock")
	}
}

func TestMistakesCountEvenAfterBackspace(t *testing.T) {
	s := NewSession("ab")
	s.Type('x') // wrong
	s.Backspace()
	s.Type('a')
	s.Type('b')

	if s.keystrokes != 3 || s.mistakes != 1 {
		t.Fatalf("keystrokes=%d mistakes=%d, want 3 and 1", s.keystrokes, s.mistakes)
	}
	if s.CorrectCount() != 2 {
		t.Fatalf("CorrectCount = %d, want 2", s.CorrectCount())
	}
}

func TestSessionRecordsEvents(t *testing.T) {
	s := NewSession("ab")
	s.Type('x') // wrong
	s.Backspace()
	s.Type('a')
	s.Type('b')

	ev := s.Events()
	if len(ev) != 3 {
		t.Fatalf("got %d events, want 3", len(ev))
	}
	if ev[0].Expected != 'a' || ev[0].Typed != 'x' {
		t.Errorf("event 0 = %+v, want expected a typed x", ev[0])
	}
	if ev[2].Expected != 'b' || ev[2].Typed != 'b' {
		t.Errorf("event 2 = %+v, want expected b typed b", ev[2])
	}
	if ev[2].At < ev[0].At {
		t.Errorf("timestamps went backwards: %v then %v", ev[0].At, ev[2].At)
	}
}
