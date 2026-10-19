package typing

import (
	"sort"
	"time"
)

type KeyStat struct {
	Presses  int           `json:"presses"`
	Mistakes int           `json:"mistakes"`
	Total    time.Duration `json:"total"`
	Timed    int           `json:"timed"`
}

func (k KeyStat) Add(o KeyStat) KeyStat {
	return KeyStat{
		Presses:  k.Presses + o.Presses,
		Mistakes: k.Mistakes + o.Mistakes,
		Total:    k.Total + o.Total,
		Timed:    k.Timed + o.Timed,
	}
}

func (k KeyStat) ErrorRate() float64 {
	if k.Presses == 0 {
		return 0
	}
	return float64(k.Mistakes) / float64(k.Presses)
}

func (k KeyStat) AvgInterval() time.Duration {
	if k.Timed == 0 {
		return 0
	}
	return k.Total / time.Duration(k.Timed)
}

func AnalyzeKeys(events []KeyEvent) map[rune]KeyStat {
	stats := map[rune]KeyStat{}
	for i, e := range events {
		st := stats[e.Expected]
		st.Presses++
		if e.Typed != e.Expected {
			st.Mistakes++
		} else if i > 0 {
			st.Total += e.At - events[i-1].At
			st.Timed++
		}
		stats[e.Expected] = st
	}
	return stats
}

type KeyRank struct {
	Key  rune
	Stat KeyStat
}

func WeakKeys(stats map[rune]KeyStat, n int) []KeyRank {
	var ranks []KeyRank
	for r, st := range stats {
		if r == ' ' || st.Presses < 2 {
			continue
		}
		ranks = append(ranks, KeyRank{Key: r, Stat: st})
	}

	sort.Slice(ranks, func(i, j int) bool {
		a, b := ranks[i].Stat, ranks[j].Stat
		if a.ErrorRate() != b.ErrorRate() {
			return a.ErrorRate() > b.ErrorRate()
		}
		if a.AvgInterval() != b.AvgInterval() {
			return a.AvgInterval() > b.AvgInterval()
		}
		return ranks[i].Key < ranks[j].Key
	})

	if len(ranks) > n {
		ranks = ranks[:n]
	}
	return ranks
}

func Weights(stats map[rune]KeyStat) map[rune]float64 {
	var sum time.Duration
	n := 0
	for _, st := range stats {
		if st.Timed >= 5 {
			sum += st.AvgInterval()
			n++
		}
	}

	var mean time.Duration
	if n > 0 {
		mean = sum / time.Duration(n)
	}

	out := map[rune]float64{}
	for r, st := range stats {
		if r == ' ' || st.Presses < 5 {
			continue
		}
		w := st.ErrorRate() * 20
		if mean > 0 && st.Timed >= 5 {
			if slow := float64(st.AvgInterval())/float64(mean) - 1; slow > 0 {
				w += slow * 2
			}
		}
		if w > 0 {
			out[r] = w
		}
	}
	return out
}
