package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tobe/go-type/internal/content"
	"github.com/tobe/go-type/internal/utils"
)

// Model represents the state of the typing test
type Model struct {
	text                  string
	userInput             string
	currentContent        *content.Content
	stateProvider         SessionState
	startTime             time.Time
	testStarted           bool
	finished              bool
	terminalWidth         int
	terminalHeight        int
	viewport              viewport.Model
	displayToNormPos      map[int]int
	lastCachedTextLen     int
	cachedNormalizedText  string
	cachedNormalizedInput string
	lastCachedInputLen    int
	nonExcessiveInText    []int
	nonExcessiveInInput   []int
	cachedRenderedText    string
	lastRenderedInputLen  int
	baselineRaw           int
	baselineEffective     int
	exitToMenu            bool
	suppressResults       bool
	sessionPersisted      bool
	cachedResultsString   string
}

// SessionState is the minimal persistence interface Model needs.
type SessionState interface {
	GetSavedCharPos() int
	SaveProgress(charPos int) error
	RecordSession(wpm, accuracy float64, errors, charTypedRaw, effectiveChars, duration int) (string, error)
}

func (m *Model) Init() tea.Cmd { return nil }

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		if m.finished {
			if key == "ctrl+c" {
				return m, tea.Quit
			}
			return m, tea.Quit
		}
		// ESC: finalize session, suppress results view, and signal runner to return to menu.
		if key == "esc" {
			m.finalizeSession()
			m.exitToMenu = true
			m.suppressResults = true
			m.finished = true
			return m, tea.Quit
		}
		if key == "ctrl+q" || key == "ctrl+s" {
			m.finalizeSession()
			m.finished = true
			return m, tea.Quit
		}
		if key == "ctrl+c" || key == "ctrl+d" {
			return m, tea.Quit
		}
		if key == "ctrl+j" {
			m.viewport.ScrollDown(1)
			return m, nil
		}
		if key == "ctrl+k" {
			m.viewport.ScrollUp(1)
			return m, nil
		}
		if key == "ctrl+f" {
			m.viewport.PageDown()
			return m, nil
		}
		if key == "ctrl+b" {
			m.viewport.PageUp()
			return m, nil
		}
		if key == "backspace" && len(m.userInput) > 0 {
			if !m.testStarted {
				m.testStarted = true
				m.startTime = time.Now()
			}
			m.userInput = m.userInput[:len(m.userInput)-1]
			m.updateCursorPosition()
			return m, nil
		}
		if key == "enter" {
			if !m.testStarted {
				m.testStarted = true
				m.startTime = time.Now()
			}
			m.userInput += "\n"
			m.updateCursorPosition()
			return m, nil
		}
		if len(key) == 1 && key[0] >= 32 && key[0] < 127 {
			if !m.testStarted {
				m.testStarted = true
				m.startTime = time.Now()
			}
			m.userInput += key
			m.updateCursorPosition()
			return m, nil
		}
	case tea.WindowSizeMsg:
		m.terminalWidth = msg.Width
		m.terminalHeight = msg.Height
		if m.viewport.Width == 0 {
			m.viewport = viewport.New(msg.Width, msg.Height-3)
			m.viewport.YPosition = 3
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - 3
		}
		m.rewrapText()
		return m, nil
	}
	return m, nil
}

func (m *Model) View() string {
	if m.finished {
		if m.suppressResults {
			return ""
		}
		return m.renderResults()
	}
	var b strings.Builder
	sourceName := "Unknown Source"
	if m.currentContent != nil {
		sourceName = m.currentContent.Name
	}
	b.WriteString(fmt.Sprintf("\nOn your mark, get set, GO TYPE! (Source: %s)\nPress Ctrl+Q or Ctrl+S when done, Ctrl+C to quit\n\n", sourceName))

	if len(m.text) != m.lastCachedTextLen {
		m.cachedNormalizedText = normalizeWhitespace(m.text)
		m.lastCachedTextLen = len(m.text)
		viewportSize := m.viewport.Height * m.viewport.Width
		if viewportSize < 500 {
			viewportSize = 500
		}
		renderUpTo := len(m.userInput) + (viewportSize * 2)
		if renderUpTo > len(m.text) {
			renderUpTo = len(m.text)
		}
		m.nonExcessiveInText = make([]int, 0, renderUpTo)
		for i := 0; i < renderUpTo; i++ {
			if !isExcessiveWhitespace(m.text, i) {
				m.nonExcessiveInText = append(m.nonExcessiveInText, i)
			}
		}
	}
	if len(m.userInput) != m.lastCachedInputLen {
		m.cachedNormalizedInput = normalizeWhitespace(m.userInput)
		m.lastCachedInputLen = len(m.userInput)
		m.nonExcessiveInInput = make([]int, 0, len(m.userInput))
		for i := 0; i < len(m.userInput); i++ {
			if !isExcessiveWhitespace(m.userInput, i) {
				m.nonExcessiveInInput = append(m.nonExcessiveInInput, i)
			}
		}
	}
	if m.displayToNormPos == nil || (len(m.displayToNormPos) == 0 && len(m.text) > 0) {
		m.displayToNormPos = make(map[int]int)
		viewportSize := m.viewport.Height * m.viewport.Width
		if viewportSize < 500 {
			viewportSize = 500
		}
		renderUpTo := len(m.userInput) + (viewportSize * 10)
		if renderUpTo > len(m.text) {
			renderUpTo = len(m.text)
		}
		normalizedPos := 0
		for displayPos := 0; displayPos < renderUpTo; displayPos++ {
			if !isExcessiveWhitespace(m.text, displayPos) {
				m.displayToNormPos[displayPos] = normalizedPos
				normalizedPos++
			} else {
				m.displayToNormPos[displayPos] = -1
			}
		}
		m.lastCachedTextLen = len(m.text)
	}
	if len(m.userInput) != m.lastRenderedInputLen || m.cachedRenderedText == "" {
		var contentBuf strings.Builder
		viewportSize := m.viewport.Height * m.viewport.Width
		if viewportSize < 500 {
			viewportSize = 500
		}
		renderUpTo := len(m.userInput) + (viewportSize * 2)
		if renderUpTo > len(m.text) {
			renderUpTo = len(m.text)
		}
		for displayPos := 0; displayPos < renderUpTo; displayPos++ {
			ch := m.text[displayPos]
			var color string
			if isExcessiveWhitespace(m.text, displayPos) {
				color = "\033[90m"
			} else {
				textCharNum := -1
				for i, pos := range m.nonExcessiveInText {
					if pos == displayPos {
						textCharNum = i
						break
					}
					if pos > displayPos {
						break
					}
				}
				if textCharNum >= 0 && textCharNum < len(m.nonExcessiveInInput) {
					userCharPos := m.nonExcessiveInInput[textCharNum]
					if m.userInput[userCharPos] == ch {
						color = "\033[32m"
					} else {
						color = "\033[31m"
					}
				} else if textCharNum == len(m.nonExcessiveInInput) {
					color = "\033[4;33m"
				} else {
					color = "\033[90m"
				}
			}
			contentBuf.WriteString(fmt.Sprintf("%s%c\033[0m", color, ch))
		}
		m.cachedRenderedText = contentBuf.String()
		m.lastRenderedInputLen = len(m.userInput)
	}
	m.viewport.SetContent(m.cachedRenderedText)
	b.WriteString(m.viewport.View())
	return b.String()
}

func (m *Model) updateCursorPosition() {}
func (m *Model) rewrapText()           {}

func (m *Model) renderResults() string {
	// Ensure session persisted and cached string computed
	m.finalizeSession()
	return "\n\n" + m.cachedResultsString + "\n\nPress any key to continue...\n"
}

func NewModel(text string, contentItem *content.Content, width, height int, provider SessionState) *Model {
	m := &Model{text: text, currentContent: contentItem, stateProvider: provider, terminalWidth: width, terminalHeight: height, viewport: viewport.New(width, height-3)}
	savedCharPos := 0
	if provider != nil {
		savedCharPos = provider.GetSavedCharPos()
	}
	if savedCharPos > 0 && savedCharPos <= len(m.text) {
		m.userInput = m.text[:savedCharPos]
	}
	// Establish baselines based on prefilled progress
	m.baselineRaw = len(m.userInput)
	if savedCharPos > 0 {
		eff := 0
		for i := 0; i < savedCharPos && i < len(m.text); i++ {
			if !isExcessiveWhitespace(m.text, i) {
				eff++
			}
		}
		m.baselineEffective = eff
	}
	m.viewport.YPosition = 3
	return m
}

// ExitToMenu indicates the session ended with a request to return to menu (via ESC).
func (m *Model) ExitToMenu() bool { return m.exitToMenu }

// finalizeSession persists session stats exactly once and caches a results string for rendering.
func (m *Model) finalizeSession() {
	if m.sessionPersisted {
		return
	}
	var duration time.Duration
	if m.testStarted {
		duration = time.Since(m.startTime)
	}
	// Compute session deltas to avoid counting prefilled progress
	sessionRaw := len(m.userInput) - m.baselineRaw
	if sessionRaw < 0 {
		sessionRaw = 0
	}
	sessionEffective := len(m.nonExcessiveInInput) - m.baselineEffective
	if sessionEffective < 0 {
		sessionEffective = 0
	}
	// Use a minimal 1s duration for WPM if we typed but duration is <1s
	adjDuration := duration
	if sessionRaw > 0 && adjDuration < time.Second {
		adjDuration = time.Second
	}
	wpm := utils.CalculateWPM(m.userInput[m.baselineRaw:], adjDuration)
	// Build effective strings for accuracy/errors
	var effInputBuilder strings.Builder
	for _, pos := range m.nonExcessiveInInput {
		if pos >= 0 && pos < len(m.userInput) {
			effInputBuilder.WriteByte(m.userInput[pos])
		}
	}
	effInput := effInputBuilder.String()
	var effTextBuilder strings.Builder
	for i, pos := range m.nonExcessiveInText {
		if i >= len(effInput) {
			break
		}
		if pos >= 0 && pos < len(m.text) {
			effTextBuilder.WriteByte(m.text[pos])
		}
	}
	effText := effTextBuilder.String()
	accuracy := utils.CalculateAccuracy(effText, effInput)
	errors := utils.CalculateErrors(effText, effInput)

	// Determine current effective progress position (filtered)
	nonExcessiveCount := len(m.nonExcessiveInInput)
	charPos := 0
	if nonExcessiveCount > 0 && nonExcessiveCount <= len(m.nonExcessiveInText) {
		charPos = m.nonExcessiveInText[nonExcessiveCount-1] + 1
	}

	sessionStats := ""
	if m.stateProvider != nil {
		if err := m.stateProvider.SaveProgress(charPos); err == nil {
			// Round duration up to at least 1 second when we have typed characters to avoid zero-second sessions.
			durSec := int(duration.Seconds())
			if sessionRaw > 0 && durSec == 0 {
				durSec = 1
			}
			if stats, err := m.stateProvider.RecordSession(wpm, accuracy, errors, sessionRaw, sessionEffective, durSec); err == nil {
				sessionStats = stats
			}
		}
	}

	totalLen := len(m.text)
	displaySeconds := duration.Seconds()
	if sessionRaw > 0 && displaySeconds == 0 {
		displaySeconds = 1
	}
	currentSessionStr := fmt.Sprintf("Duration: %.2f seconds\nWPM: %.2f\nAccuracy: %.2f%%\nErrors: %d\nTyped this session (raw/eff): %d/%d\nText Progress: %d/%d\nProgress saved!",
		displaySeconds, wpm, accuracy, errors, sessionRaw, sessionEffective, charPos, totalLen)
	m.cachedResultsString = currentSessionStr + sessionStats
	m.sessionPersisted = true
}

func normalizeWhitespace(s string) string {
	var result strings.Builder
	lastWasSpace := false
	lastWasNewline := false
	for i := 0; i < len(s); i++ {
		ch := s[i]
		switch ch {
		case '\n':
			if !lastWasNewline {
				result.WriteByte('\n')
				lastWasNewline = true
				lastWasSpace = false
			}
		case ' ', '\t':
			if !lastWasSpace {
				result.WriteByte(' ')
				lastWasSpace = true
				lastWasNewline = false
			}
		default:
			result.WriteByte(ch)
			lastWasSpace = false
			lastWasNewline = false
		}
	}
	return result.String()
}

func isExcessiveWhitespace(s string, pos int) bool {
	if pos >= len(s) {
		return false
	}
	ch := s[pos]
	if ch != ' ' && ch != '\t' && ch != '\n' {
		return false
	}
	if ch == '\n' {
		if pos > 0 && s[pos-1] == '\n' {
			return true
		}
		if pos < len(s)-1 && s[pos+1] == '\n' {
			return true
		}
		return false
	}
	count := 1
	for i := pos - 1; i >= 0 && (s[i] == ' ' || s[i] == '\t'); i-- {
		count++
	}
	for i := pos + 1; i < len(s) && (s[i] == ' ' || s[i] == '\t'); i++ {
		count++
	}
	return count >= 3
}
