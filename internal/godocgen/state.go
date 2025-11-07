package godocgen

import (
	"fmt"
	"time"

	"github.com/tobe/go-type/internal/statestore"
)

const defaultDocAppName = "doctype"

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
	store *statestore.Manager[string, DocState]
}

var docManager = newDocStateManager()

func newDocStateManager() *docStateManager {
	store := statestore.NewManager[string, DocState](
		defaultDocAppName,
		func(state *DocState) (string, bool) {
			if state.DocName == "" {
				return "", false
			}
			return state.DocName, true
		},
		func(state *DocState) {},
	)
	return &docStateManager{store: store}
}

func (dm *docStateManager) configure(appName string) error {
	return dm.store.Configure(appName)
}

// ConfigureStateFile overrides the state file name based on the provided app name.
func ConfigureStateFile(appName string) error {
	return docManager.configure(appName)
}

// GetDocState returns the state for a doc, if any
func GetDocState(docName string) *DocState {
	return docManager.store.Get(docName)
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
	if existing := docManager.store.Get(docName); existing != nil {
		sessions = existing.Sessions
	}

	state := DocState{
		DocName:         docName,
		CharacterPos:    charPos,
		TextLength:      textLength,
		PercentComplete: percentComplete,
		Sessions:        sessions,
	}

	return docManager.store.Set(state)
}

// RecordDocSession appends a session result for a doc
func RecordDocSession(docName string, wpm, accuracy float64, errors, charTyped, duration int) error {
	if docName == "" {
		return fmt.Errorf("doc name cannot be empty")
	}

	state := docManager.store.Get(docName)
	if state == nil {
		if err := docManager.store.Set(DocState{DocName: docName, Sessions: []DocSessionResult{}}); err != nil {
			return err
		}
		state = docManager.store.Get(docName)
	}
	if state == nil {
		return fmt.Errorf("doc state not initialized")
	}

	state.Sessions = append(state.Sessions, DocSessionResult{
		Timestamp: time.Now(),
		WPM:       wpm,
		Accuracy:  accuracy,
		Errors:    errors,
		CharTyped: charTyped,
		Duration:  duration,
	})

	return docManager.store.Save()
}

// GetDocStats returns aggregated statistics for a doc
func GetDocStats(docName string) map[string]interface{} {
	state := docManager.store.Get(docName)
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

	return fmt.Sprintf(
		"\n📊 DOCUMENT STATISTICS\n"+
			"──────────────────────────────────\n"+
			"Sessions Completed:  %d\n"+
			"Total Time:          %s\n"+
			"Average WPM:         %.1f\n"+
			"Best WPM:            %.1f\n"+
			"Average Accuracy:    %.1f%%\n"+
			"Total Characters:    %d\n"+
			"──────────────────────────────────\n",
		sessionsCompleted, timeStr, averageWPM, bestWPM, averageAccuracy, totalChars,
	)
}
