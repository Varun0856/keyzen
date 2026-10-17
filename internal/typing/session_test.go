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
