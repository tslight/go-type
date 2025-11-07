package cli

// selection.go centralizes interactive content selection logic used by both apps.
// It removes duplication between gutentype and doctype main.go files.

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tobe/go-type/internal/content"
)

// SelectContent presents an interactive menu and returns a Selection.
// Behavior differs based on manifest usage in the manager:
//   - Manifest-based (e.g., books): use NewMenuModel and select by numeric ID.
//   - Directory-based (e.g., godocs): use NewDocMenuModel and select by name.
//
// Returns (nil, nil) if the user aborts selection.
func SelectContent(manager *content.ContentManager, width, height int) (*Selection, error) {
	if manager == nil {
		return nil, nil
	}

	if managerUsesManifest(manager) {
		menuModel := NewMenuModel(manager, width, height)
		program := tea.NewProgram(menuModel)
		if _, err := program.Run(); err != nil {
			return nil, err
		}
		selected := menuModel.SelectedBook()
		if selected == nil {
			return nil, nil
		}
		if err := manager.SetContent(selected.ID); err != nil {
			return nil, err
		}
		text := manager.GetCurrentText()
		provider := NewBookStateProvider(manager, selected.ID, len(text))
		return &Selection{Text: text, Book: selected, Provider: provider}, nil
	}

	// Directory-based selection path (e.g., docs)
	docModel := NewDocMenuModel(manager, width, height)
	program := tea.NewProgram(docModel)
	if _, err := program.Run(); err != nil {
		return nil, err
	}
	namePtr := docModel.SelectedDocName()
	if namePtr == nil {
		return nil, nil
	}
	name := *namePtr
	if err := manager.SetContentByName(name); err != nil {
		return nil, err
	}
	text := manager.GetCurrentText()
	provider := NewDocStateProvider(manager, name, len(text))
	// The Book field historically carried content meta; reuse current content.
	current := manager.GetCurrentContent()
	if current == nil { // Fallback if not set for some reason
		current = &content.Content{ID: 0, Name: name, Text: text}
	}
	return &Selection{Text: text, Book: current, Provider: provider}, nil
}

// managerUsesManifest is a small indirection to avoid exposing useManifest outside the content package.
// We rely on presence of numeric IDs loaded from manifest (ID > 0) as heuristic; manifest loading assigns positive IDs.
func managerUsesManifest(manager *content.ContentManager) bool {
	// If there is at least one available content with ID > 0 and names do not contain slashes which we inject for directory entries
	// this is a lightweight inference; alternatively we could expose a method on ContentManager.
	contents := manager.GetAvailableContent()
	if len(contents) == 0 {
		return false
	}
	// If first item has ID > 0 and doesn't include '/' treat as manifest-based (books). Directory-based IDs start at 0 but we also transform names with '/'.
	// Books manifest intentionally assigns IDs starting at 1.
	if contents[0].ID > 0 {
		return true
	}
	return false
}
