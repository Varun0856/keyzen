package typing

import (
	"testing"
	"time"
)

func TestWPMSeries(t *testing.T) {
	ms := time.Millisecond
	events := []KeyEvent{
		{'a', 'a', 0},
		{'b', 'b', 200 * ms},
		{'c', 'x', 400 * ms}, // wrong, not counted
		{'d', 'd', 900 * ms},
		{'e', 'e', 1100 * ms},
		{'f', 'f', 1500 * ms},
	}

	got := WPMSeries(events, time.Second)
	want := []float64{36, 24} // 3 correct in second 1, 2 in second 2

	if len(got) != len(want) {
		t.Fatalf("got %d buckets, want %d", len(got), len(want))
	}
	for i := range want {
		if !near(got[i], want[i]) {
			t.Errorf("bucket %d = %.2f, want %.2f", i, got[i], want[i])
		}
	}
}

func TestWPMSeriesEmpty(t *testing.T) {
	if got := WPMSeries(nil, time.Second); got != nil {
		t.Fatalf("expected nil for no events, got %v", got)
	}
}
