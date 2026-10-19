package words

import (
	_ "embed"
	"math/rand"
	"sort"
	"strings"
)

//go:embed words.txt
var raw string
var list = strings.Fields(raw)

func Random(n int) string {
	picked := make([]string, n)
	for i := range picked {
		picked[i] = list[rand.Intn(len(list))]
	}
	return strings.Join(picked, " ")
}

func pickWeighted(list []string, n int, weights map[rune]float64, pick func() float64) string {
	cum := make([]float64, len(list))
	total := 0.0
	for i, w := range list {
		score := 1.0
		for _, r := range w {
			score += weights[r]
		}
		total += score
		cum[i] = total
	}
	picked := make([]string, n)
	for i := range picked {
		x := pick() * total
		idx := sort.Search(len(cum), func(j int) bool { return cum[j] > x })
		if idx >= len(list) {
			idx = len(list) - 1
		}
		picked[i] = list[idx]
	}
	return strings.Join(picked, " ")
}

func Weighted(n int, weights map[rune]float64, pick func() float64) string {
	return pickWeighted(list, n, weights, pick)
}
