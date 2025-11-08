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
	CharTyped int       `json:"characters_typed_raw"`
	Effective int       `json:"characters_typed_effective"`
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

func (csm *ContentStateManager) RecordSession(id, name string, wpm, accuracy float64, errors, charTypedRaw, effectiveChars, duration int) error {
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
	s.Sessions = append(s.Sessions, SessionResult{Timestamp: time.Now(), WPM: wpm, Accuracy: accuracy, Errors: errors, CharTyped: charTypedRaw, Effective: effectiveChars, Duration: duration})
	return csm.store.Save()
}

func (csm *ContentStateManager) ClearState(id string) error { return csm.store.Delete(id) }

// WipeAllStates deletes the underlying state file and clears all in-memory states.
func (csm *ContentStateManager) WipeAllStates() error { return csm.store.WipeFile() }

func (csm *ContentStateManager) GetStats(id string) map[string]interface{} {
	state := csm.store.Get(id)
	if state == nil || len(state.Sessions) == 0 {
		// Include current text progress even if there are no sessions yet.
		typed := 0
		total := 0
		if state != nil {
			typed = state.CharacterPos
			total = state.TextLength
		}
		return map[string]interface{}{
			"sessions_completed":         0,
			"total_time":                 0,
			"average_wpm":                0.0,
			"best_wpm":                   0.0,
			"average_accuracy":           0.0,
			"total_characters_raw":       0,
			"total_characters_effective": 0,
			"text_progress_typed":        typed,
			"text_progress_total":        total,
		}
	}
	totalAccuracy := 0.0
	totalTime := 0
	totalCharsRaw := 0
	totalCharsEffective := 0
	bestWPM := 0.0
	for _, session := range state.Sessions {
		// Recompute WPM from raw characters and duration to avoid legacy inflated values.
		recomputed := 0.0
		if session.Duration > 0 && session.CharTyped > 0 {
			minutes := float64(session.Duration) / 60.0
			recomputed = (float64(session.CharTyped) / 5.0) / minutes
		}
		totalAccuracy += session.Accuracy
		totalTime += session.Duration
		totalCharsRaw += session.CharTyped
		totalCharsEffective += session.Effective
		if recomputed > bestWPM {
			bestWPM = recomputed
		}
	}
	count := len(state.Sessions)
	// Aggregate WPM as total chars / total time (weighted by duration)
	avgWPM := 0.0
	if totalTime > 0 && totalCharsRaw > 0 {
		totalMinutes := float64(totalTime) / 60.0
		avgWPM = (float64(totalCharsRaw) / 5.0) / totalMinutes
	}
	return map[string]interface{}{
		"sessions_completed":         count,
		"total_time":                 totalTime,
		"average_wpm":                avgWPM,
		"best_wpm":                   bestWPM,
		"average_accuracy":           totalAccuracy / float64(count),
		"total_characters_raw":       totalCharsRaw,
		"total_characters_effective": totalCharsEffective,
		"text_progress_typed":        state.CharacterPos,
		"text_progress_total":        state.TextLength,
	}
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
	totalCharsRaw, _ := stats["total_characters_raw"].(int)
	totalCharsEff, _ := stats["total_characters_effective"].(int)
	textTyped, _ := stats["text_progress_typed"].(int)
	textTotal, _ := stats["text_progress_total"].(int)
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
	return fmt.Sprintf("\n📊 %s\n──────────────────────────────────\nSessions Completed:      %d\nTotal Time:              %s\nAverage WPM:             %.1f\nBest WPM:                %.1f\nAverage Accuracy:        %.1f%%\nText Progress:           %d/%d\nTotal Characters (raw):  %d\nTotal Characters (eff):  %d\n──────────────────────────────────\n", title, sessionsCompleted, timeStr, averageWPM, bestWPM, averageAccuracy, textTyped, textTotal, totalCharsRaw, totalCharsEff)
}

// GetGlobalStats aggregates stats across all content states.
// WPM values are recomputed per session (same logic as GetStats).
func (csm *ContentStateManager) GetGlobalStats() map[string]interface{} {
	all := csm.store.AllStates()
	totalSessions := 0
	totalTime := 0
	totalCharsRaw := 0
	totalCharsEff := 0
	totalAccuracy := 0.0
	bestWPM := 0.0
	for _, st := range all {
		for _, sess := range st.Sessions {
			totalSessions++
			totalTime += sess.Duration
			totalCharsRaw += sess.CharTyped
			totalCharsEff += sess.Effective
			recomputed := 0.0
			if sess.Duration > 0 && sess.CharTyped > 0 {
				minutes := float64(sess.Duration) / 60.0
				recomputed = (float64(sess.CharTyped) / 5.0) / minutes
			}
			totalAccuracy += sess.Accuracy
			if recomputed > bestWPM {
				bestWPM = recomputed
			}
		}
	}
	if totalSessions == 0 {
		return map[string]interface{}{"sessions_completed": 0, "total_time": 0, "average_wpm": 0.0, "best_wpm": 0.0, "average_accuracy": 0.0, "total_characters_raw": 0, "total_characters_effective": 0, "text_progress_typed": 0, "text_progress_total": 0}
	}
	return map[string]interface{}{
		"sessions_completed": totalSessions,
		"total_time":         totalTime,
		"average_wpm": func() float64 {
			if totalTime == 0 || totalCharsRaw == 0 {
				return 0.0
			}
			return (float64(totalCharsRaw) / 5.0) / (float64(totalTime) / 60.0)
		}(),
		"best_wpm":                   bestWPM,
		"average_accuracy":           totalAccuracy / float64(totalSessions),
		"total_characters_raw":       totalCharsRaw,
		"total_characters_effective": totalCharsEff,
		// Global progress (optional): treat typed as sum of per-state CharacterPos, total as sum of TextLength.
		"text_progress_typed": func() int {
			sum := 0
			for _, st := range all {
				sum += st.CharacterPos
			}
			return sum
		}(),
		"text_progress_total": func() int {
			sum := 0
			for _, st := range all {
				sum += st.TextLength
			}
			return sum
		}(),
	}
}
