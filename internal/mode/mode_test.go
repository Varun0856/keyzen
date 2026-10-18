package mode

import "testing"

func TestKeyAndLength(t *testing.T) {
	tests := []struct {
		name    string
		mode    Mode
		wantKey string
		wantLen int
	}{
		{"timed 30", timed(30), "30s", 90},
		{"timed 15", timed(15), "15s", 45},
		{"25 words", wordCount(25), "25w", 25},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.mode.Key(); got != tt.wantKey {
				t.Errorf("Key() = %q, want %q", got, tt.wantKey)
			}
			if got := tt.mode.TextLength(); got != tt.wantLen {
				t.Errorf("TextLength() = %d, want %d", got, tt.wantLen)
			}
		})
	}
}

func TestAllKeysUnique(t *testing.T) {
	seen := map[string]bool{}
	for _, m := range All() {
		if seen[m.Key()] {
			t.Fatalf("duplicate key %q", m.Key())
		}
		seen[m.Key()] = true
	}
}
