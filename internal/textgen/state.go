package textgen

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// SessionResult represents a single typing session
type SessionResult struct {
	Timestamp time.Time `json:"timestamp"`
	WPM       float64   `json:"wpm"`
	Accuracy  float64   `json:"accuracy"`
	Errors    int       `json:"errors"`
	CharTyped int       `json:"characters_typed"`
	Duration  int       `json:"duration_seconds"`
}

// BookState represents the typing progress for a book
type BookState struct {
	BookID          int             `json:"book_id"`
	BookName        string          `json:"book_name"`
	CharacterPos    int             `json:"character_position"`
	LastHash        string          `json:"last_hash"`
	TextLength      int             `json:"text_length"`
	PercentComplete float64         `json:"percent_complete"`
	Sessions        []SessionResult `json:"sessions"`
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
	_ = sm.loadStates()
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

	// Preserve existing sessions if state already exists
	var sessions []SessionResult
	if existingState := sm.states[bookID]; existingState != nil {
		sessions = existingState.Sessions
	} else {
		sessions = []SessionResult{}
	}

	sm.states[bookID] = &BookState{
		BookID:          bookID,
		BookName:        bookName,
		CharacterPos:    characterPos,
		LastHash:        lastHash,
		TextLength:      textLength,
		PercentComplete: percentComplete,
		Sessions:        sessions,
	}
	return sm.saveStates()
}

// ClearState removes the saved state for a book
func (sm *StateManager) ClearState(bookID int) error {
	delete(sm.states, bookID)
	return sm.saveStates()
}

// AddSession adds a new session result to a book's history
func (sm *StateManager) AddSession(bookID int, result SessionResult) error {
	state := sm.GetState(bookID)
	if state == nil {
		return nil // No state for this book yet
	}
	state.Sessions = append(state.Sessions, result)
	return sm.saveStates()
}

// GetStats returns cumulative statistics for a book
func (sm *StateManager) GetStats(bookID int) map[string]interface{} {
	state := sm.GetState(bookID)
	if state == nil || len(state.Sessions) == 0 {
		return map[string]interface{}{
			"sessions_completed": 0,
			"total_time":         0,
			"average_wpm":        0.0,
			"best_wpm":           0.0,
			"average_accuracy":   0.0,
			"total_characters":   0,
		}
	}

	totalWPM := 0.0
	totalAccuracy := 0.0
	totalTime := 0
	totalChars := 0
	bestWPM := 0.0

	for _, session := range state.Sessions {
		totalWPM += session.WPM
		totalAccuracy += session.Accuracy
		totalTime += session.Duration
		totalChars += session.CharTyped
		if session.WPM > bestWPM {
			bestWPM = session.WPM
		}
	}

	count := len(state.Sessions)
	return map[string]interface{}{
		"sessions_completed": count,
		"total_time":         totalTime,
		"average_wpm":        totalWPM / float64(count),
		"best_wpm":           bestWPM,
		"average_accuracy":   totalAccuracy / float64(count),
		"total_characters":   totalChars,
	}
}
