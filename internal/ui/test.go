package ui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Varun0856/keyzen/internal/mode"
	"github.com/Varun0856/keyzen/internal/storage"
	"github.com/Varun0856/keyzen/internal/typing"
	"github.com/Varun0856/keyzen/internal/words"
)

const (
	maxTextWidth = 90
	visibleLines = 3
)

type tickMsg struct{ gen int }

func tick(gen int) tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		return tickMsg{gen: gen}
	})
}

type testScreen struct {
	md            mode.Mode
	store         *storage.Store
	keys          *storage.KeyStore
	opts          words.Options
	session       *typing.Session
	remaining     int
	finished      bool
	gen           int
	width, height int
}

func newTest(md mode.Mode, store *storage.Store, keys *storage.KeyStore, opts words.Options, gen, w, h int) testScreen {
	return testScreen{
		md:        md,
		store:     store,
		keys:      keys,
		opts:      opts,
		session:   typing.NewSession(buildText(md, keys, opts)),
		remaining: md.Seconds,
		gen:       gen,
		width:     w,
		height:    h,
	}
}

func (m testScreen) Init() tea.Cmd { return nil }

func (m *testScreen) finish() tea.Cmd {
	m.finished = true
	st := m.session.Stats()
	events := m.session.Events()
	perKey := typing.AnalyzeKeys(events)

	res := result{
		md:     m.md,
		stats:  st,
		weak:   typing.Troublesome(perKey, 3),
		series: typing.WPMSeries(events, seriesBucket(st.Elapsed)),
	}

	if m.store != nil && st.WPM > 0 {
		key := m.md.Key() + m.opts.Suffix()
		best, isNew, err := m.store.Record(key, st.WPM, st.Accuracy)
		if err == nil {
			res.best, res.hasBest, res.newBest = best, true, isNew
		}
	}
	if m.keys != nil && len(events) > 0 {
		_ = m.keys.Record(perKey)
	}
	return send(showResultsMsg{res: res})
}

func (m testScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return m, send(showMenuMsg{})
		case tea.KeyTab:
			return m, send(startTestMsg{md: m.md})
		}
		if m.finished {
			return m, nil
		}

		wasStarted := m.session.Started()

		switch msg.Type {
		case tea.KeyBackspace:
			m.session.Backspace()
		case tea.KeySpace:
			m.session.Type(' ')
		case tea.KeyRunes:
			for _, r := range msg.Runes {
				m.session.Type(r)
			}
		}

		if m.session.Done() {
			cmd := m.finish()
			return m, cmd
		}
		if !wasStarted && m.session.Started() && m.md.Kind == mode.Timed {
			return m, tick(m.gen)
		}

	case tickMsg:
		if msg.gen != m.gen || m.finished {
			return m, nil
		}
		m.remaining--
		if m.remaining <= 0 {
			m.session.Finish()
			cmd := m.finish()
			return m, cmd
		}
		return m, tick(m.gen)
	}
	return m, nil
}

func (m testScreen) wordsDone() int {
	n := 0
	for i := 0; i < m.session.Cursor(); i++ {
		if m.session.Char(i) == ' ' {
			n++
		}
	}
	return n
}

func (m testScreen) textWidth() int {
	return clamp(m.width-24, 20, maxTextWidth)
}

func (m testScreen) renderChar(i int) string {
	s := string(m.session.Char(i))
	switch m.session.State(i) {
	case typing.Correct:
		return correctStyle.Render(s)
	case typing.Incorrect:
		return wrongStyle.Render(s)
	}
	if i == m.session.Cursor() {
		return cursorStyle.Render(s)
	}
	return pendingStyle.Render(s)
}

func (m testScreen) textView() string {
	runes := make([]rune, m.session.Len())
	for i := range runes {
		runes[i] = m.session.Char(i)
	}
	lines := wrapLines(runes, m.textWidth())

	cur := len(lines) - 1
	for idx, ln := range lines {
		if m.session.Cursor() >= ln[0] && m.session.Cursor() < ln[1] {
			cur = idx
			break
		}
	}

	first := cur - 1
	if first < 0 {
		first = 0
	}
	last := first + visibleLines
	if last > len(lines) {
		last = len(lines)
	}

	rendered := make([]string, 0, visibleLines)
	for _, ln := range lines[first:last] {
		var b strings.Builder
		for i := ln[0]; i < ln[1]; i++ {
			b.WriteString(m.renderChar(i))
		}
		rendered = append(rendered, b.String())
	}
	for len(rendered) < visibleLines {
		rendered = append(rendered, "")
	}
	return strings.Join(rendered, "\n\n")
}

func (m testScreen) View() string {
	cw := m.textWidth() + 2
	total := cw + panelExtraBig
	started := m.session.Started()

	progress := fmt.Sprint(m.remaining)
	frac := 0.0
	maxNumber := m.md.Seconds
	if m.md.Kind == mode.Timed {
		frac = float64(m.md.Seconds-m.remaining) / float64(m.md.Seconds)
	} else {
		progress = fmt.Sprint(m.wordsDone())
		frac = float64(m.wordsDone()) / float64(m.md.Words)
		maxNumber = m.md.Words
	}

	var stats []string
	if started {
		st := m.session.Stats()
		stats = append(stats,
			bigStyle.Render(fmt.Sprintf("%.0f", st.WPM))+dimStyle.Render(" wpm"),
			bigStyle.Render(fmt.Sprintf("%.0f%%", st.Accuracy))+dimStyle.Render(" accuracy"),
		)
	} else {
		stats = append(stats, dimStyle.Render("start typing"), dimStyle.Render("to begin"))
	}
	if m.md.Kind != mode.Timed {
		stats = append(stats, dimStyle.Render(fmt.Sprintf("of %d words", m.md.Words)))
	}
	statsBlock := strings.Join(stats, "\n")

	var top string
	if m.height >= 28 {
		from, to := currentTheme.Accent, currentTheme.Accent2
		if m.md.Kind == mode.Timed && started && m.remaining <= 5 {
			from, to = currentTheme.Bad, currentTheme.Warn
		}
		heroWidth := bigWidth(strings.Repeat("0", len(fmt.Sprint(maxNumber))), 1)
		hero := lipgloss.NewStyle().Width(heroWidth).Render(renderBig(progress, 1, from, to))
		top = lipgloss.JoinHorizontal(lipgloss.Center, hero, "      ", statsBlock)
	} else {
		top = timerStyle.Render(progress) + "    " + strings.ReplaceAll(statsBlock, "\n", "    ")
	}

	header := accentStyle.Render("keyzen") + dimStyle.Render("  ·  "+m.md.Title())
	if d := m.opts.Describe(); d != "" {
		header += dimStyle.Render("  ·  " + d)
	}
	bar := progressBar(frac, total)
	text := bigPanel(m.textView(), cw)
	foot := hints("tab", "restart", "esc", "menu", "ctrl+c", "quit")

	content := header + "\n\n" + top + "\n\n" + bar + "\n" + text + "\n\n" + foot
	return place(m.width, m.height, total, content)
}
