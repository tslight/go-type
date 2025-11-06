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
	text           string
	userInput      string
	currentBook    *textgen.Book
	sentenceCount  int // Number of sentences in the current paragraph
	startTime      time.Time
	testStarted    bool
	finished       bool
	cursorX        int // Cursor X position for wrapping
	cursorY        int // Cursor Y position for wrapping
	wrappedLines   []string
	terminalWidth  int
	terminalHeight int
	viewport       viewport.Model
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

		// Handle Enter to finish
		if key == "enter" {
			m.finished = true
			return m, tea.Quit
		}

		// Handle Ctrl-C
		if key == "ctrl+c" {
			return m, tea.Quit
		}

		// Handle backspace
		if key == "backspace" && len(m.userInput) > 0 {
			m.userInput = m.userInput[:len(m.userInput)-1]
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
	headerText := fmt.Sprintf("\nOn your mark, get set, GO TYPE! (Source: %s)\n\n", m.currentBook.Name)
	b.WriteString(headerText)

	// Build the full content to display
	var content strings.Builder
	for i := 0; i < len(m.wrappedLines); i++ {
		line := m.wrappedLines[i]

		// Calculate position of this line in the text
		lineStart := 0
		for j := 0; j < i; j++ {
			lineStart += len(m.wrappedLines[j]) + 1 // +1 for the newline between lines in original text
		}

		// Build the displayed line
		var displayLine strings.Builder
		for j, ch := range line {
			pos := lineStart + j

			if pos < len(m.userInput) {
				// Character has been typed
				expectedChar := m.text[pos]
				if m.userInput[pos] == expectedChar {
					displayLine.WriteString(fmt.Sprintf("\033[32m%c\033[0m", expectedChar)) // Green
				} else {
					displayLine.WriteString(fmt.Sprintf("\033[31m%c\033[0m", expectedChar)) // Red
				}
			} else if pos == len(m.userInput) && pos < len(m.text) {
				// Cursor position - show next character with underline
				displayLine.WriteString(fmt.Sprintf("\033[4;33m%c\033[0m", ch)) // Yellow underline
			} else if pos < len(m.text) {
				// Character not typed yet - show in gray
				displayLine.WriteString(fmt.Sprintf("\033[90m%c\033[0m", ch)) // Gray
			} else {
				// Beyond text end
				break
			}
		}

		content.WriteString(displayLine.String())
		content.WriteString("\n")
	}

	// Set viewport content and render
	m.viewport.SetContent(content.String())
	b.WriteString(m.viewport.View())

	return b.String()
} // Helper functions

func (m *Model) rewrapText() {
	if m.terminalWidth == 0 {
		m.terminalWidth = 80
	}
	wrappedText := wrapTextManually(m.text, m.terminalWidth)
	m.wrappedLines = strings.Split(wrappedText, "\n")
}

func (m *Model) updateCursorPosition() {
	// Calculate cursor position based on userInput length
	pos := len(m.userInput)
	currentPos := 0

	m.cursorY = 0
	m.cursorX = 0

	for i, line := range m.wrappedLines {
		lineLen := len(line) + 1 // +1 for the newline separator between lines
		if currentPos+lineLen > pos {
			// Cursor is on this line
			m.cursorY = i
			m.cursorX = pos - currentPos
			break
		}
		currentPos += lineLen
	}
}

func (m *Model) renderResults() string {
	var duration time.Duration
	if m.testStarted {
		duration = time.Since(m.startTime)
	}

	wpm := calculateWPM(m.userInput, duration)
	accuracy := calculateAccuracy(m.text, m.userInput)
	errors := calculateErrors(m.text, m.userInput)

	// Save progress - absolute position of where user completed typing
	if m.currentBook != nil {
		charPos := len(m.userInput)
		_ = textgen.SaveProgress(charPos, "")
	}

	return fmt.Sprintf("\n\nDuration: %.2f seconds\nWPM: %.2f\nAccuracy: %.2f%%\nErrors: %d\nTyped: %d/%d characters\nProgress saved!\n",
		duration.Seconds(), wpm, accuracy, errors, len(m.userInput), len(m.text))
}

// NewModel creates a new typing test model
func NewModel(text string, book *textgen.Book, sentenceCount, width, height int) *Model {
	m := &Model{
		text:           toASCII(text),
		currentBook:    book,
		sentenceCount:  sentenceCount,
		terminalWidth:  width,
		terminalHeight: height,
		viewport:       viewport.New(width, height-3),
	}

	// On resume, pre-fill userInput with already-completed characters
	// This will show them as "typed" (grayed out) so user sees what they've completed
	savedCharPos := textgen.GetCurrentCharPos()
	if savedCharPos > 0 && savedCharPos <= len(text) {
		m.userInput = text[:savedCharPos]
	}

	m.viewport.YPosition = 3 // Position below header
	m.rewrapText()
	return m
}

// toASCII filters out non-ASCII characters to avoid UTF-8 encoding issues
func toASCII(s string) string {
	var result []byte
	for i := 0; i < len(s); i++ {
		if s[i] < 128 {
			result = append(result, s[i])
		}
	}
	return string(result)
}

// wrapTextManually wraps text at specified width by breaking words
func wrapTextManually(text string, width int) string {
	words := strings.Fields(text)
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
