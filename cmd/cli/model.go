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
	startTime      time.Time
	testStarted    bool
	finished       bool
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
	for i := 0; i < len(m.text); i++ {
		ch := m.text[i]
		if i < len(m.userInput) {
			// Character has been typed - check if it matches
			if m.userInput[i] == ch {
				content.WriteString(fmt.Sprintf("\033[32m%c\033[0m", ch)) // Green
			} else {
				content.WriteString(fmt.Sprintf("\033[31m%c\033[0m", ch)) // Red
			}
		} else if i == len(m.userInput) && i < len(m.text) {
			// Cursor position - show character with underline
			content.WriteString(fmt.Sprintf("\033[4;33m%c\033[0m", ch)) // Yellow underline
		} else if i < len(m.text) {
			// Not yet typed - show in gray
			content.WriteString(fmt.Sprintf("\033[90m%c\033[0m", ch)) // Gray
		} else {
			// Beyond text
			break
		}
	}

	// Set viewport content and render
	m.viewport.SetContent(content.String())
	b.WriteString(m.viewport.View())

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

	// Save progress - absolute position of where user completed typing
	if m.currentBook != nil {
		charPos := len(m.userInput)
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
		text:           toASCII(text),
		currentBook:    book,
		terminalWidth:  width,
		terminalHeight: height,
		viewport:       viewport.New(width, height-3),
	}

	// On resume, pre-fill userInput with already-completed characters
	savedCharPos := textgen.GetCurrentCharPos()
	if savedCharPos > 0 && savedCharPos <= len(text) {
		m.userInput = text[:savedCharPos]
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
