package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/Varun0856/keyzen/internal/storage"
	"github.com/Varun0856/keyzen/internal/typing"
	tea "github.com/charmbracelet/bubbletea"
)

var keyboardRows = []string{"qwertyuiop", "asdfghjkl", "zxcvbnm"}

type statsScreen struct {
	stats         map[rune]typing.KeyStat
	median        time.Duration
	total         int
	bySpeed       bool
	width, height int
}

func newStats(keys *storage.KeyStore, w, h int) statsScreen {
	stats := map[rune]typing.KeyStat{}
	if keys != nil {
		if t, err := keys.Load(); err == nil {
			stats = t.Runes()
		}
	}

	total := 0
	var avgs []time.Duration
	for _, st := range stats {
		total += st.Presses
		if st.Timed >= 5 {
			avgs = append(avgs, st.AvgInterval())
		}
	}
	return statsScreen{stats: stats, median: median(avgs), total: total, width: w, height: h}
}

func (m statsScreen) Init() tea.Cmd { return nil }

func (m statsScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			m.bySpeed = !m.bySpeed
		case "esc", "m":
			return m, send(showMenuMsg{})
		case "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m statsScreen) levelFor(r rune) level {
	st, ok := m.stats[r]
	if !ok {
		return levelNone
	}
	if m.bySpeed {
		return levelForSpeed(st.AvgInterval(), m.median, st.Timed)
	}
	return levelForError(st.ErrorRate(), st.Presses)
}

func (m statsScreen) View() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("stats") + "\n")

	if m.total == 0 {
		b.WriteString(dimStyle.Render("no data yet, finish a test first") + "\n\n")
		b.WriteString(dimStyle.Render("esc menu • q quit") + "\n")
		return place(m.width, m.height, 44, b.String())
	}

	metric := "accuracy"
	if m.bySpeed {
		metric = "speed"
	}
	b.WriteString(dimStyle.Render(fmt.Sprintf("%d keystrokes recorded", m.total)) + "\n\n")
	b.WriteString(groupStyle.Render("by "+metric) + "\n\n")

	for i, row := range keyboardRows {
		b.WriteString(strings.Repeat(" ", i*2))
		for _, r := range row {
			b.WriteString(levelStyle(m.levelFor(r)).Render(" "+string(r)+" ") + " ")
		}
		b.WriteString("\n")
	}

	b.WriteString("\n" +
		levelStyle(levelGood).Render(" good ") + " " +
		levelStyle(levelOK).Render(" ok ") + " " +
		levelStyle(levelBad).Render(" weak ") + " " +
		levelStyle(levelWorst).Render(" worst ") + " " +
		levelStyle(levelNone).Render(" no data "+"\n"))

	if weak := typing.WeakKeys(m.stats, 5); len(weak) > 0 {
		b.WriteString("\n" + groupStyle.Render("weakest keys") + "\n")
		for _, k := range weak {
			fmt.Fprintf(&b, "%s %3.0f%% acc %4dms\n",
				bigStyle.Render(string(k.Key)),
				(1-k.Stat.ErrorRate())*100,
				k.Stat.AvgInterval().Milliseconds())
		}
	}
	b.WriteString("\n" + dimStyle.Render("tab accuracy/speed • esc menu • q quit") + "\n")
	return place(m.width, m.height, 44, b.String())
}
