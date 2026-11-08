package typing

func Troublesome(stats map[rune]KeyStat, n int) []KeyRank {
	withMistakes := make(map[rune]KeyStat, len(stats))
	for r, st := range stats {
		if st.Mistakes > 0 {
			withMistakes[r] = st
		}
	}
	return WeakKeys(withMistakes, n)
}
