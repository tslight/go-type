package textgen

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// BookState represents the typing progress for a book
type BookState struct {
	BookID          int     `json:"book_id"`
	BookName        string  `json:"book_name"`
	CharacterPos    int     `json:"character_position"`
	LastHash        string  `json:"last_hash"`
	TextLength      int     `json:"text_length"`
	PercentComplete float64 `json:"percent_complete"`
}

// StateManager handles loading and saving book progress
type StateManager struct {
	stateFile string
	states    map[int]*BookState
}

// NewStateManager creates a new state manager
func NewStateManager() *StateManager {
	sm := &StateManager{
		states: make(map[int]*BookState),
	}
	sm.stateFile = sm.getStateFilePath()
	sm.loadStates()
	return sm
}

// getStateFilePath returns the path to the state file
func (sm *StateManager) getStateFilePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	return filepath.Join(home, ".go-type-state.json")
}

// loadStates loads the state file from disk
func (sm *StateManager) loadStates() error {
	data, err := os.ReadFile(sm.stateFile)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist yet, that's fine
			return nil
		}
		return err
	}

	var states []BookState
	if err := json.Unmarshal(data, &states); err != nil {
		return err
	}

	for i := range states {
		// Migrate old format to new format
		migrateBookState(&states[i])
		sm.states[states[i].BookID] = &states[i]
	}
	return nil
}

// migrateBookState handles backward compatibility with old state format
func migrateBookState(bs *BookState) {
	// No migration needed - CharacterPos is our primary field now
} // saveStates saves the current states to disk
func (sm *StateManager) saveStates() error {
	states := make([]BookState, 0, len(sm.states))
	for _, state := range sm.states {
		states = append(states, *state)
	}

	data, err := json.MarshalIndent(states, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(sm.stateFile, data, 0644)
}

// GetState returns the saved state for a book
func (sm *StateManager) GetState(bookID int) *BookState {
	return sm.states[bookID]
}

// SaveState saves the state for a book
func (sm *StateManager) SaveState(bookID int, bookName string, characterPos int, lastHash string) error {
	// Calculate text length and percent complete
	textLength := len(GetFullText())
	percentComplete := 0.0
	if textLength > 0 && characterPos > 0 {
		percentComplete = (float64(characterPos) / float64(textLength)) * 100.0
	}

	sm.states[bookID] = &BookState{
		BookID:          bookID,
		BookName:        bookName,
		CharacterPos:    characterPos,
		LastHash:        lastHash,
		TextLength:      textLength,
		PercentComplete: percentComplete,
	}
	return sm.saveStates()
}

// ClearState removes the saved state for a book
func (sm *StateManager) ClearState(bookID int) error {
	delete(sm.states, bookID)
	return sm.saveStates()
}
