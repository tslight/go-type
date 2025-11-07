package cli

// Unified content selection logic using a single MenuModel.

import (
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tobe/go-type/internal/content"
)

// SelectContent runs the interactive menu, loads chosen content, and builds a Selection.
// (nil, nil) is returned if the user aborts.
func SelectContent(manager *content.ContentManager, width, height int) (*Selection, error) {
	if manager == nil {
		return nil, nil
	}
	menuModel := NewMenuModel(manager, width, height)
	program := tea.NewProgram(menuModel)
	if _, err := program.Run(); err != nil {
		return nil, err
	}
	selected := menuModel.SelectedContent()
	if selected == nil {
		return nil, nil
	}

	// Try ID-based load first (manifest scenario), fallback to name-based, tracking path.
	usedID := true
	if err := manager.SetContent(selected.ID); err != nil {
		usedID = false
		if err2 := manager.SetContentByName(selected.Name); err2 != nil {
			return nil, err // keep original error context
		}
	}

	text := manager.GetCurrentText()

	var contentID, statsTitle string
	if usedID {
		contentID = strconv.Itoa(selected.ID)
		statsTitle = "BOOK STATISTICS"
	} else {
		contentID = selected.Name
		statsTitle = "DOCUMENT STATISTICS"
	}

	provider := NewContentStateProvider(manager, contentID, len(text), statsTitle)
	return &Selection{Text: text, Book: manager.GetCurrentContent(), Provider: provider}, nil
}
