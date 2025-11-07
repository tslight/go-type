package godocgen

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// DocSessionResult represents a single typing session for a doc
type DocSessionResult struct {
	Timestamp time.Time `json:"timestamp"`
	WPM       float64   `json:"wpm"`
	Accuracy  float64   `json:"accuracy"`
	Errors    int       `json:"errors"`
	CharTyped int       `json:"characters_typed"`
	Duration  int       `json:"duration_seconds"`
}

// DocState stores progress and sessions for a documentation module
type DocState struct {
	DocName         string             `json:"doc_name"`
	CharacterPos    int                `json:"character_position"`
	TextLength      int                `json:"text_length"`
	PercentComplete float64            `json:"percent_complete"`
	Sessions        []DocSessionResult `json:"sessions"`
}

type docStateManager struct {
	stateFile string
	states    map[string]*DocState
}

var docManager = newDocStateManager()

func newDocStateManager() *docStateManager {
	dm := &docStateManager{states: make(map[string]*DocState)}
	dm.stateFile = dm.getStateFilePath()
	_ = dm.loadStates()
	return dm
}

func (dm *docStateManager) getStateFilePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	return filepath.Join(home, ".go-type-docs-state.json")
}

func (dm *docStateManager) loadStates() error {
	data, err := os.ReadFile(dm.stateFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	var states []DocState
	if err := json.Unmarshal(data, &states); err != nil {
		return err
	}

	for i := range states {
		state := states[i]
		if state.DocName == "" {
			continue
		}
		dm.states[state.DocName] = &states[i]
	}
	return nil
}

func (dm *docStateManager) saveStates() error {
	states := make([]DocState, 0, len(dm.states))
	for _, s := range dm.states {
		states = append(states, *s)
	}

	data, err := json.MarshalIndent(states, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(dm.stateFile, data, 0644)
}

// GetDocState returns the state for a doc, if any
func GetDocState(docName string) *DocState {
	return docManager.states[docName]
}

// GetSavedCharPos returns the saved character position for a doc
func GetSavedCharPos(docName string) int {
	if state := GetDocState(docName); state != nil {
		return state.CharacterPos
	}
	return 0
}

// SaveDocProgress stores the current progress for a doc
func SaveDocProgress(docName string, charPos int, textLength int) error {
	if docName == "" {
		return fmt.Errorf("doc name cannot be empty")
	}

	percentComplete := 0.0
	if textLength > 0 && charPos > 0 {
		percentComplete = (float64(charPos) / float64(textLength)) * 100.0
	}

	sessions := []DocSessionResult{}
	if existing := docManager.states[docName]; existing != nil {
		sessions = existing.Sessions
	}

	docManager.states[docName] = &DocState{
		DocName:         docName,
		CharacterPos:    charPos,
		TextLength:      textLength,
		PercentComplete: percentComplete,
		Sessions:        sessions,
	}

	return docManager.saveStates()
}

// RecordDocSession appends a session result for a doc
func RecordDocSession(docName string, wpm, accuracy float64, errors, charTyped, duration int) error {
	if docName == "" {
		return fmt.Errorf("doc name cannot be empty")
	}

	state := docManager.states[docName]
	if state == nil {
		state = &DocState{DocName: docName, Sessions: []DocSessionResult{}}
		docManager.states[docName] = state
	}

	state.Sessions = append(state.Sessions, DocSessionResult{
		Timestamp: time.Now(),
		WPM:       wpm,
		Accuracy:  accuracy,
		Errors:    errors,
		CharTyped: charTyped,
		Duration:  duration,
	})

	return docManager.saveStates()
}

// GetDocStats returns aggregated statistics for a doc
func GetDocStats(docName string) map[string]interface{} {
	state := docManager.states[docName]
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

// FormatDocStats returns a formatted stats string for display
func FormatDocStats(stats map[string]interface{}) string {
	if len(stats) == 0 {
		return ""
	}

	sessionsCompleted, _ := stats["sessions_completed"].(int)
	if sessionsCompleted == 0 {
		return "\nNo previous sessions recorded for this document."
	}

	totalTime, _ := stats["total_time"].(int)
	averageWPM, _ := stats["average_wpm"].(float64)
	bestWPM, _ := stats["best_wpm"].(float64)
	averageAccuracy, _ := stats["average_accuracy"].(float64)
	totalChars, _ := stats["total_characters"].(int)

	return fmt.Sprintf("\nPrevious Sessions:\nSessions completed: %d\nTotal time: %d seconds\nAverage WPM: %.2f\nBest WPM: %.2f\nAverage accuracy: %.2f%%\nTotal characters typed: %d\n",
		sessionsCompleted, totalTime, averageWPM, bestWPM, averageAccuracy, totalChars)
}
