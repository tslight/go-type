package statestore

import (
	"fmt"
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

// ContentState represents the typing progress for any content (book, doc, etc)
type ContentState struct {
	ID              string          `json:"id"`
	Name            string          `json:"name"`
	CharacterPos    int             `json:"character_position"`
	LastHash        string          `json:"last_hash,omitempty"`
	TextLength      int             `json:"text_length"`
	PercentComplete float64         `json:"percent_complete"`
	Sessions        []SessionResult `json:"sessions"`
}

// ContentStateManager handles loading and saving content progress
type ContentStateManager struct {
	store *Manager[string, ContentState]
}

// NewContentStateManager creates a new unified state manager
func NewContentStateManager(appName string) *ContentStateManager {
	store := NewManager[string, ContentState](
		appName,
		func(state *ContentState) (string, bool) {
			if state.ID == "" {
				return "", false
			}
			return state.ID, true
		},
		func(state *ContentState) {
			// Optional migration logic
		},
	)
	return &ContentStateManager{store: store}
}

// Configure allows changing the app name / state file path
func (csm *ContentStateManager) Configure(appName string) error {
	return csm.store.Configure(appName)
}

// GetState returns the saved state for content by ID
func (csm *ContentStateManager) GetState(id string) *ContentState {
	return csm.store.Get(id)
}

// GetCharPos returns the saved character position for content
func (csm *ContentStateManager) GetCharPos(id string) int {
	if state := csm.GetState(id); state != nil {
		return state.CharacterPos
	}
	return 0
}

// SaveProgress stores the current progress for content
func (csm *ContentStateManager) SaveProgress(id, name string, charPos, textLength int, lastHash string) error {
	if id == "" {
		return fmt.Errorf("statestore: content ID cannot be empty")
	}

	percentComplete := 0.0
	if textLength > 0 && charPos > 0 {
		percentComplete = (float64(charPos) / float64(textLength)) * 100.0
	}

	sessions := []SessionResult{}
	if existing := csm.store.Get(id); existing != nil {
		sessions = existing.Sessions
	}

	state := ContentState{
		ID:              id,
		Name:            name,
		CharacterPos:    charPos,
		LastHash:        lastHash,
		TextLength:      textLength,
		PercentComplete: percentComplete,
		Sessions:        sessions,
	}

	return csm.store.Set(state)
}

// RecordSession appends a session result to content history
func (csm *ContentStateManager) RecordSession(id, name string, wpm, accuracy float64, errors, charTyped, duration int) error {
	if id == "" {
		return fmt.Errorf("statestore: content ID cannot be empty")
	}

	state := csm.store.Get(id)
	if state == nil {
		// Initialize state if it doesn't exist
		if err := csm.store.Set(ContentState{
			ID:       id,
			Name:     name,
			Sessions: []SessionResult{},
		}); err != nil {
			return err
		}
		state = csm.store.Get(id)
	}
	if state == nil {
		return fmt.Errorf("statestore: content state not initialized")
	}

	state.Sessions = append(state.Sessions, SessionResult{
		Timestamp: time.Now(),
		WPM:       wpm,
		Accuracy:  accuracy,
		Errors:    errors,
		CharTyped: charTyped,
		Duration:  duration,
	})

	return csm.store.Save()
}

// ClearState removes the saved state for content
func (csm *ContentStateManager) ClearState(id string) error {
	return csm.store.Delete(id)
}

// GetStats returns aggregated statistics for content
func (csm *ContentStateManager) GetStats(id string) map[string]interface{} {
	state := csm.store.Get(id)
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

// FormatStats returns a formatted stats string for display
func (csm *ContentStateManager) FormatStats(stats map[string]interface{}, title string) string {
	if len(stats) == 0 {
		return ""
	}

	sessionsCompleted := 0
	if v, ok := stats["sessions_completed"].(int); ok {
		sessionsCompleted = v
	}

	totalTime := 0
	if v, ok := stats["total_time"].(int); ok {
		totalTime = v
	}

	averageWPM := 0.0
	if v, ok := stats["average_wpm"].(float64); ok {
		averageWPM = v
	}

	bestWPM := 0.0
	if v, ok := stats["best_wpm"].(float64); ok {
		bestWPM = v
	}

	averageAccuracy := 0.0
	if v, ok := stats["average_accuracy"].(float64); ok {
		averageAccuracy = v
	}

	totalChars := 0
	if v, ok := stats["total_characters"].(int); ok {
		totalChars = v
	}

	hours := totalTime / 3600
	minutes := (totalTime % 3600) / 60
	seconds := totalTime % 60

	var timeStr string
	if hours > 0 {
		timeStr = fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
	} else if minutes > 0 {
		timeStr = fmt.Sprintf("%dm %ds", minutes, seconds)
	} else {
		timeStr = fmt.Sprintf("%ds", seconds)
	}

	if title == "" {
		title = "STATISTICS"
	}

	return fmt.Sprintf(
		"\n📊 %s\n"+
			"──────────────────────────────────\n"+
			"Sessions Completed:  %d\n"+
			"Total Time:          %s\n"+
			"Average WPM:         %.1f\n"+
			"Best WPM:            %.1f\n"+
			"Average Accuracy:    %.1f%%\n"+
			"Total Characters:    %d\n"+
			"──────────────────────────────────\n",
		title, sessionsCompleted, timeStr, averageWPM, bestWPM, averageAccuracy, totalChars,
	)
}
