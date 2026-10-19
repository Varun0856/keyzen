package typing

import "time"

func WPMSeries(events []KeyEvent, bucket time.Duration) []float64 {
	if len(events) == 0 || bucket <= 0 {
		return nil
	}
	last := events[len(events)-1].At
	counts := make([]int, int(last/bucket)+1)
	for _, e := range events {
		if e.Typed != e.Expected {
			continue
		}
		counts[int(e.At/bucket)]++
	}
	out := make([]float64, len(counts))
	for i, c := range counts {
		out[i] = float64(c) / 5 / bucket.Minutes()
	}
	return out
}
