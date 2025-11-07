package cli

// Unified content selection logic using a single MenuModel.

import (
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

	// Load based on manager mode.
	if manager.IsManifestBased() {
		if err := manager.SetContent(selected.ID); err != nil {
			return nil, err
		}
	} else {
		if err := manager.SetContentByName(selected.Name); err != nil {
			return nil, err
		}
	}

	text := manager.GetCurrentText()

	// Keep terminology content-agnostic throughout the CLI/package layer.
	provider := NewContentStateProvider(manager, manager.StateKeyFor(*selected), len(text), "CONTENT STATISTICS")
	return &Selection{Text: text, Content: manager.GetCurrentContent(), Provider: provider}, nil
}
