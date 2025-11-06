package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tobe/go-type/internal/textgen"
)

// MenuModel represents the book selection menu state
type MenuModel struct {
	books          []textgen.Book
	selectedIndex  int
	viewport       viewport.Model
	terminalWidth  int
	terminalHeight int
	selectedBook   *textgen.Book
	done           bool
}

// NewMenuModel creates a new book selection menu
func NewMenuModel(width, height int) *MenuModel {
	books := textgen.GetAvailableBooks()
	m := &MenuModel{
		books:          books,
		selectedIndex:  0,
		terminalWidth:  width,
		terminalHeight: height,
		viewport:       viewport.New(width, height-4),
	}
	m.viewport.YPosition = 3
	m.renderMenu()
	return m
}

// Init initializes the menu
func (m *MenuModel) Init() tea.Cmd {
	return nil
}

// Update handles input
func (m *MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()

		switch key {
		case "j", "down":
			// Move down
			if m.selectedIndex < len(m.books)-1 {
				m.selectedIndex++
				m.syncViewport()
			}
		case "k", "up":
			// Move up
			if m.selectedIndex > 0 {
				m.selectedIndex--
				m.syncViewport()
			}
		case "g":
			// Go to start (Vi style)
			m.selectedIndex = 0
			m.syncViewport()
		case "G":
			// Go to end (Vi style)
			m.selectedIndex = len(m.books) - 1
			m.syncViewport()
		case "enter":
			// Select book
			m.selectedBook = &m.books[m.selectedIndex]
			m.done = true
			return m, tea.Quit
		case "q", "esc", "ctrl+c":
			// Quit without selecting
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.terminalWidth = msg.Width
		m.terminalHeight = msg.Height
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - 4
		m.renderMenu()
		m.syncViewport()
	}

	return m, nil
}

// View renders the menu
func (m *MenuModel) View() string {
	var b strings.Builder

	// Header
	headerText := "\nSelect a book (j/k to navigate, Enter to select, q to quit)\n\n"
	b.WriteString(headerText)

	// Books list
	var content strings.Builder
	for i, book := range m.books {
		if i == m.selectedIndex {
			// Highlight selected book
			content.WriteString(fmt.Sprintf("\033[1;33m▶ %3d: %s\033[0m\n", book.ID, book.Name))
		} else {
			content.WriteString(fmt.Sprintf("  %3d: %s\n", book.ID, book.Name))
		}
	}

	m.viewport.SetContent(content.String())
	b.WriteString(m.viewport.View())

	return b.String()
}

// SelectedBook returns the selected book (if any)
func (m *MenuModel) SelectedBook() *textgen.Book {
	return m.selectedBook
}

// renderMenu updates the viewport content
func (m *MenuModel) renderMenu() {
	var content strings.Builder
	for i, book := range m.books {
		if i == m.selectedIndex {
			// Highlight selected book with yellow background and arrow
			content.WriteString(fmt.Sprintf("\033[1;33m▶ %3d: %s\033[0m\n", book.ID, book.Name))
		} else {
			content.WriteString(fmt.Sprintf("  %3d: %s\n", book.ID, book.Name))
		}
	}
	m.viewport.SetContent(content.String())
}

// syncViewport ensures the selected item is visible in the viewport
func (m *MenuModel) syncViewport() {
	m.renderMenu()

	// Calculate which line the selected item is on
	// Each book is one line
	selectedLine := m.selectedIndex

	// Ensure the selected line is visible in the viewport
	if selectedLine < m.viewport.YOffset {
		// Selected item is above visible area, scroll up
		m.viewport.YOffset = selectedLine
	} else if selectedLine >= m.viewport.YOffset+m.viewport.Height {
		// Selected item is below visible area, scroll down
		m.viewport.YOffset = selectedLine - m.viewport.Height + 1
	}
}
