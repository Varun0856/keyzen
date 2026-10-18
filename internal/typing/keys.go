package typing

import (
	"sort"
	"time"
)

type KeyStat struct {
	Presses  int
	Mistakes int
	Total    time.Duration
	Timed    int
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
