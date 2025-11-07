package cli

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
	showingStats    bool // True when displaying stats for a book
	statsBookID     int  // ID of book whose stats are being shown
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

		// Handle stats view - if showing stats, only escape key closes it
		if m.showingStats {
			switch key {
			case "esc", "i", "q":
				m.showingStats = false
				m.renderMenu()
			}
			return m, nil
		}

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
		case "i":
			// Show stats for selected book
			m.showingStats = true
			m.statsBookID = m.books[m.selectedIndex].ID
			m.renderMenu()
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

	// Show stats view if requested
	if m.showingStats {
		stats := textgen.GetBookStats(&m.books[m.selectedIndex])
		statsStr := textgen.FormatBookStats(stats)
		headerText := "\n\nBook: " + m.books[m.selectedIndex].Name + "\n"
		b.WriteString(headerText)
		b.WriteString(statsStr)
		b.WriteString("\nPress any key to continue...\n")
		return b.String()
	}

	// Header
	var headerText string
	if m.searchMode {
		searchPrefix := "/"
		if m.searchDirection == -1 {
			searchPrefix = "?"
		}
		headerText = fmt.Sprintf("\nSelect a book (searching... Press Enter to search, Esc to cancel)\n%s%s\n\n", searchPrefix, m.searchQuery)
	} else {
		headerText = "\nSelect a book (j/k navigate, / search, n/N next/prev result, i info, Enter select, q quit)\n\n"
	}
	b.WriteString(headerText)

	// Books list with progress
	var content strings.Builder
	for i, book := range m.books {
		// Get progress for this book
		progress := textgen.GetProgressForBook(&book)

		// Format the book entry with percent complete if available
		bookEntry := book.Name
		if progress != nil && progress.CharacterPos > 0 {
			// Show percent with 1 decimal place
			if progress.PercentComplete > 0 {
				bookEntry = fmt.Sprintf("%s (%.1f%%)", book.Name, progress.PercentComplete)
			} else {
				bookEntry = fmt.Sprintf("%s (0.0%%)", book.Name)
			}
		}

		if i == m.selectedIndex {
			// Highlight selected book
			content.WriteString(fmt.Sprintf("\033[1;33m▶ %s\033[0m\n", bookEntry))
		} else {
			content.WriteString(fmt.Sprintf("  %s\n", bookEntry))
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
		// Get progress for this book
		progress := textgen.GetProgressForBook(&book)

		// Format the book entry with percent complete if available
		bookEntry := book.Name
		if progress != nil && progress.CharacterPos > 0 {
			// Show percent with 1 decimal place
			if progress.PercentComplete > 0 {
				bookEntry = fmt.Sprintf("%s (%.1f%%)", book.Name, progress.PercentComplete)
			} else {
				bookEntry = fmt.Sprintf("%s (0.0%%)", book.Name)
			}
		}

		if i == m.selectedIndex {
			// Highlight selected book with yellow arrow
			content.WriteString(fmt.Sprintf("\033[1;33m▶ %s\033[0m\n", bookEntry))
		} else {
			content.WriteString(fmt.Sprintf("  %s\n", bookEntry))
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

// DocMenuModel represents the documentation selection menu state
type DocMenuModel struct {
	docs           []interface{} // Will hold Doc structs
	selectedIndex  int
	viewport       viewport.Model
	terminalWidth  int
	terminalHeight int
	selectedDoc    *string // Pointer to selected doc name
	done           bool
	searchMode     bool
	searchQuery    string
	searchResults  []int
	searchIndex    int
}

// NewDocMenuModel creates a new documentation selection menu
func NewDocMenuModel(docs []string, width, height int) *DocMenuModel {
	docInterfaces := make([]interface{}, len(docs))
	for i, doc := range docs {
		docInterfaces[i] = doc
	}

	m := &DocMenuModel{
		docs:           docInterfaces,
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
func (m *DocMenuModel) Init() tea.Cmd {
	return nil
}

// Update handles input
func (m *DocMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			if m.selectedIndex < len(m.docs)-1 {
				m.selectedIndex++
				m.syncViewport()
			}
		case "k", "up":
			// Move up
			if m.selectedIndex > 0 {
				m.selectedIndex--
				m.syncViewport()
			}
		case "enter":
			// Select current doc
			if m.selectedIndex < len(m.docs) {
				docName := m.docs[m.selectedIndex].(string)
				m.selectedDoc = &docName
				m.done = true
				return m, tea.Quit
			}
		case "/", "?":
			// Start search
			m.searchMode = true
			m.searchQuery = ""
		case "q", "esc":
			// Quit without selecting
			m.done = true
			return m, tea.Quit
		case "g":
			// Go to top
			m.selectedIndex = 0
			m.syncViewport()
		case "G":
			// Go to bottom
			m.selectedIndex = len(m.docs) - 1
			m.syncViewport()
		}

	case tea.WindowSizeMsg:
		m.terminalWidth = msg.Width
		m.terminalHeight = msg.Height
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - 4
		m.renderMenu()
	}

	return m, nil
}

// View renders the menu
func (m *DocMenuModel) View() string {
	return m.viewport.View()
}

// SelectedDocName returns the selected documentation name
func (m *DocMenuModel) SelectedDocName() *string {
	return m.selectedDoc
}

// Sync viewport to selected item
func (m *DocMenuModel) syncViewport() {
	m.renderMenu()

	if m.selectedIndex < len(m.docs) {
		// Ensure selected item is visible
		m.viewport.YPosition = 3
		itemHeight := 2 // Each item takes 2 lines (title + space)
		visibleItems := m.viewport.Height / itemHeight

		if m.selectedIndex < m.viewport.YOffset/itemHeight {
			m.viewport.YOffset = m.selectedIndex * itemHeight
		} else if m.selectedIndex >= (m.viewport.YOffset/itemHeight)+visibleItems {
			m.viewport.YOffset = (m.selectedIndex - visibleItems + 1) * itemHeight
		}
	}
}

// Render the menu content
func (m *DocMenuModel) renderMenu() {
	var buf strings.Builder
	buf.WriteString("Available Go Documentation\n")
	buf.WriteString("============================\n\n")

	for i, doc := range m.docs {
		docName := doc.(string)
		if i == m.selectedIndex {
			// Highlight selected doc with yellow arrow
			buf.WriteString(fmt.Sprintf("\033[1;33m▶ %s\033[0m\n", docName))
		} else {
			buf.WriteString(fmt.Sprintf("  %s\n", docName))
		}
	}

	m.viewport.SetContent(buf.String())
}

// Perform search on docs
func (m *DocMenuModel) performSearch() {
	if m.searchQuery == "" {
		m.searchResults = nil
		return
	}

	query := strings.ToLower(m.searchQuery)
	m.searchResults = nil

	// Search through all docs
	for i, doc := range m.docs {
		docName := doc.(string)
		if strings.Contains(strings.ToLower(docName), query) {
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
