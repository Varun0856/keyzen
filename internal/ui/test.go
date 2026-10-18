package ui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/Varun0856/keyzen/internal/mode"
	"github.com/Varun0856/keyzen/internal/storage"
	"github.com/Varun0856/keyzen/internal/typing"
	"github.com/Varun0856/keyzen/internal/words"
)

const (
	maxTextWidth = 60
	visibleLines = 3
)

// tickMsg carries the generation of the test it belongs to, so a tick
// left over from an earlier test is recognised and ignored.
type tickMsg struct{ gen int }

func tick(gen int) tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		return tickMsg{gen: gen}
	})
}

type testScreen struct {
	md            mode.Mode
	store         *storage.Store
	session       *typing.Session
	remaining     int
	finished      bool
	gen           int
	width, height int
}

func newTest(md mode.Mode, store *storage.Store, gen, w, h int) testScreen {
	return testScreen{
		md:        md,
		store:     store,
		session:   typing.NewSession(words.Random(md.TextLength())),
		remaining: md.Seconds,
		gen:       gen,
		width:     w,
		height:    h,
	}
}

func (m testScreen) Init() tea.Cmd { return nil }

// finish records the score and asks the app to show the results screen.
// Pointer receiver: it changes m.finished.
func (m *testScreen) finish() tea.Cmd {
	m.finished = true
	st := m.session.Stats()
	res := result{md: m.md, stats: st}

	if m.store != nil && st.WPM > 0 {
		best, isNew, err := m.store.Record(m.md.Key(), st.WPM, st.Accuracy)
		if err == nil {
			res.best, res.hasBest, res.newBest = best, true, isNew
		}
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

// wordsDone counts the spaces already passed, i.e. completed words.
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
	w := m.width - 8
	if w > maxTextWidth {
		w = maxTextWidth
	}
	if w < 10 {
		w = 10
	}
	return w
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

// textView draws exactly visibleLines lines, scrolling with the cursor.
// Always the same height, so the screen doesn't jump near the end.
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

	var b strings.Builder
	shown := 0
	for _, ln := range lines[first:last] {
		for i := ln[0]; i < ln[1]; i++ {
			b.WriteString(m.renderChar(i))
		}
		b.WriteString("\n")
		shown++
	}
	for ; shown < visibleLines; shown++ {
		b.WriteString("\n")
	}
	return b.String()
}

func (m testScreen) View() string {
	status := dimStyle.Render("start typing to begin")
	if m.session.Started() {
		if m.md.Kind == mode.Timed {
			status = timerStyle.Render(fmt.Sprintf("%ds", m.remaining))
		} else {
			status = timerStyle.Render(fmt.Sprintf("%d/%d", m.wordsDone(), m.md.Words))
		}
	}

	content := dimStyle.Render(m.md.Title()) + "\n" +
		status + "\n\n" +
		m.textView() + "\n" +
		dimStyle.Render("tab restart • esc menu • ctrl+c quit")

	return place(m.width, m.height, m.textWidth(), content)
}
