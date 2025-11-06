package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tobe/go-type/internal/textgen"
)

// Model represents the state of the typing test
type Model struct {
	text                  string
	userInput             string
	currentBook           *textgen.Book
	startTime             time.Time
	testStarted           bool
	finished              bool
	terminalWidth         int
	terminalHeight        int
	viewport              viewport.Model
	displayToNormPos      map[int]int // Cached mapping of display positions to normalized positions
	lastCachedTextLen     int         // Track when cache was last built
	cachedNormalizedText  string      // Cache of normalized text
	cachedNormalizedInput string      // Cache of normalized user input
	lastCachedInputLen    int         // Track when input cache was last built
	nonExcessiveInText    []int       // Cached: indices in text that are not excessive whitespace
	nonExcessiveInInput   []int       // Cached: indices in userInput that are not excessive whitespace
}

// Init initializes the model
func (m *Model) Init() tea.Cmd {
	return nil
}

// Update handles messages/input
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()

		// Handle Ctrl+Q to finish session
		if key == "ctrl+q" {
			m.finished = true
			return m, tea.Quit
		}

		// Handle Ctrl-C to quit without saving
		if key == "ctrl+c" {
			return m, tea.Quit
		}

		// Handle backspace
		if key == "backspace" && len(m.userInput) > 0 {
			m.userInput = m.userInput[:len(m.userInput)-1]
			m.updateCursorPosition()
			return m, nil
		}

		// Handle Enter key - add newline to userInput (typing test continues)
		if key == "enter" {
			m.userInput += "\n"
			m.updateCursorPosition()
			return m, nil
		}

		// Handle regular characters - also check if it's a rune
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

		// Initialize or update viewport
		if m.viewport.Width == 0 {
			m.viewport = viewport.New(msg.Width, msg.Height-3) // -3 for header
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

// View renders the UI
func (m *Model) View() string {
	if m.finished {
		return m.renderResults()
	}

	var b strings.Builder

	// Header takes 3 lines: newline + title + blank line
	headerText := fmt.Sprintf("\nOn your mark, get set, GO TYPE! (Source: %s)\nPress Ctrl+Q when done, Ctrl+C to quit\n\n", m.currentBook.Name)
	b.WriteString(headerText)

	// Render text character by character with validation
	var content strings.Builder

	// Cache normalized versions only when they change
	if len(m.text) != m.lastCachedTextLen {
		m.cachedNormalizedText = normalizeWhitespace(m.text)
		m.lastCachedTextLen = len(m.text)

		// Rebuild non-excessive indices for text
		m.nonExcessiveInText = make([]int, 0, len(m.text))
		for i := 0; i < len(m.text); i++ {
			if !isExcessiveWhitespace(m.text, i) {
				m.nonExcessiveInText = append(m.nonExcessiveInText, i)
			}
		}
	}

	if len(m.userInput) != m.lastCachedInputLen {
		m.cachedNormalizedInput = normalizeWhitespace(m.userInput)
		m.lastCachedInputLen = len(m.userInput)

		// Rebuild non-excessive indices for user input
		m.nonExcessiveInInput = make([]int, 0, len(m.userInput))
		for i := 0; i < len(m.userInput); i++ {
			if !isExcessiveWhitespace(m.userInput, i) {
				m.nonExcessiveInInput = append(m.nonExcessiveInInput, i)
			}
		}
	}

	// Build or update the display-to-normalized position map when text changes (including lazy loading)
	if m.displayToNormPos == nil || len(m.displayToNormPos) == 0 || m.lastCachedTextLen != len(m.text) {
		m.displayToNormPos = make(map[int]int)
		normalizedPos := 0
		for displayPos := 0; displayPos < len(m.text); displayPos++ {
			// Check if this is excessive whitespace (should be skipped in normalized version)
			if !isExcessiveWhitespace(m.text, displayPos) {
				m.displayToNormPos[displayPos] = normalizedPos
				normalizedPos++
			} else {
				m.displayToNormPos[displayPos] = -1 // Marker for excessive whitespace
			}
		}
		m.lastCachedTextLen = len(m.text)
	}

	// Render each display character - render enough to fill the viewport
	// Estimate: need at least viewport.Height * viewport.Width characters
	// Plus a buffer beyond user input for lookahead
	viewportSize := m.viewport.Height * m.viewport.Width
	if viewportSize < 500 {
		viewportSize = 500 // Minimum buffer
	}
	endPos := len(m.userInput) + viewportSize
	if endPos > len(m.text) {
		endPos = len(m.text)
	}

	for displayPos := 0; displayPos < endPos; displayPos++ {
		ch := m.text[displayPos]

		// Determine color based on validation
		var color string

		// Check if this character is excessive whitespace (don't require typing)
		if isExcessiveWhitespace(m.text, displayPos) {
			// Excessive whitespace - always show in gray (user doesn't type it)
			color = "\033[90m" // Gray
		} else {
			// This is a character user should type
			// Find which non-excessive character number this is
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
				// User has typed this character - check if it matches
				userCharPos := m.nonExcessiveInInput[textCharNum]
				if m.userInput[userCharPos] == ch {
					color = "\033[32m" // Green - correct
				} else {
					color = "\033[31m" // Red - incorrect
				}
			} else if textCharNum == len(m.nonExcessiveInInput) {
				// Cursor position
				color = "\033[4;33m" // Yellow underline
			} else {
				// Not yet typed
				color = "\033[90m" // Gray
			}
		}

		// Display the original character with its color (show spaces and tabs as-is, not as symbols)
		content.WriteString(fmt.Sprintf("%s%c\033[0m", color, ch))
	}

	// Set viewport content and render
	m.viewport.SetContent(content.String())
	b.WriteString(m.viewport.View())

	// Check if we need to load more text (lazy loading)
	// If user has typed past 80% of loaded text, load more
	if len(m.userInput) > int(float64(len(m.text))*0.8) {
		// Expand text by loading more from the book
		fullText := textgen.GetFullText()
		if len(fullText) > len(m.text) {
			// Load next chunk (add 50KB more)
			newEnd := len(m.text) + 50000
			if newEnd > len(fullText) {
				newEnd = len(fullText)
			}
			m.text = fullText[:newEnd]
		}
	}

	return b.String()
}

func (m *Model) updateCursorPosition() {
	// Viewport handles everything, no action needed
}

func (m *Model) rewrapText() {
	// Viewport handles wrapping, no action needed
}

func (m *Model) renderResults() string {
	var duration time.Duration
	if m.testStarted {
		duration = time.Since(m.startTime)
	}

	wpm := calculateWPM(m.userInput, duration)
	accuracy := calculateAccuracy(m.text, m.userInput)
	errors := calculateErrors(m.text, m.userInput)

	// Save progress - save the position based on non-excessive characters only
	// This ensures we restore to the exact spot the user typed to, skipping excessive whitespace
	if m.currentBook != nil {
		// Count how many non-excessive characters the user has typed
		nonExcessiveCount := len(m.nonExcessiveInInput)

		// Find the corresponding position in m.text
		// We need to find the position of the nonExcessiveCount-th non-excessive character in m.text
		charPos := 0
		if nonExcessiveCount > 0 && nonExcessiveCount <= len(m.nonExcessiveInText) {
			charPos = m.nonExcessiveInText[nonExcessiveCount-1] + 1
		}
		_ = textgen.SaveProgress(charPos, "")
	}

	return fmt.Sprintf("\n\nDuration: %.2f seconds\nWPM: %.2f\nAccuracy: %.2f%%\nErrors: %d\nTyped: %d/%d characters\nProgress saved!\n",
		duration.Seconds(), wpm, accuracy, errors, len(m.userInput), len(m.text))
}

// toASCII filters out non-ASCII characters to avoid UTF-8 encoding issues
// Preserves newlines for paragraph formatting
func toASCII(s string) string {
	var result []byte
	for i := 0; i < len(s); i++ {
		// Keep newlines and ASCII characters
		if s[i] == '\n' || (s[i] < 128 && s[i] >= 32) || s[i] == '\t' {
			result = append(result, s[i])
		}
	}
	return string(result)
}

// wrapTextManually wraps text at specified width while preserving paragraphs
func wrapTextManually(text string, width int) string {
	// Split by double newlines to preserve paragraph structure
	paragraphs := strings.Split(text, "\n\n")
	var result []string

	for _, para := range paragraphs {
		// Skip empty paragraphs
		para = strings.TrimSpace(para)
		if para == "" {
			result = append(result, "")
			continue
		}

		// Remove any internal line breaks in the paragraph (preserve as single paragraph)
		para = strings.ReplaceAll(para, "\n", " ")
		// Collapse multiple spaces into single space
		para = strings.Join(strings.Fields(para), " ")

		// Wrap this paragraph
		wrapped := wrapParagraph(para, width)
		result = append(result, wrapped)
		// Add blank line after paragraph (except last)
		result = append(result, "")
	}

	return strings.Join(result, "\n")
}

// NewModel creates a new typing test model
func NewModel(text string, book *textgen.Book, width, height int) *Model {
	m := &Model{
		text:           text, // Text is already ASCII-filtered from textgen
		currentBook:    book,
		terminalWidth:  width,
		terminalHeight: height,
		viewport:       viewport.New(width, height-3),
	}

	// On resume, pre-fill userInput with already-completed characters
	// The saved position is at the character AFTER the last non-excessive char typed
	savedCharPos := textgen.GetCurrentCharPos()
	if savedCharPos > 0 && savedCharPos <= len(m.text) {
		// We saved the position right after the last non-excessive char
		// So we can restore directly to that position
		m.userInput = m.text[:savedCharPos]
	}

	m.viewport.YPosition = 3
	return m
}

// wrapParagraph wraps a single paragraph at specified width
func wrapParagraph(text string, width int) string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return ""
	}

	var lines []string
	var currentLine string

	for _, word := range words {
		// If adding this word would exceed width, start a new line
		if len(currentLine) > 0 && len(currentLine)+1+len(word) > width {
			lines = append(lines, currentLine)
			currentLine = word
		} else {
			if len(currentLine) > 0 {
				currentLine += " "
			}
			currentLine += word
		}
	}

	if len(currentLine) > 0 {
		lines = append(lines, currentLine)
	}

	return strings.Join(lines, "\n")
}

// normalizeWhitespace collapses excessive whitespace to make typing easier
// Collapses 3+ consecutive spaces/tabs into single space
// Collapses 2+ consecutive newlines into single newline
func normalizeWhitespace(s string) string {
	var result strings.Builder
	lastWasSpace := false
	lastWasNewline := false

	for i := 0; i < len(s); i++ {
		ch := s[i]

		if ch == '\n' {
			// Handle newline
			if !lastWasNewline {
				result.WriteByte('\n')
				lastWasNewline = true
				lastWasSpace = false
			}
			// Skip this newline if we just wrote one
		} else if ch == ' ' || ch == '\t' {
			// Handle spaces and tabs
			if !lastWasSpace {
				result.WriteByte(' ')
				lastWasSpace = true
				lastWasNewline = false
			}
			// Skip this space/tab if we just wrote one
		} else {
			// Regular character
			result.WriteByte(ch)
			lastWasSpace = false
			lastWasNewline = false
		}
	}

	return result.String()
}

// isExcessiveWhitespace checks if the character at position i in string s
// is part of excessive whitespace (3+ spaces/tabs or 2+ newlines in a row)
func isExcessiveWhitespace(s string, pos int) bool {
	if pos >= len(s) {
		return false
	}

	ch := s[pos]
	if ch != ' ' && ch != '\t' && ch != '\n' {
		return false // Not whitespace
	}

	// Count how many of the same whitespace type are consecutive
	if ch == '\n' {
		// Excessive newlines: 2 or more
		if pos > 0 && s[pos-1] == '\n' {
			return true
		}
		if pos < len(s)-1 && s[pos+1] == '\n' {
			return true
		}
		return false
	}

	// For spaces and tabs, count consecutive
	// Excessive: 3+ spaces/tabs
	count := 1
	// Count backwards
	for i := pos - 1; i >= 0 && (s[i] == ' ' || s[i] == '\t'); i-- {
		count++
	}
	// Count forwards
	for i := pos + 1; i < len(s) && (s[i] == ' ' || s[i] == '\t'); i++ {
		count++
	}

	return count >= 3
}
