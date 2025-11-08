package state

import (
	"strings"
	"testing"
)

func TestContentState_SaveProgressAndGet(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	mgr := NewContentStateManager("test-app-state")
	id := "item-1"
	name := "First"
	if err := mgr.SaveProgress(id, name, 50, 200, ""); err != nil {
		t.Fatalf("SaveProgress failed: %v", err)
	}
	st := mgr.GetState(id)
	if st == nil {
		t.Fatalf("expected state to be saved and retrievable")
	}
	if st.PercentComplete <= 0.0 || st.PercentComplete >= 100.0 {
		t.Fatalf("unexpected percent complete: %v", st.PercentComplete)
	}
	if got := mgr.GetCharPos(id); got != 50 {
		t.Fatalf("GetCharPos mismatch: got %d", got)
	}
}

func TestContentState_RecordSessionAndStats(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	mgr := NewContentStateManager("test-app-state")
	id := "item-2"
	name := "Second"
	// Record a couple of sessions
	if err := mgr.RecordSession(id, name, 50.0, 95.0, 2, 300, 250, 60); err != nil {
		t.Fatalf("RecordSession failed: %v", err)
	}
	if err := mgr.RecordSession(id, name, 70.0, 97.0, 1, 500, 450, 120); err != nil {
		t.Fatalf("RecordSession failed: %v", err)
	}
	stats := mgr.GetStats(id)
	if stats["sessions_completed"].(int) != 2 {
		t.Fatalf("expected 2 sessions, got %v", stats["sessions_completed"])
	}
	if stats["total_time"].(int) != 180 {
		t.Fatalf("expected total time 180, got %v", stats["total_time"])
	}
	if stats["total_characters_raw"].(int) != 800 {
		t.Fatalf("expected total raw chars 800, got %v", stats["total_characters_raw"])
	}
	if stats["total_characters_effective"].(int) != 700 {
		t.Fatalf("expected total effective chars 700, got %v", stats["total_characters_effective"])
	}
	// Ensure formatting is non-empty and contains the title
	formatted := mgr.FormatStats(stats, "TEST STATS")
	if formatted == "" || !strings.Contains(formatted, "TEST STATS") {
		t.Fatalf("unexpected formatted stats: %q", formatted)
	}
}

func TestContentState_AggregatedWPMRecompute(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	mgr := NewContentStateManager("wpm-app")
	id := "wpm-item"
	// Two sessions: 250 chars in 120s -> (250/5)/2 = 25 WPM; 100 chars in 30s -> (100/5)/0.5 = 40 WPM
	if err := mgr.RecordSession(id, "Name", 999.0, 90.0, 0, 250, 240, 120); err != nil { // inflate stored WPM to test recompute
		t.Fatalf("RecordSession failed: %v", err)
	}
	if err := mgr.RecordSession(id, "Name", 888.0, 92.0, 1, 100, 95, 30); err != nil {
		t.Fatalf("RecordSession failed: %v", err)
	}
	stats := mgr.GetStats(id)
	avg := stats["average_wpm"].(float64)
	best := stats["best_wpm"].(float64)
	// Time-weighted average: total chars 350 -> 70 words; total time 150s -> 2.5 minutes; 70/2.5 = 28
	if avg < 27.9 || avg > 28.1 {
		t.Fatalf("expected average WPM ~28.0, got %f", avg)
	}
	if best < 39.9 || best > 40.1 {
		t.Fatalf("expected best WPM ~40, got %f", best)
	}
}

func TestContentState_WhitespaceOnlySession(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	mgr := NewContentStateManager("test-app-state")
	id := "item-ws"
	name := "WS"
	// Record whitespace-only session: raw>0, effective=0
	if err := mgr.RecordSession(id, name, 10.0, 100.0, 0, 20, 0, 30); err != nil {
		t.Fatalf("RecordSession failed: %v", err)
	}
	stats := mgr.GetStats(id)
	if stats["total_characters_raw"].(int) != 20 {
		t.Fatalf("expected raw 20, got %v", stats["total_characters_raw"])
	}
	if stats["total_characters_effective"].(int) != 0 {
		t.Fatalf("expected effective 0, got %v", stats["total_characters_effective"])
	}
}

func TestContentState_ClearState(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	mgr := NewContentStateManager("test-app-state")
	id := "item-3"
	if err := mgr.SaveProgress(id, "Name", 10, 100, ""); err != nil {
		t.Fatalf("SaveProgress failed: %v", err)
	}
	if err := mgr.ClearState(id); err != nil {
		t.Fatalf("ClearState failed: %v", err)
	}
	if st := mgr.GetState(id); st != nil {
		t.Fatalf("expected state to be cleared")
	}
}

func TestContentState_EmptyIDErrors(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	mgr := NewContentStateManager("test-app-state")
	if err := mgr.SaveProgress("", "", 0, 0, ""); err == nil {
		t.Fatalf("expected error for empty ID in SaveProgress")
	}
	if err := mgr.RecordSession("", "", 0, 0, 0, 0, 0, 0); err == nil {
		t.Fatalf("expected error for empty ID in RecordSession")
	}
}

func TestBuildStateFileName_Variants(t *testing.T) {
	name, err := BuildStateFileName("App")
	if err != nil || name == "" {
		t.Fatalf("unexpected error or empty name: %v %q", err, name)
	}
	if _, err := BuildStateFileName(""); err == nil {
		t.Fatalf("expected error for empty app name")
	}
}

func TestManager_ConfigureAndAllStates(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	// Use the underlying generic manager through ContentStateManager
	mgr := NewContentStateManager("app")
	// Configure with new app name
	if err := mgr.Configure("another"); err != nil {
		t.Fatalf("configure: %v", err)
	}
	// Save two states
	_ = mgr.SaveProgress("id1", "n1", 1, 10, "")
	_ = mgr.SaveProgress("id2", "n2", 2, 10, "")
	// Access the underlying AllStates by retrieving and counting
	// (We can't call AllStates directly; assert through Get on both IDs)
	if mgr.GetState("id1") == nil || mgr.GetState("id2") == nil {
		t.Fatalf("expected both states present")
	}
	// StateFilePath is non-empty after configure since we wrote
	if mgr.store.StateFilePath() == "" {
		t.Fatalf("expected non-empty state file path")
	}
}

func TestManager_LoadStatesOnNewInstance(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	mgr1 := NewContentStateManager("app")
	_ = mgr1.SaveProgress("idA", "NameA", 3, 10, "")
	_ = mgr1.SaveProgress("idB", "NameB", 4, 10, "")
	// New manager should load from same state file
	mgr2 := NewContentStateManager("app")
	if mgr2.GetState("idA") == nil || mgr2.GetState("idB") == nil {
		t.Fatalf("expected states to be loaded by new manager instance")
	}
}

func TestContentState_TextProgressInStats(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	mgr := NewContentStateManager("app")
	// Save progress to set CharacterPos and TextLength
	_ = mgr.SaveProgress("idP", "NameP", 42, 100, "")
	// No sessions yet; GetStats should still include text progress
	stats := mgr.GetStats("idP")
	if stats["text_progress_typed"].(int) != 42 || stats["text_progress_total"].(int) != 100 {
		t.Fatalf("expected text progress 42/100, got %v/%v", stats["text_progress_typed"], stats["text_progress_total"])
	}
	formatted := mgr.FormatStats(stats, "TEST STATS")
	if formatted == "" || !strings.Contains(formatted, "Text Progress:           42/100") {
		t.Fatalf("expected formatted text progress, got %q", formatted)
	}
}

func TestContentState_WipeAllStates(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	mgr := NewContentStateManager("wipe-app")
	_ = mgr.SaveProgress("id1", "Name1", 10, 100, "")
	_ = mgr.SaveProgress("id2", "Name2", 20, 100, "")
	// Ensure states exist
	if mgr.GetState("id1") == nil || mgr.GetState("id2") == nil {
		t.Fatalf("expected states before wipe")
	}
	// Wipe
	if err := mgr.WipeAllStates(); err != nil {
		t.Fatalf("wipe failed: %v", err)
	}
	if mgr.GetState("id1") != nil || mgr.GetState("id2") != nil {
		t.Fatalf("expected all states cleared after wipe")
	}
	// Subsequent save should recreate file cleanly
	if err := mgr.SaveProgress("id3", "Name3", 5, 50, ""); err != nil {
		t.Fatalf("SaveProgress after wipe failed: %v", err)
	}
	if mgr.GetState("id3") == nil {
		t.Fatalf("expected new state after wipe and save")
	}
}

func TestContentState_GlobalStatsAggregation(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	mgr := NewContentStateManager("global-app")
	// Two different content items with sessions
	if err := mgr.RecordSession("idA", "NameA", 25.0, 95.0, 1, 250, 240, 120); err != nil { // 250 raw in 120s -> 25 WPM
		t.Fatalf("RecordSession A failed: %v", err)
	}
	if err := mgr.RecordSession("idB", "NameB", 40.0, 90.0, 2, 100, 95, 30); err != nil { // 100 raw in 30s -> 40 WPM
		t.Fatalf("RecordSession B failed: %v", err)
	}
	g := mgr.GetGlobalStats()
	if g["sessions_completed"].(int) != 2 {
		t.Fatalf("expected 2 sessions, got %v", g["sessions_completed"])
	}
	if g["total_characters_raw"].(int) != 350 {
		t.Fatalf("expected raw 350, got %v", g["total_characters_raw"])
	}
	avg := g["average_wpm"].(float64)
	best := g["best_wpm"].(float64)
	// Time-weighted average: same sessions -> ~28.0
	if avg < 27.9 || avg > 28.1 {
		t.Fatalf("expected avg wpm ~28.0, got %f", avg)
	}
	if best < 39.9 || best > 40.1 {
		t.Fatalf("expected best wpm ~40, got %f", best)
	}
	formatted := mgr.FormatStats(g, "GLOBAL STATS")
	if formatted == "" || !strings.Contains(formatted, "GLOBAL STATS") {
		t.Fatalf("expected formatted global stats, got %q", formatted)
	}
}

func TestContentStateManager_GlobalStatsEmpty(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	cm := NewContentStateManager("app")
	stats := cm.GetGlobalStats()
	if stats["sessions_completed"].(int) != 0 {
		t.Fatalf("expected 0 sessions, got %+v", stats)
	}
	if stats["average_wpm"].(float64) != 0 || stats["average_accuracy"].(float64) != 0 {
		t.Fatalf("expected zeroed averages, got %+v", stats)
	}
}
