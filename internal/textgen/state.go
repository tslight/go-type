package textgen

import (
	"time"

	"github.com/tobe/go-type/internal/statestore"
)

const defaultAppName = "gutentype"

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

// StateManager handles loading and saving book progress using a shared statestore manager
type StateManager struct {
	store *statestore.Manager[int, BookState]
}

// NewStateManager creates a new state manager backed by statestore
func NewStateManager() *StateManager {
	store := statestore.NewManager[int, BookState](
		defaultAppName,
		func(state *BookState) (int, bool) {
			return state.BookID, true
		},
		migrateBookState,
	)
	return &StateManager{store: store}
}

func (sm *StateManager) configure(appName string) error {
	return sm.store.Configure(appName)
}

// ConfigureStateFile allows callers to override the state file based on app name
func ConfigureStateFile(appName string) error {
	return stateManager.configure(appName)
}

// migrateBookState handles backward compatibility with old state format
func migrateBookState(bs *BookState) {
	// No migration currently required.
}

// GetState returns the saved state for a book
func (sm *StateManager) GetState(bookID int) *BookState {
	return sm.store.Get(bookID)
}

// SaveState saves the state for a book
func (sm *StateManager) SaveState(bookID int, bookName string, characterPos int, lastHash string) error {
	textLength := len(GetFullText())
	percentComplete := 0.0
	if textLength > 0 && characterPos > 0 {
		percentComplete = (float64(characterPos) / float64(textLength)) * 100.0
	}

	sessions := []SessionResult{}
	if existingState := sm.store.Get(bookID); existingState != nil {
		sessions = existingState.Sessions
	}

	state := BookState{
		BookID:          bookID,
		BookName:        bookName,
		CharacterPos:    characterPos,
		LastHash:        lastHash,
		TextLength:      textLength,
		PercentComplete: percentComplete,
		Sessions:        sessions,
	}

	return sm.store.Set(state)
}

// ClearState removes the saved state for a book
func (sm *StateManager) ClearState(bookID int) error {
	return sm.store.Delete(bookID)
}

// AddSession adds a new session result to a book's history
func (sm *StateManager) AddSession(bookID int, result SessionResult) error {
	state := sm.store.Get(bookID)
	if state == nil {
		return nil
	}
	state.Sessions = append(state.Sessions, result)
	return sm.store.Save()
}

// GetStats returns cumulative statistics for a book
func (sm *StateManager) GetStats(bookID int) map[string]interface{} {
	state := sm.store.Get(bookID)
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
