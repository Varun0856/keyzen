package ui

import (
	"testing"

	"github.com/Varun0856/keyzen/internal/mode"
)

func TestGroupModes(t *testing.T) {
	got := groupModes(mode.All())

	wantNames := []string{"Timed", "Word count", "Practice"}
	if len(got) != len(wantNames) {
		t.Fatalf("got %d groups, want %d", len(got), len(wantNames))
	}
	for i, n := range wantNames {
		if got[i].name != n {
			t.Errorf("group %d = %q, want %q", i, got[i].name, n)
		}
	}
	if len(got[0].modes) != 4 || len(got[2].modes) != 1 {
		t.Errorf("unexpected group sizes: %d, %d", len(got[0].modes), len(got[2].modes))
	}
}

func TestGroupModesNonConsecutive(t *testing.T) {
	in := []mode.Mode{{Group: "A"}, {Group: "B"}, {Group: "A"}}
	got := groupModes(in)

	if len(got) != 2 || len(got[0].modes) != 2 || len(got[1].modes) != 1 {
		t.Fatalf("groupModes merged wrongly: %+v", got)
	}
}
