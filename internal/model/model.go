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
}

// SessionState is the minimal persistence interface Model needs.
type SessionState interface {
	GetSavedCharPos() int
	SaveProgress(charPos int) error
	RecordSession(wpm, accuracy float64, errors, charTyped, duration int) (string, error)
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
		if key == "ctrl+q" || key == "ctrl+s" {
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
			m.userInput = m.userInput[:len(m.userInput)-1]
			m.updateCursorPosition()
			return m, nil
		}
		if key == "enter" {
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
	var duration time.Duration
	if m.testStarted {
		duration = time.Since(m.startTime)
	}
	wpm := utils.CalculateWPM(m.userInput, duration)
	accuracy := utils.CalculateAccuracy(m.text, m.userInput)
	errors := utils.CalculateErrors(m.text, m.userInput)
	sessionStats := ""
	if m.stateProvider != nil {
		nonExcessiveCount := len(m.nonExcessiveInInput)
		charPos := 0
		if nonExcessiveCount > 0 && nonExcessiveCount <= len(m.nonExcessiveInText) {
			charPos = m.nonExcessiveInText[nonExcessiveCount-1] + 1
		}
		if err := m.stateProvider.SaveProgress(charPos); err == nil {
			if stats, err := m.stateProvider.RecordSession(wpm, accuracy, errors, len(m.userInput), int(duration.Seconds())); err == nil {
				sessionStats = stats
			}
		}
	}
	currentSessionStr := fmt.Sprintf("Duration: %.2f seconds\nWPM: %.2f\nAccuracy: %.2f%%\nErrors: %d\nTyped: %d/%d characters\nProgress saved!",
		duration.Seconds(), wpm, accuracy, errors, len(m.userInput), len(m.text))
	return "\n\n" + currentSessionStr + sessionStats + "\n\nPress any key to continue...\n"
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
	m.viewport.YPosition = 3
	return m
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
