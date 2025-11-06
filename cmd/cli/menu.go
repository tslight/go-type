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
	books           []textgen.Book
	selectedIndex   int
	viewport        viewport.Model
	terminalWidth   int
	terminalHeight  int
	selectedBook    *textgen.Book
	done            bool
	searchMode      bool
	searchQuery     string
	searchDirection int // 1 for forward (/), -1 for backward (?)
	searchResults   []int
	searchIndex     int
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

		// Handle search mode input separately
		if m.searchMode {
			switch key {
			case "enter":
				// Execute search and exit search mode
				m.performSearch()
				m.searchMode = false
				m.renderMenu()
			case "backspace":
				if len(m.searchQuery) > 0 {
					m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
				}
			case "esc":
				// Exit search mode without selecting
				m.searchMode = false
				m.searchQuery = ""
				m.searchResults = nil
				m.renderMenu()
			default:
				// Add character to search query
				if len(key) == 1 && key[0] >= 32 && key[0] < 127 {
					m.searchQuery += key
				}
			}
			return m, nil
		}

		// Normal navigation mode
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
		case "n":
			// Next search result
			if len(m.searchResults) > 0 {
				m.searchIndex = (m.searchIndex + 1) % len(m.searchResults)
				m.selectedIndex = m.searchResults[m.searchIndex]
				m.syncViewport()
			}
		case "N":
			// Previous search result
			if len(m.searchResults) > 0 {
				m.searchIndex = (m.searchIndex - 1 + len(m.searchResults)) % len(m.searchResults)
				m.selectedIndex = m.searchResults[m.searchIndex]
				m.syncViewport()
			}
		case "/":
			// Forward search
			m.searchMode = true
			m.searchQuery = ""
			m.searchDirection = 1
		case "?":
			// Backward search
			m.searchMode = true
			m.searchQuery = ""
			m.searchDirection = -1
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
		if m.searchMode {
			m.viewport.Height = msg.Height - 6 // Extra line for search input
		}
		m.renderMenu()
		m.syncViewport()
	}

	return m, nil
}

// View renders the menu
func (m *MenuModel) View() string {
	var b strings.Builder

	// Header
	var headerText string
	if m.searchMode {
		searchPrefix := "/"
		if m.searchDirection == -1 {
			searchPrefix = "?"
		}
		headerText = fmt.Sprintf("\nSelect a book (searching... Press Enter to search, Esc to cancel)\n%s%s\n\n", searchPrefix, m.searchQuery)
	} else {
		headerText = "\nSelect a book (j/k navigate, / search, n/N next/prev result, Enter select, q quit)\n\n"
	}
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

// performSearch searches for books matching the query
func (m *MenuModel) performSearch() {
	if m.searchQuery == "" {
		m.searchResults = nil
		return
	}

	query := strings.ToLower(m.searchQuery)
	m.searchResults = nil

	// Search through all books
	for i, book := range m.books {
		if strings.Contains(strings.ToLower(book.Name), query) {
			m.searchResults = append(m.searchResults, i)
		}
	}

	// If results found, select the first one
	if len(m.searchResults) > 0 {
		m.searchIndex = 0
		m.selectedIndex = m.searchResults[0]
		m.syncViewport()
	}
}
