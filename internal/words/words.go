package words

import (
	_ "embed"
	"math/rand"
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
