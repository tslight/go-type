package textgen

import (
	"strconv"

	"github.com/tobe/go-type/internal/statestore"
)

// StateManager handles loading and saving book progress
type StateManager struct {
	csm *statestore.ContentStateManager
}

// NewStateManager creates a new state manager for book progress
func NewStateManager() *StateManager {
	return &StateManager{
		csm: statestore.NewContentStateManager("gutentype"),
	}
}

// GetState returns the saved state for a book
func (sm *StateManager) GetState(bookID int) *statestore.ContentState {
	return sm.csm.GetState(strconv.Itoa(bookID))
}

// SaveState stores the current book progress
func (sm *StateManager) SaveState(bookID int, bookName string, characterPos int, lastHash string) error {
	textLength := len(GetFullText())
	return sm.csm.SaveProgress(strconv.Itoa(bookID), bookName, characterPos, textLength, lastHash)
}

// ClearState removes the saved state for a book
func (sm *StateManager) ClearState(bookID int) error {
	return sm.csm.ClearState(strconv.Itoa(bookID))
}

// AddSession adds a new session result to a book's history
func (sm *StateManager) AddSession(bookID int, result statestore.SessionResult) error {
	// Note: Converting SessionResult which was in textgen to use statestore.SessionResult
	return sm.csm.RecordSession(
		strconv.Itoa(bookID),
		"", // bookName not needed for recording session
		result.WPM,
		result.Accuracy,
		result.Errors,
		result.CharTyped,
		result.Duration,
	)
}

// GetStats returns aggregated statistics for a book
func (sm *StateManager) GetStats(bookID int) map[string]interface{} {
	return sm.csm.GetStats(strconv.Itoa(bookID))
}
