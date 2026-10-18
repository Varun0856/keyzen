package ui

func wrapLines(text []rune, width int) [][2]int {
	if width < 1 {
		width = 1
	}

	var lines [][2]int
	start := 0
	lastSpace := -1

	for i, r := range text {
		if r == ' ' {
			lastSpace = i
			continue
		}
		if i-start+1 > width && lastSpace >= start {
			lines = append(lines, [2]int{start, lastSpace + 1})
			start = lastSpace + 1
			lastSpace = -1
		}
	}
	return append(lines, [2]int{start, len(text)})
}
