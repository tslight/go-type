package cli

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tobe/go-type/internal/content"
)

// MenuModel represents the content selection menu state (books, docs, etc.)
type MenuModel struct {
	items           []content.Content
	selectedIndex   int
	viewport        viewport.Model
	terminalWidth   int
	terminalHeight  int
	selectedContent *content.Content
	done            bool
	searchMode      bool
	searchQuery     string
	searchDirection int // 1 for forward (/), -1 for backward (?)
	searchResults   []int
	searchIndex     int
	showingStats    bool // True when displaying stats for a content item
	statsIndex      int  // Index of item whose stats are being shown
	manager         *content.ContentManager
}

// NewMenuModel creates a new content selection menu
func NewMenuModel(manager *content.ContentManager, width, height int) *MenuModel {
	items := manager.GetAvailableContent()
	m := &MenuModel{
		items:          items,
		selectedIndex:  0,
		terminalWidth:  width,
		terminalHeight: height,
		viewport:       viewport.New(width, height-4),
		manager:        manager,
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
			if m.selectedIndex < len(m.items)-1 {
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
			m.selectedIndex = len(m.items) - 1
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
			// Show stats for selected item
			m.showingStats = true
			m.statsIndex = m.selectedIndex
			m.renderMenu()
		case "enter":
			// Select content
			m.selectedContent = &m.items[m.selectedIndex]
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
		if m.statsIndex < 0 || m.statsIndex >= len(m.items) {
			m.showingStats = false
			return m.View()
		}
		item := m.items[m.statsIndex]
		key := m.manager.StateKeyFor(item)
		stats := m.manager.StateManager.GetStats(key)
		statsStr := m.manager.StateManager.FormatStats(stats, "CONTENT STATISTICS")
		headerText := "\n\nContent: " + item.Name + "\n"
		b.WriteString(headerText)
		b.WriteString(statsStr)
		b.WriteString("\nPress any key to continue...\n")
		return b.String()
	}

	// Header
	if m.searchMode {
		prefix := "/"
		if m.searchDirection == -1 {
			prefix = "?"
		}
		b.WriteString(fmt.Sprintf("\nSelect content (searching... Press Enter to search, Esc to cancel)\n%s%s\n\n", prefix, m.searchQuery))
	} else {
		b.WriteString("\nSelect content (j/k navigate, / search, n/N next/prev result, i info, Enter select, q quit)\n\n")
	}
	m.viewport.SetContent(m.buildListContent())
	b.WriteString(m.viewport.View())

	return b.String()
}

// SelectedContent returns the selected content (if any)
func (m *MenuModel) SelectedContent() *content.Content {
	return m.selectedContent
}

// renderMenu updates the viewport content
func (m *MenuModel) renderMenu() {
	m.viewport.SetContent(m.buildListContent())
}

// syncViewport ensures the selected item is visible in the viewport
func (m *MenuModel) syncViewport() {
	m.renderMenu()

	// Calculate which line the selected item is on (one line per item)
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

// performSearch searches for items matching the query
func (m *MenuModel) performSearch() {
	if m.searchQuery == "" {
		m.searchResults = nil
		return
	}

	query := strings.ToLower(m.searchQuery)
	m.searchResults = nil

	// Search through all items
	for i, item := range m.items {
		if strings.Contains(strings.ToLower(item.Name), query) {
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

// buildListContent renders the selectable list with progress and highlighting.
func (m *MenuModel) buildListContent() string {
	var buf strings.Builder
	for i, item := range m.items {
		// Determine state key using manager helper
		key := m.manager.StateKeyFor(item)
		progress := m.manager.StateManager.GetState(key)

		entry := item.Name
		if progress != nil && progress.CharacterPos > 0 {
			if progress.PercentComplete > 0 {
				entry = fmt.Sprintf("%s (%.1f%%)", item.Name, progress.PercentComplete)
			} else {
				entry = fmt.Sprintf("%s (0.0%%)", item.Name)
			}
		}

		if i == m.selectedIndex {
			buf.WriteString(fmt.Sprintf("\033[1;33m▶ %s\033[0m\n", entry))
		} else {
			buf.WriteString(fmt.Sprintf("  %s\n", entry))
		}
	}
	return buf.String()
}
