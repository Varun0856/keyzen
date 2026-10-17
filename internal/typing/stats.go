package typing

import "time"

type Stats struct {
	WPM      float64
	Accuracy float64
	Correct  int
	Elapsed  time.Duration
}

func Compute(correct, keystrokes, mistakes int, elapsed time.Duration) Stats {
	st := Stats{Correct: correct, Elapsed: elapsed}

	if minutes := elapsed.Minutes(); minutes > 0 {
		st.WPM = (float64(correct) / 5) / minutes
	}

	if keystrokes > 0 {
		st.Accuracy = float64(keystrokes-mistakes) / float64(keystrokes) * 100
	}
	return st
}
