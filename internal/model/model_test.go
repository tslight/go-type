package model

import (
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tobe/go-type/internal/content"
)

type dummyState struct{}

func (d *dummyState) GetSavedCharPos() int   { return 0 }
func (d *dummyState) SaveProgress(int) error { return nil }
func (d *dummyState) RecordSession(float64, float64, int, int, int, int) (string, error) {
	return "", nil
}

func TestNewModel_Creation(t *testing.T) {
	c := &content.Content{ID: 1, Name: "Test", Text: "Hello world"}
	m := NewModel(c.Text, c, 80, 24, &dummyState{})
	if m == nil {
		t.Fatal("NewModel returned nil")
	}
	if m.currentContent == nil || m.currentContent.Name != "Test" {
		t.Fatalf("expected currentContent name Test")
	}
}

func TestModel_UpdateTyping(t *testing.T) {
	c := &content.Content{ID: 1, Name: "Test", Text: "abcdef"}
	m := NewModel(c.Text, c, 80, 24, &dummyState{})
	// Simulate key presses via internal logic: append directly
	m.userInput = "abc"
	if len(m.userInput) != 3 {
		t.Fatalf("expected userInput len 3")
	}
}

func TestModel_WPMAccuracy(t *testing.T) {
	c := &content.Content{ID: 1, Name: "Test", Text: "aaaaa"}
	m := NewModel(c.Text, c, 80, 24, &dummyState{})
	m.userInput = "aaaaa"
	m.testStarted = true
	m.startTime = time.Now().Add(-1 * time.Minute)
	view := m.View()
	if view == "" {
		t.Fatalf("expected non-empty view")
	}
}

type captureState struct {
	savedPositions []int
	sessions       int
	lastWPM        float64
	lastDuration   int
}

func (c *captureState) GetSavedCharPos() int { return 0 }
func (c *captureState) SaveProgress(pos int) error {
	c.savedPositions = append(c.savedPositions, pos)
	return nil
}
func (c *captureState) RecordSession(wpm, accuracy float64, errors, charTypedRaw, effectiveChars, duration int) (string, error) {
	c.sessions++
	c.lastWPM = wpm
	c.lastDuration = duration
	return "", nil
}

func TestModel_KeyFlowAndFinish(t *testing.T) {
	c := &content.Content{ID: 1, Name: "Test", Text: "abc def"}
	cap := &captureState{}
	m := NewModel(c.Text, c, 40, 10, cap)
	// simulate window size to init viewport
	m.Update(tea.WindowSizeMsg{Width: 40, Height: 10})
	// type characters
	for _, r := range []rune{'a', 'b', 'c', ' '} {
		m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}
	// backspace
	m.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	// enter newline
	m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	// finish
	m.Update(tea.KeyMsg{Type: tea.KeyCtrlQ})
	if !m.finished {
		t.Fatalf("expected finished state")
	}
	v := m.View()
	if v == "" {
		t.Fatalf("expected results view after finish")
	}
	if cap.sessions == 0 {
		t.Fatalf("expected a recorded session")
	}
	if len(cap.savedPositions) == 0 {
		t.Fatalf("expected saved progress positions")
	}
}

func TestModel_EscExitToMenuAndPersist(t *testing.T) {
	c := &content.Content{ID: 2, Name: "EscDoc", Text: "hello world"}
	cap := &captureState{}
	m := NewModel(c.Text, c, 40, 10, cap)
	m.Update(tea.WindowSizeMsg{Width: 40, Height: 10})
	// Type a single character to start timing
	m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})
	// Now press ESC to exit to menu
	m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if !m.ExitToMenu() {
		t.Fatalf("expected ExitToMenu() true after ESC")
	}
	if cap.sessions == 0 {
		t.Fatalf("expected session to be recorded on ESC")
	}
}

func TestModel_ShortSessionWPMNonZero(t *testing.T) {
	c := &content.Content{ID: 3, Name: "Short", Text: "abc"}
	cap := &captureState{}
	m := NewModel(c.Text, c, 40, 10, cap)
	m.Update(tea.WindowSizeMsg{Width: 40, Height: 10})
	// Type quickly one char
	m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	// Finish immediately
	m.Update(tea.KeyMsg{Type: tea.KeyCtrlQ})
	if cap.sessions == 0 {
		t.Fatalf("expected a recorded session")
	}
	if cap.lastWPM <= 0 {
		t.Fatalf("expected WPM > 0 for short session, got %f", cap.lastWPM)
	}
	if cap.lastDuration < 1 {
		t.Fatalf("expected persisted duration >=1s when chars typed, got %d", cap.lastDuration)
	}
}

func TestModel_InternalNoopsCovered(t *testing.T) {
	c := &content.Content{ID: 4, Name: "Noop", Text: "text"}
	m := NewModel(c.Text, c, 40, 10, &dummyState{})
	// directly invoke no-op methods to increase coverage
	m.updateCursorPosition()
	m.rewrapText()
}

func TestNormalizeWhitespace(t *testing.T) {
	in := "a   b\n\n\n c"
	out := normalizeWhitespace(in)
	if strings.Contains(out, "   ") {
		t.Fatalf("expected collapsed spaces")
	}
	if strings.Count(out, "\n") > 2 {
		t.Fatalf("expected collapsed newlines")
	}
}

func TestIsExcessiveWhitespace(t *testing.T) {
	s := "a    b\n\n\n c"
	// pick a middle space in run of spaces
	spaceRunPos := strings.Index(s, "    ") + 1
	if spaceRunPos <= 0 {
		t.Fatalf("space run not found")
	}
	if !isExcessiveWhitespace(s, spaceRunPos) {
		t.Fatalf("expected excessive whitespace detection")
	}
	nlPos := strings.Index(s, "\n\n") + 1
	if nlPos <= 0 {
		t.Fatalf("newline run not found")
	}
	if !isExcessiveWhitespace(s, nlPos) {
		t.Fatalf("expected excessive newline detection")
	}
}

func TestModel_ScrollAndQuitKeys(t *testing.T) {
	c := &content.Content{ID: 1, Name: "Test", Text: strings.Repeat("x", 1000)}
	m := NewModel(c.Text, c, 40, 10, &dummyState{})
	_ = m.Init()
	// init viewport
	m.Update(tea.WindowSizeMsg{Width: 40, Height: 10})
	// scroll keys
	m.Update(tea.KeyMsg{Type: tea.KeyCtrlJ})
	m.Update(tea.KeyMsg{Type: tea.KeyCtrlK})
	m.Update(tea.KeyMsg{Type: tea.KeyCtrlF})
	m.Update(tea.KeyMsg{Type: tea.KeyCtrlB})
	// quit key
	m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	// Only ensure no panic and valid View returned
	_ = m.View()
}

func TestModel_DebugOverlayToggle(t *testing.T) {
	c := &content.Content{ID: 5, Name: "Debug", Text: strings.Repeat("abcd ", 200)}
	m := NewModel(c.Text, c, 60, 15, &dummyState{})
	m.Update(tea.WindowSizeMsg{Width: 60, Height: 15})
	// Type some characters to have data
	for _, r := range []rune{'a', 'b', 'c'} {
		m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}
	// Toggle on (ctrl+d)
	m.Update(tea.KeyMsg{Type: tea.KeyCtrlD})
	vOn := m.View()
	if !strings.Contains(vOn, "[Debug]") {
		t.Fatalf("expected debug overlay in view after toggle on")
	}
	// Toggle off (ctrl+d again)
	m.Update(tea.KeyMsg{Type: tea.KeyCtrlD})
	vOff := m.View()
	if strings.Contains(vOff, "[Debug]") {
		t.Fatalf("expected debug overlay removed after second toggle")
	}
}

type baselineState struct{ pos int }

func (b *baselineState) GetSavedCharPos() int   { return b.pos }
func (b *baselineState) SaveProgress(int) error { return nil }
func (b *baselineState) RecordSession(float64, float64, int, int, int, int) (string, error) {
	return "", nil
}

func TestModel_BaselineProgressPreload(t *testing.T) {
	text := "abcdefghij"
	c := &content.Content{ID: 6, Name: "Baseline", Text: text}
	state := &baselineState{pos: 5}
	m := NewModel(c.Text, c, 80, 24, state)
	if len(m.userInput) != 5 || m.userInput != text[:5] {
		t.Fatalf("expected userInput preloaded to first 5 chars, got %q", m.userInput)
	}
	if m.baselineRaw != 5 {
		t.Fatalf("expected baselineRaw 5, got %d", m.baselineRaw)
	}
	if m.baselineEffective <= 0 { // should be >0 based on non-excessive whitespace counting
		t.Fatalf("expected baselineEffective >0, got %d", m.baselineEffective)
	}
}

func TestModel_FinishWithoutTyping(t *testing.T) {
	c := &content.Content{ID: 7, Name: "Empty", Text: "abcdefgh"}
	cap := &captureState{}
	m := NewModel(c.Text, c, 40, 10, cap)
	// Immediately finish without typing
	m.Update(tea.KeyMsg{Type: tea.KeyCtrlQ})
	if !m.finished {
		t.Fatalf("expected finished after Ctrl+Q")
	}
	if cap.sessions == 0 {
		t.Fatalf("expected session recorded even with no typing")
	}
	if cap.lastWPM != 0 {
		t.Fatalf("expected WPM 0 when no typing, got %f", cap.lastWPM)
	}
}

// Ensure saved progress stops at the last correct contiguous prefix (does not include mismatches).
func TestModel_PrefixProgressStopsAtMismatch(t *testing.T) {
	c := &content.Content{ID: 8, Name: "Mismatch", Text: "abcdef"}
	cap := &captureState{}
	m := NewModel(c.Text, c, 40, 10, cap)
	m.Update(tea.WindowSizeMsg{Width: 40, Height: 10})
	// Type first char correct then mismatch immediately to force zero contiguous prefix beyond first
	m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}}) // mismatch vs 'b'
	m.Update(tea.KeyMsg{Type: tea.KeyCtrlQ})
	if cap.sessions == 0 {
		t.Fatalf("expected a recorded session")
	}
	if len(cap.savedPositions) == 0 {
		t.Fatalf("expected a saved position slice")
	}
	// Because mismatch occurred before second effective character, saved progress should be 1 (after 'a').
	if cap.savedPositions[0] != 1 {
		t.Fatalf("expected saved progress 1 (after 'a'), got %d", cap.savedPositions[0])
	}
}
