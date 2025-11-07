package state

import (
	"fmt"
	"time"
)

type SessionResult struct {
	Timestamp time.Time `json:"timestamp"`
	WPM       float64   `json:"wpm"`
	Accuracy  float64   `json:"accuracy"`
	Errors    int       `json:"errors"`
	CharTyped int       `json:"characters_typed"`
	Duration  int       `json:"duration_seconds"`
}

type ContentState struct {
	ID              string          `json:"id"`
	Name            string          `json:"name"`
	CharacterPos    int             `json:"character_position"`
	LastHash        string          `json:"last_hash,omitempty"`
	TextLength      int             `json:"text_length"`
	PercentComplete float64         `json:"percent_complete"`
	Sessions        []SessionResult `json:"sessions"`
}

type ContentStateManager struct {
	store *Manager[string, ContentState]
}

func NewContentStateManager(appName string) *ContentStateManager {
	store := NewManager[string, ContentState](appName, func(state *ContentState) (string, bool) {
		if state.ID == "" {
			return "", false
		}
		return state.ID, true
	}, func(state *ContentState) {})
	return &ContentStateManager{store: store}
}

func (csm *ContentStateManager) Configure(appName string) error   { return csm.store.Configure(appName) }
func (csm *ContentStateManager) GetState(id string) *ContentState { return csm.store.Get(id) }
func (csm *ContentStateManager) GetCharPos(id string) int {
	if state := csm.GetState(id); state != nil {
		return state.CharacterPos
	}
	return 0
}

func (csm *ContentStateManager) SaveProgress(id, name string, charPos, textLength int, lastHash string) error {
	if id == "" {
		return fmt.Errorf("state: content ID cannot be empty")
	}
	percentComplete := 0.0
	if textLength > 0 && charPos > 0 {
		percentComplete = (float64(charPos) / float64(textLength)) * 100.0
	}
	sessions := []SessionResult{}
	if existing := csm.store.Get(id); existing != nil {
		sessions = existing.Sessions
	}
	state := ContentState{ID: id, Name: name, CharacterPos: charPos, LastHash: lastHash, TextLength: textLength, PercentComplete: percentComplete, Sessions: sessions}
	return csm.store.Set(state)
}

func (csm *ContentStateManager) RecordSession(id, name string, wpm, accuracy float64, errors, charTyped, duration int) error {
	if id == "" {
		return fmt.Errorf("state: content ID cannot be empty")
	}
	s := csm.store.Get(id)
	if s == nil {
		if err := csm.store.Set(ContentState{ID: id, Name: name, Sessions: []SessionResult{}}); err != nil {
			return err
		}
		s = csm.store.Get(id)
	}
	if s == nil {
		return fmt.Errorf("state: content state not initialized")
	}
	s.Sessions = append(s.Sessions, SessionResult{Timestamp: time.Now(), WPM: wpm, Accuracy: accuracy, Errors: errors, CharTyped: charTyped, Duration: duration})
	return csm.store.Save()
}

func (csm *ContentStateManager) ClearState(id string) error { return csm.store.Delete(id) }

func (csm *ContentStateManager) GetStats(id string) map[string]interface{} {
	state := csm.store.Get(id)
	if state == nil || len(state.Sessions) == 0 {
		return map[string]interface{}{"sessions_completed": 0, "total_time": 0, "average_wpm": 0.0, "best_wpm": 0.0, "average_accuracy": 0.0, "total_characters": 0}
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
	return map[string]interface{}{"sessions_completed": count, "total_time": totalTime, "average_wpm": totalWPM / float64(count), "best_wpm": bestWPM, "average_accuracy": totalAccuracy / float64(count), "total_characters": totalChars}
}

func (csm *ContentStateManager) FormatStats(stats map[string]interface{}, title string) string {
	if len(stats) == 0 {
		return ""
	}
	sessionsCompleted, _ := stats["sessions_completed"].(int)
	totalTime, _ := stats["total_time"].(int)
	averageWPM, _ := stats["average_wpm"].(float64)
	bestWPM, _ := stats["best_wpm"].(float64)
	averageAccuracy, _ := stats["average_accuracy"].(float64)
	totalChars, _ := stats["total_characters"].(int)
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
	return fmt.Sprintf("\n📊 %s\n──────────────────────────────────\nSessions Completed:  %d\nTotal Time:          %s\nAverage WPM:         %.1f\nBest WPM:            %.1f\nAverage Accuracy:    %.1f%%\nTotal Characters:    %d\n──────────────────────────────────\n", title, sessionsCompleted, timeStr, averageWPM, bestWPM, averageAccuracy, totalChars)
}
